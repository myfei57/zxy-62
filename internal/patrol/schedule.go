package patrol

func (s *Service) SetDue(cabinID string, due string) {
	s.due[cabinID] = due
}

func (s *Service) Due(cabinID string) string {
	return s.due[cabinID]
}

func (s *Service) Overdue(now string) []string {
	checked := make(map[string]bool)
	for _, check := range s.checks {
		checked[check.CabinID] = true
	}
	out := make([]string, 0)
	for _, checkpoint := range s.route {
		due := s.due[checkpoint.CabinID]
		if due == "" {
			continue
		}
		if !checked[checkpoint.CabinID] && due < now {
			out = append(out, checkpoint.CabinID)
		}
	}
	return out
}

func (s *Service) Progress() float64 {
	if len(s.route) == 0 {
		return 0
	}
	return float64(len(s.checks)) / float64(len(s.route))
}
