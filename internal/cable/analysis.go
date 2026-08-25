package cable

func (r *Registry) DangerCount(threshold float64) int {
	count := 0
	for _, item := range r.items {
		if item.Temp >= threshold {
			count++
		}
	}
	return count
}

func (r *Registry) HealthyCount(threshold float64) int {
	count := 0
	for _, item := range r.items {
		if item.Temp < threshold {
			count++
		}
	}
	return count
}

func (r *Registry) LoadFactor() float64 {
	if len(r.items) == 0 {
		return 0
	}
	total := 0.0
	for _, item := range r.items {
		total += item.Temp
	}
	return total / float64(len(r.items))
}
