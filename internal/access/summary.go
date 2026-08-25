package access

func (s *Service) Summary() map[string]int {
	pending := 0
	granted := 0
	for _, request := range s.requests {
		if request.Status == "pending" {
			pending++
		}
		if request.Status == "granted" {
			granted++
		}
	}
	return map[string]int{"requests": len(s.requests), "pending": pending, "granted": granted}
}
