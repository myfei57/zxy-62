package cable

import "sort"

type Registry struct {
	items   map[string]*Cable
	history map[string][]float64
}

func NewRegistry() *Registry {
	return &Registry{items: make(map[string]*Cable), history: make(map[string][]float64)}
}

func (r *Registry) Add(item Cable) {
	copy := item
	r.items[item.ID] = &copy
}

func (r *Registry) Get(id string) (Cable, bool) {
	item, ok := r.items[id]
	if !ok {
		return Cable{}, false
	}
	return *item, true
}

func (r *Registry) SetTemp(id string, temp float64) bool {
	item, ok := r.items[id]
	if !ok {
		return false
	}
	item.Temp = temp
	if temp >= 75 {
		item.Status = "over"
	} else {
		item.Status = "normal"
	}
	return true
}

func (r *Registry) List() []Cable {
	ids := make([]string, 0, len(r.items))
	for id := range r.items {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]Cable, 0, len(ids))
	for _, id := range ids {
		out = append(out, *r.items[id])
	}
	return out
}

func (r *Registry) Count() int {
	return len(r.items)
}
