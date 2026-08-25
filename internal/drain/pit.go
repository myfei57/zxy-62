package drain

import "sort"

func (s *Service) AddPit(pit Pit) {
	copy := pit
	s.pits[pit.ID] = &copy
}

func (s *Service) SetLevel(id string, level float64) bool {
	pit, ok := s.pits[id]
	if !ok {
		return false
	}
	pit.WaterLevel = level
	return true
}

func (s *Service) Level(id string) (float64, bool) {
	pit, ok := s.pits[id]
	if !ok {
		return 0, false
	}
	return pit.WaterLevel, true
}

func (s *Service) Pits() []Pit {
	ids := make([]string, 0, len(s.pits))
	for id := range s.pits {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]Pit, 0, len(ids))
	for _, id := range ids {
		out = append(out, *s.pits[id])
	}
	return out
}

func (s *Service) PitCount() int {
	return len(s.pits)
}

func (s *Service) Levels() map[string]float64 {
	out := make(map[string]float64)
	for id, pit := range s.pits {
		out[id] = pit.WaterLevel
	}
	return out
}

func (s *Service) SetLevels(levels map[string]float64) {
	for id, value := range levels {
		s.SetLevel(id, value)
	}
}
