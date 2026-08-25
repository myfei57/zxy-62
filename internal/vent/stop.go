package vent

func (s *Service) Stop() {
	if s.running {
		s.running = false
	}
}
