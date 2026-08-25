package fire

func (s *Service) RaiseAlarm(eventKey string) bool {
	if !s.deduper.Record(eventKey) {
		return false
	}
	s.alarms = append(s.alarms, eventKey)
	return true
}
