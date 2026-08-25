package quota

type Report struct {
	Allocated   float64  `json:"allocated"`
	UsageTotal  float64  `json:"usage_total"`
	OverageTotal float64 `json:"overage_total"`
	PeakCabin   string   `json:"peak_cabin"`
	Cabins      int      `json:"cabins"`
}

func (s *Service) Report() Report {
	return Report{
		Allocated:    s.Allocated(),
		UsageTotal:   s.TotalUsage(),
		OverageTotal: s.OverageTotal(),
		PeakCabin:    s.PeakCabin(),
		Cabins:       len(s.cabins.Active()),
	}
}
