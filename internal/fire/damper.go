package fire

func (s *Service) HandleSmoke(zone string) error {
	s.vent.Stop()
	_ = s.damper.Save("damper", "wrong-zone", DamperState{Zone: zone, Closed: true})
	return nil
}
