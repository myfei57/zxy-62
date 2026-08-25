package console

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"tunnelnet/internal/audit"
)

func (s *Server) handleCabins(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": s.cabins.All()})
}

func (s *Server) handleEnv(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"alarms": s.env.Alarms(), "route_sensor_1": s.env.RouteAlarm("sensor-1")})
}

func (s *Server) handlePatrols(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"route":       s.patrol.Route(),
		"missed":      s.patrol.Missed(),
		"route_ids":   s.patrol.RouteIDs(),
		"checked_ids": s.patrol.CheckedCabinIDs(),
		"sorted_route": s.patrol.SortedRoute(),
		"checkins":    s.patrol.CheckIns(),
	})
}

func (s *Server) handleAlarms(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"fire_alarm_count": s.fire.AlarmCount(), "env_alarms": s.env.Alarms()})
}

func (s *Server) handleNamespaces(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": s.namespaces.List(), "names": s.namespaces.Names()})
}

func (s *Server) handleQuotaReport(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": s.quota.LoadReport()})
}

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": s.audit.List()})
}

func (s *Server) handleAccessApprove(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CabinID   string `json:"cabin_id"`
		Requester string `json:"requester"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	approval, err := s.access.Approve(body.CabinID, body.Requester)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	s.recordAudit("access_approved", approval.ID)
	writeJSON(w, http.StatusOK, approval)
}

func (s *Server) handleDrainStart(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CabinID string `json:"cabin_id"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if err := s.drain.StartDrain(body.CabinID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"runs": s.drain.RunCount()})
}

func (s *Server) handleDrainExecute(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CommandKey string `json:"command_key"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if err := s.drain.Execute(body.CommandKey, func() error { return nil }); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"runs": s.drain.RunCount()})
}

func (s *Server) handleFireSmoke(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Zone string `json:"zone"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if err := s.fire.HandleSmoke(body.Zone); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	s.recordAudit("fire_smoke", body.Zone)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleFireZone(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Zone string `json:"zone"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	s.fire.ActivateZones(body.Zone)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleFireAlarm(w http.ResponseWriter, r *http.Request) {
	var body struct {
		EventKey string `json:"event_key"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	raised := s.fire.RaiseAlarm(body.EventKey)
	writeJSON(w, http.StatusOK, map[string]any{"raised": raised, "count": s.fire.AlarmCount()})
}

func (s *Server) handleQuotaRefresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Watts float64 `json:"watts"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	s.quota.RefreshBaseline(body.Watts)
	writeJSON(w, http.StatusOK, map[string]any{"items": s.quota.LoadReport()})
}

func (s *Server) handlePatrolCheckIn(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CabinID string `json:"cabin_id"`
		At      string `json:"at"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	s.patrol.CheckIn(body.CabinID, body.At)
	writeJSON(w, http.StatusOK, map[string]any{"missed": s.patrol.Missed()})
}

func (s *Server) handlePatrolRebuild(w http.ResponseWriter, r *http.Request) {
	s.patrol.RebuildRoute()
	writeJSON(w, http.StatusOK, map[string]any{"route": s.patrol.Route()})
}

func (s *Server) handleEnvPartition(w http.ResponseWriter, r *http.Request) {
	var body struct {
		PartitionID string `json:"partition_id"`
		CabinID     string `json:"cabin_id"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	s.env.Repartition(body.PartitionID, body.CabinID)
	writeJSON(w, http.StatusOK, map[string]any{"route_sensor_1": s.env.RouteAlarm("sensor-1")})
}

func (s *Server) recordAudit(event string, ref string) {
	s.audit.Append(audit.Record{
		ID:    uuid.NewString(),
		Event: event,
		At:    s.now().UTC().Format("2006-01-02T15:04:05Z"),
	})
	if ref == "" {
		return
	}
	s.audit.Append(audit.Record{
		ID:    uuid.NewString(),
		Event: event + "_ref",
		At:    s.now().UTC().Format("2006-01-02T15:04:05Z"),
	})
}

func decodeBody(r *http.Request, target any) error {
	return json.NewDecoder(r.Body).Decode(target)
}
