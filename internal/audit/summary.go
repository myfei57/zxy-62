package audit

func (r *Registry) CountByEvent() map[string]int {
	out := make(map[string]int)
	for _, record := range r.records {
		out[record.Event]++
	}
	return out
}

func (r *Registry) Summary() map[string]int {
	out := make(map[string]int)
	for _, record := range r.records {
		out[record.Event]++
	}
	out["total"] = len(r.records)
	return out
}

func (r *Registry) Between(start string, end string) []Record {
	out := make([]Record, 0)
	for _, record := range r.records {
		if record.At >= start && record.At <= end {
			out = append(out, record)
		}
	}
	return out
}
