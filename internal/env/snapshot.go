package env

func (s *Service) Snapshot() map[string]any {
	sensors := make(map[string]any)
	for _, sensorID := range s.sensorOrder {
		history := s.readings[sensorID]
		latest := 0.0
		if len(history) > 0 {
			latest = history[len(history)-1].Value
		}
		sensors[sensorID] = map[string]any{
			"partition": s.sensorPartition[sensorID],
			"latest":    latest,
			"count":     len(history),
		}
	}
	return map[string]any{
		"sensors": sensors,
		"alarms":  len(s.alarms),
	}
}
