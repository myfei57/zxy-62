package drain

type DrainCycle struct {
	CabinID  string `json:"cabin_id"`
	StartedAt string `json:"started_at"`
}

func (s *Service) RecordCycle(cabinID string, at string) {
	s.cycles = append(s.cycles, DrainCycle{CabinID: cabinID, StartedAt: at})
	s.state = "idle"
}

func (s *Service) Cycles() []DrainCycle {
	return append([]DrainCycle{}, s.cycles...)
}

func (s *Service) CycleCount() int {
	return len(s.cycles)
}

func (s *Service) SetThreshold(value float64) {
	s.levelThreshold = value
}

func (s *Service) HighPits() []Pit {
	out := make([]Pit, 0)
	for _, pit := range s.Pits() {
		if pit.WaterLevel >= s.levelThreshold {
			out = append(out, pit)
		}
	}
	return out
}

func (s *Service) IsHigh(id string) bool {
	level, ok := s.Level(id)
	if !ok {
		return false
	}
	return level >= s.levelThreshold
}
