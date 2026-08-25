package cabin

type PowerBaseline struct {
	CabinID string  `json:"cabin_id"`
	Watts   float64 `json:"watts"`
}

type PowerBook struct {
	baselines map[string]float64
}

func NewPowerBook() *PowerBook {
	return &PowerBook{baselines: make(map[string]float64)}
}

func (b *PowerBook) SetBaseline(cabinID string, watts float64) {
	b.baselines[cabinID] = watts
}

func (b *PowerBook) Baseline(cabinID string) float64 {
	return b.baselines[cabinID]
}

func (b *PowerBook) Total() float64 {
	total := 0.0
	for _, watts := range b.baselines {
		total += watts
	}
	return total
}
