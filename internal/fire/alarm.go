package fire

func (s *Service) Alarms() []string {
	return append([]string{}, s.alarms...)
}

func (s *Service) ClearAlarms() {
	s.alarms = make([]string, 0)
}

func (s *Service) DistinctAlarms() int {
	return len(s.alarms)
}
