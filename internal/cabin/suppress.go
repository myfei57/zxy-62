package cabin

type Suppression struct {
	Zone   string `json:"zone"`
	Active bool   `json:"active"`
}

type SuppressionLog struct {
	entries []Suppression
}

func NewSuppressionLog() *SuppressionLog {
	return &SuppressionLog{entries: make([]Suppression, 0)}
}

func (l *SuppressionLog) Activate(zone string) {
	for _, e := range l.entries {
		if e.Zone == zone {
			return
		}
	}
	l.entries = append(l.entries, Suppression{Zone: zone, Active: true})
}

func (l *SuppressionLog) Zones() []string {
	out := make([]string, 0, len(l.entries))
	for _, entry := range l.entries {
		out = append(out, entry.Zone)
	}
	return out
}
