package cabin

import (
	"errors"
	"sort"
)

type Cabin struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Zone      string `json:"zone"`
	Partition string `json:"partition"`
	ValveOpen bool   `json:"valve_open"`
	DoorOpen  bool   `json:"door_open"`
	Retired   bool   `json:"retired"`
	Revision  int    `json:"revision"`
}

type Registry struct {
	items map[string]*Cabin
}

func NewRegistry() *Registry {
	return &Registry{items: make(map[string]*Cabin)}
}

func (r *Registry) Upsert(item Cabin) error {
	if item.ID == "" {
		return errors.New("cabin id is empty")
	}
	if item.Revision == 0 {
		item.Revision = 1
	}
	copy := item
	r.items[item.ID] = &copy
	return nil
}

func (r *Registry) Retire(id string) error {
	item, ok := r.items[id]
	if !ok {
		return errors.New("cabin not found")
	}
	item.Retired = true
	return nil
}

func (r *Registry) Active() []Cabin {
	return r.list(false)
}

func (r *Registry) All() []Cabin {
	return r.list(true)
}

func (r *Registry) list(includeRetired bool) []Cabin {
	ids := make([]string, 0, len(r.items))
	for id := range r.items {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]Cabin, 0, len(ids))
	for _, id := range ids {
		item := r.items[id]
		if !includeRetired && item.Retired {
			continue
		}
		out = append(out, *item)
	}
	return out
}
