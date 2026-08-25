package quota

func (s *Service) Overage(cabinID string) float64 {
	usage := s.usage[cabinID]
	baseline := s.book.Baseline(cabinID)
	if usage <= baseline {
		return 0
	}
	return usage - baseline
}

func (s *Service) Allocated() float64 {
	return s.book.Total()
}

func (s *Service) OverageTotal() float64 {
	total := 0.0
	for cabinID := range s.usage {
		total += s.Overage(cabinID)
	}
	return total
}

func (s *Service) Forecast(cabinID string, factor float64) float64 {
	return s.usage[cabinID] * factor
}

func (s *Service) PeakCabin() string {
	best := ""
	bestValue := -1.0
	for cabinID, value := range s.usage {
		if value > bestValue {
			best = cabinID
			bestValue = value
		}
	}
	return best
}
