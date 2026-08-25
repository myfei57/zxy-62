package cabin

import "sort"

func (r *Registry) Count() int {
	return len(r.items)
}

func (r *Registry) ActiveCount() int {
	count := 0
	for _, item := range r.items {
		if !item.Retired {
			count++
		}
	}
	return count
}

func (r *Registry) ZoneCabins(zone string) []Cabin {
	out := make([]Cabin, 0)
	for _, item := range r.items {
		if item.Zone == zone && !item.Retired {
			out = append(out, *item)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (r *Registry) Zones() []string {
	seen := make(map[string]bool)
	out := make([]string, 0)
	for _, item := range r.items {
		if seen[item.Zone] {
			continue
		}
		seen[item.Zone] = true
		out = append(out, item.Zone)
	}
	sort.Strings(out)
	return out
}

func (r *Registry) Extend(item Cabin) error {
	item.Revision = 0
	return r.Upsert(item)
}
