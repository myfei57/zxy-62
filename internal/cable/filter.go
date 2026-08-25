package cable

type Service struct {
	window []float64
	size   int
}

func NewService(size int) *Service {
	if size < 1 {
		size = 1
	}
	return &Service{window: make([]float64, 0, size), size: size}
}

func (s *Service) Filter(sample float64) float64 {
	s.window = append(s.window, sample)
	if len(s.window) > s.size {
		s.window = s.window[1:]
	}
	total := 0.0
	for _, value := range s.window {
		total += value
	}
	return total / float64(len(s.window))
}
