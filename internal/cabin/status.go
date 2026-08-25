package cabin

func (r *Registry) StatusSummary() map[string]int {
	openDoors := 0
	openValves := 0
	for _, item := range r.items {
		if item.DoorOpen {
			openDoors++
		}
		if item.ValveOpen {
			openValves++
		}
	}
	return map[string]int{
		"total":       len(r.items),
		"active":      r.ActiveCount(),
		"open_doors":  openDoors,
		"open_valves": openValves,
	}
}
