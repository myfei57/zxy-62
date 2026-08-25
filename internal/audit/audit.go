package audit

import "sort"

type Record struct {
	ID    string `json:"id"`
	Event string `json:"event"`
	At    string `json:"at"`
}

type Registry struct {
	records []Record
}

func NewRegistry() *Registry {
	return &Registry{records: make([]Record, 0)}
}

func (r *Registry) Append(record Record) {
	r.records = append(r.records, record)
}

func (r *Registry) List() []Record {
	out := make([]Record, len(r.records))
	copy(out, r.records)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (r *Registry) Count() int {
	return len(r.records)
}

func (r *Registry) Clear() {
	r.records = make([]Record, 0)
}

func (r *Registry) LatestID() string {
	if len(r.records) == 0 {
		return ""
	}
	return r.records[len(r.records)-1].ID
}
