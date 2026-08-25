package fire

func (s *Service) Alarms() []string {
	return append([]string{}, s.alarms...)
}

func (s *Service) ClearAlarms() {
	s.alarms = make([]string, 0)
}

func (s *Service) DistinctAlarms() int {
	seen := make(map[string]bool)
	for _, alarm := range s.alarms {
		seen[alarm] = true
	}
	return len(seen)
}
