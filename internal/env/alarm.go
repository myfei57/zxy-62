package env

func (s *Service) TemperatureAlarm(sensorID string, raw float64) bool {
	filtered := s.filter.Filter(raw)
	if filtered < s.alarmThreshold {
		return false
	}
	s.alarms = append(s.alarms, Alarm{SensorID: sensorID, Value: filtered})
	return true
}
