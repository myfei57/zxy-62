package env

import "sort"

func (s *Service) Variance(sensorID string) float64 {
	history := s.readings[sensorID]
	if len(history) == 0 {
		return 0
	}
	mean := s.Average(sensorID)
	total := 0.0
	for _, reading := range history {
		diff := reading.Value - mean
		total += diff * diff
	}
	return total / float64(len(history))
}

func (s *Service) StdDev(sensorID string) float64 {
	variance := s.Variance(sensorID)
	if variance <= 0 {
		return 0
	}
	x := variance
	for i := 0; i < 24; i++ {
		x = (x + variance/x) / 2
	}
	return x
}

func (s *Service) Percentile(sensorID string, p float64) float64 {
	history := s.readings[sensorID]
	if len(history) == 0 {
		return 0
	}
	values := make([]float64, 0, len(history))
	for _, reading := range history {
		values = append(values, reading.Value)
	}
	sort.Float64s(values)
	index := int(p * float64(len(values)-1))
	if index < 0 {
		index = 0
	}
	if index >= len(values) {
		index = len(values) - 1
	}
	return values[index]
}
