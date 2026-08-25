package audit

func (r *Registry) Filter(event string) []Record {
	out := make([]Record, 0)
	for _, record := range r.records {
		if record.Event == event {
			out = append(out, record)
		}
	}
	return out
}

func (r *Registry) Last(n int) []Record {
	if n <= 0 {
		return []Record{}
	}
	if n > len(r.records) {
		n = len(r.records)
	}
	start := len(r.records) - n
	out := make([]Record, n)
	copy(out, r.records[start:])
	return out
}

func (r *Registry) Events() []string {
	seen := make(map[string]bool)
	out := make([]string, 0)
	for _, record := range r.records {
		if seen[record.Event] {
			continue
		}
		seen[record.Event] = true
		out = append(out, record.Event)
	}
	return out
}

func (r *Registry) MostRecent() (Record, bool) {
	if len(r.records) == 0 {
		return Record{}, false
	}
	return r.records[len(r.records)-1], true
}
