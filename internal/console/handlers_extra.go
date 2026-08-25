package console

import (
	"net/http"
	"strconv"

	"tunnelnet/internal/cabin"
	"tunnelnet/internal/drain"
)

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"namespaces":          s.namespaces.Count(),
		"cabins":              s.cabins.Count(),
		"active_cabins":       s.cabins.ActiveCount(),
		"cables":              s.cables.Count(),
		"pits":                s.drain.PitCount(),
		"power_baseline_total": s.powerBook.Total(),
		"power_usage_total":   s.quota.TotalUsage(),
		"fire_alarms":         s.fire.AlarmCount(),
		"patrol":              s.patrol.Summary(),
		"audit":               s.audit.Count(),
		"vent":                s.ventilation.Status(),
		"vent_running":        s.ventilation.Running(),
		"suppressed_zones":    s.suppression.Zones(),
		"store":               s.storage.Stats(),
		"files":               s.storage.FileCount(),
		"bytes":               s.storage.TotalBytes(),
	})
}

func (s *Server) handleCables(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": s.cables.List(), "summary": s.cables.Summary()})
}

func (s *Server) handleCableTemp(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CableID string  `json:"cable_id"`
		Temp    float64 `json:"temp"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if !s.cables.SetTemp(body.CableID, body.Temp) {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "cable not found"})
		return
	}
	item, _ := s.cables.Get(body.CableID)
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) handleEnvSensors(w http.ResponseWriter, r *http.Request) {
	sensorID := r.URL.Query().Get("sensor")
	payload := map[string]any{"sensors": s.env.Sensors(), "snapshot": s.env.Snapshot()}
	if sensorID != "" {
		payload["history"] = s.env.History(sensorID)
	}
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handleEnvSample(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SensorID string  `json:"sensor_id"`
		Value    float64 `json:"value"`
		At       string  `json:"at"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	alarmed := s.env.Sample(body.SensorID, body.Value, body.At)
	writeJSON(w, http.StatusOK, map[string]any{"alarmed": alarmed})
}

func (s *Server) handleDrainPits(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	payload := map[string]any{"items": s.drain.Pits(), "count": s.drain.PitCount(), "levels": s.drain.Levels()}
	if id != "" {
		payload["is_high"] = s.drain.IsHigh(id)
	}
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handleDrainPit(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID         string  `json:"id"`
		Cabin      string  `json:"cabin"`
		WaterLevel float64 `json:"water_level"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	s.drain.AddPit(drain.Pit{ID: body.ID, Cabin: body.Cabin, WaterLevel: body.WaterLevel})
	writeJSON(w, http.StatusOK, map[string]any{"items": s.drain.Pits()})
}

func (s *Server) handleQuotaUsage(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"usage":     s.quota.UsageByCabin(),
		"total":     s.quota.TotalUsage(),
		"peak":      s.quota.Peak(),
		"ranked":    s.quota.RankedCabinIDs(),
		"baseline":  s.quota.LoadReport(),
	})
}

func (s *Server) handleQuotaAddUsage(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CabinID string  `json:"cabin_id"`
		Watts   float64 `json:"watts"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	s.quota.AddUsage(body.CabinID, body.Watts)
	writeJSON(w, http.StatusOK, map[string]any{"usage": s.quota.Usage(body.CabinID)})
}

func (s *Server) handleAuditQuery(w http.ResponseWriter, r *http.Request) {
	event := r.URL.Query().Get("event")
	last := r.URL.Query().Get("last")
	payload := map[string]any{
		"count":  s.audit.Count(),
		"events": s.audit.Events(),
	}
	if event != "" {
		payload["filtered"] = s.audit.Filter(event)
	}
	if last != "" {
		if n, err := strconv.Atoi(last); err == nil {
			payload["last"] = s.audit.Last(n)
		}
	}
	if recent, ok := s.audit.MostRecent(); ok {
		payload["most_recent"] = recent
	}
	payload["latest_id"] = s.audit.LatestID()
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handleDrainLevels(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Levels map[string]float64 `json:"levels"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	s.drain.SetLevels(body.Levels)
	writeJSON(w, http.StatusOK, map[string]any{"levels": s.drain.Levels()})
}

func (s *Server) handleAuditClear(w http.ResponseWriter, r *http.Request) {
	s.audit.Clear()
	writeJSON(w, http.StatusOK, map[string]any{"count": s.audit.Count()})
}

func (s *Server) handleStoreExport(w http.ResponseWriter, r *http.Request) {
	kind := r.URL.Query().Get("kind")
	if kind == "" {
		all, err := s.storage.ExportAll()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": all})
		return
	}
	records, err := s.storage.Export(kind)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": records})
}

func (s *Server) handleFireAlarms(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": s.fire.Alarms(), "count": s.fire.AlarmCount(), "distinct": s.fire.DistinctAlarms()})
}

func (s *Server) handleFireClear(w http.ResponseWriter, r *http.Request) {
	s.fire.ClearAlarms()
	writeJSON(w, http.StatusOK, map[string]any{"count": s.fire.AlarmCount()})
}

func (s *Server) handleVentToggle(w http.ResponseWriter, r *http.Request) {
	running := s.ventilation.Toggle()
	writeJSON(w, http.StatusOK, map[string]any{"running": running, "status": s.ventilation.Status()})
}

func (s *Server) handleAccessApprovals(w http.ResponseWriter, r *http.Request) {
	approvals, err := s.storage.ListApprovals()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	out := make([]map[string]any, 0, len(approvals))
	for _, approval := range approvals {
		out = append(out, map[string]any{
			"id":        approval.ID,
			"cabin_id":  approval.CabinID,
			"requester": approval.Requester,
			"durable":   s.access.Approved(approval.ID),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": out})
}

func (s *Server) handleAccessRevoke(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ApprovalID string `json:"approval_id"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if err := s.storage.Delete("approval", body.ApprovalID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleCabinExtend(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Zone      string `json:"zone"`
		Partition string `json:"partition"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if err := s.cabins.Extend(cabin.Cabin{ID: body.ID, Name: body.Name, Zone: body.Zone, Partition: body.Partition}); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": s.cabins.All()})
}

func (s *Server) handleCabinZones(w http.ResponseWriter, r *http.Request) {
	zone := r.URL.Query().Get("zone")
	id := r.URL.Query().Get("id")
	payload := map[string]any{"zones": s.cabins.Zones()}
	if zone != "" {
		payload["cabins"] = s.cabins.ZoneCabins(zone)
	}
	if id != "" {
		payload["found"] = s.cabins.Has(id)
		payload["cabin"] = s.cabins.Find(id)
	}
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handleNamespaceOps(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if body.ID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "id required"})
		return
	}
	if s.namespaces.Exists(body.ID) {
		s.namespaces.Remove(body.ID)
		writeJSON(w, http.StatusOK, map[string]any{"removed": true, "count": s.namespaces.Count()})
		return
	}
	s.namespaces.Add(body.ID, body.Name)
	writeJSON(w, http.StatusOK, map[string]any{"added": true, "count": s.namespaces.Count()})
}
