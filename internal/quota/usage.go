package quota

import "sort"

func (s *Service) AddUsage(cabinID string, watts float64) {
	s.usage[cabinID] += watts
}

func (s *Service) Usage(cabinID string) float64 {
	return s.usage[cabinID]
}

func (s *Service) TotalUsage() float64 {
	total := 0.0
	for _, value := range s.usage {
		total += value
	}
	return total
}

func (s *Service) Peak() float64 {
	peak := 0.0
	for _, value := range s.usage {
		if value > peak {
			peak = value
		}
	}
	return peak
}

func (s *Service) UsageByCabin() map[string]float64 {
	out := make(map[string]float64)
	for id, value := range s.usage {
		out[id] = value
	}
	return out
}

func (s *Service) RankedCabinIDs() []string {
	ids := make([]string, 0, len(s.usage))
	for id := range s.usage {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return s.usage[ids[i]] > s.usage[ids[j]] })
	return ids
}
