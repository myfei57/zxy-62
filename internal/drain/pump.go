package drain

func (s *Service) StartDrain(cabinID string) error {
	if err := s.valve.OpenValve(cabinID); err != nil {
		return err
	}
	s.state = "valve_open"
	if err := s.pump.Start(); err != nil {
		return err
	}
	s.state = "pump_running"
	s.runs++
	return nil
}
