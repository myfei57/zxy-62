package cabin

func (r *Registry) Find(id string) Cabin {
	item, ok := r.items[id]
	if !ok {
		return Cabin{}
	}
	return *item
}

func (r *Registry) Has(id string) bool {
	_, ok := r.items[id]
	return ok
}
