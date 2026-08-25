package console

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"tunnelnet/internal/access"
	"tunnelnet/internal/audit"
	"tunnelnet/internal/cabin"
	"tunnelnet/internal/cable"
	"tunnelnet/internal/drain"
	"tunnelnet/internal/env"
	"tunnelnet/internal/fire"
	"tunnelnet/internal/ns"
	"tunnelnet/internal/patrol"
	"tunnelnet/internal/quota"
	"tunnelnet/internal/store"
	"tunnelnet/internal/vent"
)

type pumpUnit struct {
	runs int
}

func (p *pumpUnit) Start() error {
	p.runs++
	return nil
}

type Server struct {
	webDir   string
	namespaces *ns.Registry
	cabins   *cabin.Registry
	access   *access.Service
	env      *env.Service
	drain    *drain.Service
	fire     *fire.Service
	patrol   *patrol.Service
	quota    *quota.Service
	audit    *audit.Registry
	storage  *store.Store
	cables   *cable.Registry
	ventilation *vent.Service
	suppression *cabin.SuppressionLog
	powerBook  *cabin.PowerBook
	now      func() time.Time
}

func New(webDir string, storage *store.Store) *Server {
	namespaces := ns.NewRegistry()
	namespaces.Add("tunnel-a", "综合管廊A段")
	namespaces.Add("tunnel-b", "综合管廊B段")

	cabins := cabin.NewRegistry()
	cabins.Upsert(cabin.Cabin{ID: "cabin-1", Name: "东舱", Zone: "zone-1", Partition: "part-1"})
	cabins.Upsert(cabin.Cabin{ID: "cabin-2", Name: "西舱", Zone: "zone-2", Partition: "part-2"})
	cabins.Upsert(cabin.Cabin{ID: "cabin-3", Name: "南舱", Zone: "zone-3", Partition: "part-3"})

	partitions := cabin.NewPartitionTable()
	partitions.RePartition("part-1", "cabin-1")
	partitions.RePartition("part-2", "cabin-2")
	partitions.RePartition("part-3", "cabin-3")

	cableService := cable.NewService(3)
	environment := env.NewService(cableService, partitions, 75.0)
	environment.BindSensorAt("sensor-1", "part-1")
	environment.BindSensorAt("sensor-2", "part-2")
	environment.BindSensorAt("sensor-3", "part-3")

	accessService := access.NewService(storage, cabins)
	pump := &pumpUnit{}
	drainService := drain.NewService(cabins, pump, storage)
	drainService.AddPit(drain.Pit{ID: "pit-1", Cabin: "cabin-1", WaterLevel: 0.2})
	drainService.AddPit(drain.Pit{ID: "pit-2", Cabin: "cabin-2", WaterLevel: 0.1})

	ventilation := vent.NewService()
	suppression := cabin.NewSuppressionLog()
	dedupe := audit.NewMarkSet()
	fireService := fire.NewService(storage, ventilation, suppression, dedupe, []string{"zone-1", "zone-2", "zone-3"})
	fireService.AddDetector(fire.Detector{ID: "detector-1", Zone: "zone-1"})
	fireService.AddDetector(fire.Detector{ID: "detector-2", Zone: "zone-2"})

	patrolService := patrol.NewService(cabins)
	powerBook := cabin.NewPowerBook()
	powerBook.SetBaseline("cabin-1", 4200)
	powerBook.SetBaseline("cabin-2", 3800)
	powerBook.SetBaseline("cabin-3", 4600)
	quotaService := quota.NewService(cabins, powerBook)
	auditRegistry := audit.NewRegistry()
	cableRegistry := cable.NewRegistry()
	cableRegistry.Add(cable.Cable{ID: "cable-1", Cabin: "cabin-1", Temp: 42, Status: "normal"})
	cableRegistry.Add(cable.Cable{ID: "cable-2", Cabin: "cabin-2", Temp: 46, Status: "normal"})

	server := &Server{
		webDir:     webDir,
		namespaces: namespaces,
		cabins:     cabins,
		access:     accessService,
		env:        environment,
		drain:      drainService,
		fire:       fireService,
		patrol:     patrolService,
		quota:      quotaService,
		audit:      auditRegistry,
		storage:    storage,
		cables:     cableRegistry,
		ventilation: ventilation,
		suppression: suppression,
		powerBook:  powerBook,
		now:        time.Now,
	}
	server.seedDemo()
	return server
}

func (s *Server) Handler() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/cabins", http.StatusFound)
	})
	router.Get("/cabins", s.servePage("cabins.html"))
	router.Get("/env", s.servePage("env.html"))
	router.Get("/patrols", s.servePage("patrols.html"))
	router.Get("/alarms", s.servePage("alarms.html"))

	router.Get("/api/cabins", s.handleCabins)
	router.Get("/api/env", s.handleEnv)
	router.Get("/api/patrols", s.handlePatrols)
	router.Get("/api/alarms", s.handleAlarms)
	router.Get("/api/namespaces", s.handleNamespaces)
	router.Get("/api/quota/report", s.handleQuotaReport)
	router.Get("/api/audit", s.handleAudit)
	router.Post("/api/access/approve", s.handleAccessApprove)
	router.Post("/api/drain/start", s.handleDrainStart)
	router.Post("/api/fire/smoke", s.handleFireSmoke)
	router.Post("/api/fire/zone", s.handleFireZone)
	router.Post("/api/fire/alarm", s.handleFireAlarm)
	router.Post("/api/quota/refresh", s.handleQuotaRefresh)
	router.Post("/api/patrol/checkin", s.handlePatrolCheckIn)
	router.Post("/api/patrol/rebuild", s.handlePatrolRebuild)
	router.Post("/api/env/partition", s.handleEnvPartition)
	router.Post("/api/drain/execute", s.handleDrainExecute)
	router.Get("/api/status", s.handleStatus)
	router.Get("/api/report", s.handleReport)
	router.Get("/api/cables", s.handleCables)
	router.Post("/api/cables/temp", s.handleCableTemp)
	router.Get("/api/env/sensors", s.handleEnvSensors)
	router.Post("/api/env/sample", s.handleEnvSample)
	router.Get("/api/drain/pits", s.handleDrainPits)
	router.Post("/api/drain/pit", s.handleDrainPit)
	router.Get("/api/quota/usage", s.handleQuotaUsage)
	router.Post("/api/quota/usage", s.handleQuotaAddUsage)
	router.Get("/api/audit/query", s.handleAuditQuery)
	router.Get("/api/fire/alarms", s.handleFireAlarms)
	router.Post("/api/fire/clear", s.handleFireClear)
	router.Post("/api/vent/toggle", s.handleVentToggle)
	router.Get("/api/access/approvals", s.handleAccessApprovals)
	router.Post("/api/cabins/extend", s.handleCabinExtend)
	router.Get("/api/cabins/zones", s.handleCabinZones)
	router.Post("/api/access/revoke", s.handleAccessRevoke)
	router.Post("/api/namespaces", s.handleNamespaceOps)
	router.Get("/api/env/trend", s.handleEnvTrend)
	router.Get("/api/drain/high", s.handleDrainHigh)
	router.Post("/api/drain/threshold", s.handleDrainThreshold)
	router.Post("/api/drain/cycle", s.handleDrainCycle)
	router.Get("/api/fire/detectors", s.handleFireDetectors)
	router.Post("/api/fire/detector", s.handleFireDetectorAdd)
	router.Post("/api/fire/detect", s.handleFireDetect)
	router.Get("/api/patrol/progress", s.handlePatrolProgress)
	router.Post("/api/patrol/due", s.handlePatrolDue)
	router.Get("/api/cables/stats", s.handleCableStats)
	router.Post("/api/cables/sample", s.handleCableSample)
	router.Get("/api/audit/summary", s.handleAuditSummary)
	router.Post("/api/cabins/rename", s.handleCabinRename)
	router.Post("/api/cabins/move", s.handleCabinMove)
	router.Post("/api/cabins/zone", s.handleCabinZone)
	router.Get("/api/access/requests", s.handleAccessRequests)
	router.Post("/api/access/submit", s.handleAccessSubmit)
	router.Post("/api/access/grant", s.handleAccessGrant)
	router.Get("/api/quota/advanced", s.handleQuotaAdvanced)
	router.Get("/api/store/export", s.handleStoreExport)
	router.Post("/api/drain/levels", s.handleDrainLevels)
	router.Post("/api/audit/clear", s.handleAuditClear)
	router.Post("/api/cabins/door", s.handleCabinDoor)
	router.Post("/api/cabins/valve", s.handleCabinValve)
	router.Post("/api/cabins/retire", s.handleCabinRetire)

	return router
}

func (s *Server) servePage(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(s.webDir, name))
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
