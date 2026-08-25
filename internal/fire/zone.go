package fire

import "sort"

func (s *Service) ActivateZones(source string) {
	s.activate(source)
	for _, zone := range s.adjacent(source) {
		s.activate(zone)
	}
}

func (s *Service) activate(zone string) {
	s.suppressor.Activate(zone)
}

func (s *Service) adjacent(source string) []string {
	type ranked struct {
		zone     string
		distance int
	}
	index := -1
	for i, zone := range s.zones {
		if zone == source {
			index = i
			break
		}
	}
	items := make([]ranked, 0)
	for i, zone := range s.zones {
		if zone == source {
			continue
		}
		distance := i - index
		if distance < 0 {
			distance = -distance
		}
		items = append(items, ranked{zone: zone, distance: distance})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].distance != items[j].distance {
			return items[i].distance < items[j].distance
		}
		return items[i].zone < items[j].zone
	})
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, item.zone)
	}
	return out
}
