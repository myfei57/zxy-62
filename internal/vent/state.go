package vent

func (s *Service) Toggle() bool {
	s.running = !s.running
	return s.running
}

func (s *Service) Status() string {
	if s.running {
		return "running"
	}
	return "stopped"
}
