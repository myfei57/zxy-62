package fire

func (s *Service) HandleSmoke(zone string) error {
	if err := s.damper.Save("damper", zone, DamperState{Zone: zone, Closed: true}); err != nil {
		return err
	}
	s.vent.Stop()
	return nil
}
