package vent

func (s *Service) Stop() {
	s.running = false
}
