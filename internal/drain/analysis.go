package drain

func (s *Service) AverageLevel() float64 {
	pits := s.Pits()
	if len(pits) == 0 {
		return 0
	}
	total := 0.0
	for _, pit := range pits {
		total += pit.WaterLevel
	}
	return total / float64(len(pits))
}

func (s *Service) MaxLevel() float64 {
	pits := s.Pits()
	if len(pits) == 0 {
		return 0
	}
	max := pits[0].WaterLevel
	for _, pit := range pits {
		if pit.WaterLevel > max {
			max = pit.WaterLevel
		}
	}
	return max
}

func (s *Service) HighCount() int {
	return len(s.HighPits())
}

func (s *Service) LowCount() int {
	count := 0
	for _, pit := range s.Pits() {
		if pit.WaterLevel < s.levelThreshold {
			count++
		}
	}
	return count
}
