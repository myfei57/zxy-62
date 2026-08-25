package env

func (s *Service) Average(sensorID string) float64 {
	history := s.readings[sensorID]
	if len(history) == 0 {
		return 0
	}
	total := 0.0
	for _, reading := range history {
		total += reading.Value
	}
	return total / float64(len(history))
}

func (s *Service) Max(sensorID string) float64 {
	history := s.readings[sensorID]
	if len(history) == 0 {
		return 0
	}
	max := history[0].Value
	for _, reading := range history {
		if reading.Value > max {
			max = reading.Value
		}
	}
	return max
}

func (s *Service) Min(sensorID string) float64 {
	history := s.readings[sensorID]
	if len(history) == 0 {
		return 0
	}
	min := history[0].Value
	for _, reading := range history {
		if reading.Value < min {
			min = reading.Value
		}
	}
	return min
}

func (s *Service) Trend(sensorID string) float64 {
	history := s.readings[sensorID]
	if len(history) < 2 {
		return 0
	}
	return history[len(history)-1].Value - history[0].Value
}

func (s *Service) ReadingCount(sensorID string) int {
	return len(s.readings[sensorID])
}
