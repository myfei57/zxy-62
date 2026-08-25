package cabin

type Checkpoint struct {
	CabinID string `json:"cabin_id"`
	Order   int    `json:"order"`
}

func (r *Registry) Checkpoints() []Checkpoint {
	active := r.Active()
	out := make([]Checkpoint, 0, len(active))
	for index, item := range active {
		out = append(out, Checkpoint{CabinID: item.ID, Order: index + item.Revision - 1})
	}
	return out
}
