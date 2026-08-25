package cable

func (r *Registry) Sample(id string, value float64) bool {
	if _, ok := r.items[id]; !ok {
		return false
	}
	r.history[id] = append(r.history[id], value)
	r.SetTemp(id, value)
	return true
}

func (r *Registry) History(id string) []float64 {
	return append([]float64{}, r.history[id]...)
}

func (r *Registry) Average(id string) float64 {
	values := r.history[id]
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value
	}
	return total / float64(len(values))
}

func (r *Registry) Max(id string) float64 {
	values := r.history[id]
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, value := range values {
		if value > max {
			max = value
		}
	}
	return max
}

func (r *Registry) Min(id string) float64 {
	values := r.history[id]
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, value := range values {
		if value < min {
			min = value
		}
	}
	return min
}

func (r *Registry) Over(id string, threshold float64) int {
	count := 0
	for _, value := range r.history[id] {
		if value >= threshold {
			count++
		}
	}
	return count
}

func (r *Registry) SampleCount() int {
	total := 0
	for _, values := range r.history {
		total += len(values)
	}
	return total
}

func (r *Registry) Summary() map[string]int {
	return map[string]int{"total": len(r.items), "samples": r.SampleCount()}
}
