package patrol

import "tunnelnet/internal/cabin"

type CheckIn struct {
	CabinID string `json:"cabin_id"`
	At      string `json:"at"`
}

type Service struct {
	cabins *cabin.Registry
	route  []cabin.Checkpoint
	checks []CheckIn
	due    map[string]string
}

func NewService(cabins *cabin.Registry) *Service {
	return &Service{cabins: cabins, route: cabins.Checkpoints(), checks: make([]CheckIn, 0), due: make(map[string]string)}
}

func (s *Service) Route() []cabin.Checkpoint {
	return s.route
}

func (s *Service) CheckIn(cabinID string, at string) {
	s.checks = append(s.checks, CheckIn{CabinID: cabinID, At: at})
}

func (s *Service) CheckIns() []CheckIn {
	return s.checks
}

func (s *Service) Missed() []string {
	seen := make(map[string]bool)
	for _, check := range s.checks {
		seen[check.CabinID] = true
	}
	out := make([]string, 0)
	for _, checkpoint := range s.route {
		if !seen[checkpoint.CabinID] {
			out = append(out, checkpoint.CabinID)
		}
	}
	return out
}
