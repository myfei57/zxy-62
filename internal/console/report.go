package console

import "net/http"

type DashboardReport struct {
	Namespaces map[string]int         `json:"namespaces"`
	Cabins     map[string]int         `json:"cabins"`
	Cables     map[string]any         `json:"cables"`
	Drain      map[string]any         `json:"drain"`
	Fire       map[string]any         `json:"fire"`
	Patrol     map[string]any         `json:"patrol"`
	Quota      map[string]any         `json:"quota"`
	Audit      map[string]any         `json:"audit"`
	Access     map[string]int         `json:"access"`
}

func (s *Server) handleReport(w http.ResponseWriter, r *http.Request) {
	fireReport := s.fire.Report()
	quotaReport := s.quota.Report()
	report := DashboardReport{
		Namespaces: map[string]int{"count": s.namespaces.Count()},
		Cabins:     s.cabins.StatusSummary(),
		Cables: map[string]any{
			"count":      s.cables.Count(),
			"load_factor": s.cables.LoadFactor(),
			"danger":     s.cables.DangerCount(75),
			"healthy":    s.cables.HealthyCount(75),
		},
		Drain: map[string]any{
			"pits":          s.drain.PitCount(),
			"average_level": s.drain.AverageLevel(),
			"max_level":     s.drain.MaxLevel(),
			"high":          s.drain.HighCount(),
			"low":           s.drain.LowCount(),
			"cycles":        s.drain.CycleCount(),
		},
		Fire: map[string]any{
			"detectors": fireReport.DetectorCount,
			"alarms":    fireReport.AlarmCount,
			"zones":     fireReport.ZoneCount,
		},
		Patrol: map[string]any{
			"route":    s.patrol.RouteLen(),
			"checked":  s.patrol.CheckedCount(),
			"missed":   len(s.patrol.Missed()),
			"progress": s.patrol.Progress(),
		},
		Quota: map[string]any{
			"allocated":     quotaReport.Allocated,
			"usage_total":   quotaReport.UsageTotal,
			"overage_total": quotaReport.OverageTotal,
			"peak_cabin":    quotaReport.PeakCabin,
		},
		Audit: map[string]any{
			"count": s.audit.Count(),
			"by_event": s.audit.CountByEvent(),
		},
		Access: s.access.Summary(),
	}
	writeJSON(w, http.StatusOK, report)
}
