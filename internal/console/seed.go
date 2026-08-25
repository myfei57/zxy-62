package console

func (s *Server) seedDemo() {
	s.env.Sample("sensor-1", 48.0, "2026-08-25T08:00:00Z")
	s.env.Sample("sensor-1", 51.0, "2026-08-25T08:05:00Z")
	s.env.Sample("sensor-2", 55.0, "2026-08-25T08:05:00Z")
	s.patrol.CheckIn("cabin-1", "2026-08-25T08:00:00Z")
	s.patrol.SetDue("cabin-2", "2026-08-25T07:00:00Z")
	s.quota.AddUsage("cabin-1", 1200)
	s.quota.AddUsage("cabin-2", 900)
	s.drain.RecordCycle("cabin-1", "2026-08-25T08:00:00Z")
	s.drain.RecordCycle("cabin-1", "2026-08-25T08:10:00Z")
	s.cables.Sample("cable-1", 42.0)
	s.cables.Sample("cable-1", 47.0)
}
