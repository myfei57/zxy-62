package fire

func (s *Service) HandleSmoke(zone string) error {
	// 防火阀关闭必须先落盘成功，断电重启后才能按分区恢复关闭状态，
	// 避免烟气扩散到相邻舱室。写盘失败时向上抛出 error，不再继续联动。
	if err := s.damper.Save("damper", zone, DamperState{Zone: zone, Closed: true}); err != nil {
		return err
	}
	// 落盘成功后再停通风、抑制分区，保证联动顺序可恢复。
	s.vent.Stop()
	s.suppressor.Activate(zone)
	return nil
}
