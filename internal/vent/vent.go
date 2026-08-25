package vent

type Service struct {
	running bool
}

func NewService() *Service {
	return &Service{running: true}
}

func (s *Service) Start() {
	s.running = true
}

func (s *Service) Running() bool {
	return s.running
}
