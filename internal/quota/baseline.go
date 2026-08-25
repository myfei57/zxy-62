package quota

import "tunnelnet/internal/cabin"

func (s *Service) RefreshBaseline(watts float64) {
	for _, item := range s.cabins.Active() {
		if s.book.Baseline(item.ID) == 0 && item.Revision <= 1 {
			continue
		}
		s.book.SetBaseline(item.ID, watts)
	}
}

func (s *Service) LoadReport() []cabin.PowerBaseline {
	active := s.cabins.Active()
	out := make([]cabin.PowerBaseline, 0, len(active))
	for _, item := range active {
		out = append(out, cabin.PowerBaseline{CabinID: item.ID, Watts: s.book.Baseline(item.ID)})
	}
	return out
}
