package quota

import "tunnelnet/internal/cabin"

type Service struct {
	cabins *cabin.Registry
	book   *cabin.PowerBook
	usage  map[string]float64
}

func NewService(cabins *cabin.Registry, book *cabin.PowerBook) *Service {
	return &Service{cabins: cabins, book: book, usage: make(map[string]float64)}
}
