package ns

import "sort"

type Namespace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Registry struct {
	items map[string]Namespace
}

func NewRegistry() *Registry {
	return &Registry{items: make(map[string]Namespace)}
}

func (r *Registry) Add(id string, name string) Namespace {
	item := Namespace{ID: id, Name: name}
	r.items[id] = item
	return item
}

func (r *Registry) List() []Namespace {
	ids := make([]string, 0, len(r.items))
	for id := range r.items {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]Namespace, 0, len(ids))
	for _, id := range ids {
		out = append(out, r.items[id])
	}
	return out
}

func (r *Registry) Remove(id string) bool {
	if _, ok := r.items[id]; !ok {
		return false
	}
	delete(r.items, id)
	return true
}

func (r *Registry) Exists(id string) bool {
	_, ok := r.items[id]
	return ok
}

func (r *Registry) Count() int {
	return len(r.items)
}

func (r *Registry) Names() []string {
	out := make([]string, 0, len(r.items))
	for id := range r.items {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}
