package patrol

import "sort"

func (s *Service) RouteLen() int {
	return len(s.route)
}

func (s *Service) CheckedCount() int {
	return len(s.checks)
}

func (s *Service) Summary() map[string]int {
	return map[string]int{
		"route":   len(s.route),
		"checked": len(s.checks),
		"missed":  len(s.Missed()),
	}
}

func (s *Service) SortedRoute() []string {
	ids := append([]string{}, s.RouteIDs()...)
	sort.Strings(ids)
	return ids
}
