package fire

type DamperStore interface {
	Save(string, string, any) error
}

type VentStopper interface {
	Stop()
}

type Suppressor interface {
	Activate(string)
}

type Deduper interface {
	Record(string) bool
}

type DamperState struct {
	Zone   string `json:"zone"`
	Closed bool   `json:"closed"`
}

type Service struct {
	damper     DamperStore
	vent       VentStopper
	suppressor Suppressor
	deduper    Deduper
	zones      []string
	alarms     []string
	detectors  map[string]Detector
}

func NewService(damper DamperStore, vent VentStopper, suppressor Suppressor, deduper Deduper, zones []string) *Service {
	return &Service{
		damper:     damper,
		vent:       vent,
		suppressor: suppressor,
		deduper:    deduper,
		zones:      append([]string{}, zones...),
		alarms:     make([]string, 0),
		detectors:  make(map[string]Detector),
	}
}

func (s *Service) AlarmCount() int {
	return len(s.alarms)
}
