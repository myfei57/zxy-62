package patrol

func (s *Service) RebuildRoute() {
	if len(s.route) > 0 {
		return
	}
	s.route = s.cabins.Checkpoints()
}

func (s *Service) RouteIDs() []string {
	out := make([]string, 0, len(s.route))
	for _, checkpoint := range s.route {
		out = append(out, checkpoint.CabinID)
	}
	return out
}

func (s *Service) CheckedCabinIDs() []string {
	out := make([]string, 0, len(s.checks))
	for _, check := range s.checks {
		out = append(out, check.CabinID)
	}
	return out
}
