package console

import (
	"net/http"
	"strconv"

	"tunnelnet/internal/fire"
)

func (s *Server) handleEnvTrend(w http.ResponseWriter, r *http.Request) {
	sensorID := r.URL.Query().Get("sensor")
	if sensorID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "sensor required"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"sensor":    sensorID,
		"count":     s.env.ReadingCount(sensorID),
		"average":   s.env.Average(sensorID),
		"max":       s.env.Max(sensorID),
		"min":       s.env.Min(sensorID),
		"trend":     s.env.Trend(sensorID),
		"variance":  s.env.Variance(sensorID),
		"stddev":    s.env.StdDev(sensorID),
		"percentile": s.env.Percentile(sensorID, 0.95),
	})
}

func (s *Server) handleDrainHigh(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"high":   s.drain.HighPits(),
		"cycles": s.drain.Cycles(),
		"count":  s.drain.CycleCount(),
	})
}

func (s *Server) handleDrainThreshold(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Threshold float64 `json:"threshold"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	s.drain.SetThreshold(body.Threshold)
	writeJSON(w, http.StatusOK, map[string]any{"high": s.drain.HighPits()})
}

func (s *Server) handleDrainCycle(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CabinID string `json:"cabin_id"`
		At      string `json:"at"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	s.drain.RecordCycle(body.CabinID, body.At)
	writeJSON(w, http.StatusOK, map[string]any{"count": s.drain.CycleCount()})
}

func (s *Server) handleFireDetectors(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": s.fire.Detectors(), "count": s.fire.DetectorCount()})
}

func (s *Server) handleFireDetectorAdd(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID   string `json:"id"`
		Zone string `json:"zone"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	s.fire.AddDetector(fire.Detector{ID: body.ID, Zone: body.Zone})
	writeJSON(w, http.StatusOK, map[string]any{"items": s.fire.Detectors()})
}

func (s *Server) handleFireDetect(w http.ResponseWriter, r *http.Request) {
	var body struct {
		DetectorID string `json:"detector_id"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if err := s.fire.Detect(body.DetectorID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handlePatrolProgress(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	writeJSON(w, http.StatusOK, map[string]any{
		"progress": s.patrol.Progress(),
		"overdue":  s.patrol.Overdue(now),
	})
}

func (s *Server) handlePatrolDue(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CabinID string `json:"cabin_id"`
		Due     string `json:"due"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	s.patrol.SetDue(body.CabinID, body.Due)
	writeJSON(w, http.StatusOK, map[string]any{"due": s.patrol.Due(body.CabinID)})
}

func (s *Server) handleCableStats(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "id required"})
		return
	}
	threshold := 75.0
	if value := r.URL.Query().Get("threshold"); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			threshold = parsed
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id":       id,
		"history":  s.cables.History(id),
		"average":  s.cables.Average(id),
		"max":      s.cables.Max(id),
		"min":      s.cables.Min(id),
		"over":     s.cables.Over(id, threshold),
	})
}

func (s *Server) handleCableSample(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CableID string  `json:"cable_id"`
		Value   float64 `json:"value"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if !s.cables.Sample(body.CableID, body.Value) {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "cable not found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"history": s.cables.History(body.CableID)})
}

func (s *Server) handleAuditSummary(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	payload := map[string]any{
		"summary":  s.audit.Summary(),
		"by_event": s.audit.CountByEvent(),
	}
	if start != "" && end != "" {
		payload["between"] = s.audit.Between(start, end)
	}
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handleCabinRename(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if err := s.cabins.Rename(body.ID, body.Name); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"cabin": s.cabins.Find(body.ID)})
}

func (s *Server) handleCabinMove(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID        string `json:"id"`
		Zone      string `json:"zone"`
		Partition string `json:"partition"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if err := s.cabins.MoveCabin(body.ID, body.Zone, body.Partition); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"cabin": s.cabins.Find(body.ID)})
}

func (s *Server) handleCabinZone(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID   string `json:"id"`
		Zone string `json:"zone"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if err := s.cabins.AssignZone(body.ID, body.Zone); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"cabin": s.cabins.Find(body.ID)})
}

func (s *Server) handleAccessRequests(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"pending": s.access.Pending(), "count": s.access.RequestCount()})
}

func (s *Server) handleAccessSubmit(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CabinID   string `json:"cabin_id"`
		Requester string `json:"requester"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	request := s.access.Submit(body.CabinID, body.Requester)
	writeJSON(w, http.StatusOK, request)
}

func (s *Server) handleAccessGrant(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RequestID string `json:"request_id"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	approval, err := s.access.Grant(body.RequestID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, approval)
}

func (s *Server) handleQuotaAdvanced(w http.ResponseWriter, r *http.Request) {
	factor := 1.0
	if value := r.URL.Query().Get("factor"); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			factor = parsed
		}
	}
	cabinID := r.URL.Query().Get("cabin")
	writeJSON(w, http.StatusOK, map[string]any{
		"allocated":   s.quota.Allocated(),
		"overage_total": s.quota.OverageTotal(),
		"peak_cabin":  s.quota.PeakCabin(),
		"forecast":    s.quota.Forecast(cabinID, factor),
		"overage":     s.quota.Overage(cabinID),
	})
}

func (s *Server) handleCabinDoor(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID      string `json:"id"`
		Release bool   `json:"release"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	var err error
	if body.Release {
		err = s.cabins.ReleaseDoor(body.ID)
	} else {
		err = s.cabins.SealDoor(body.ID)
	}
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	state, _ := s.cabins.DoorState(body.ID)
	writeJSON(w, http.StatusOK, map[string]any{"id": body.ID, "open": state})
}

func (s *Server) handleCabinValve(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID   string `json:"id"`
		Open bool   `json:"open"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	var err error
	if body.Open {
		err = s.cabins.OpenValve(body.ID)
	} else {
		err = s.cabins.CloseValve(body.ID)
	}
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	state, _ := s.cabins.ValveState(body.ID)
	writeJSON(w, http.StatusOK, map[string]any{"id": body.ID, "open": state})
}

func (s *Server) handleCabinRetire(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID string `json:"id"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if err := s.cabins.Retire(body.ID); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"active": s.cabins.Active()})
}
