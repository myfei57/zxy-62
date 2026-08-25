package drain

import (
	"time"

	"tunnelnet/internal/store"
)

type Pit struct {
	ID         string  `json:"id"`
	Cabin      string  `json:"cabin"`
	WaterLevel float64 `json:"water_level"`
}

type ValveController interface {
	OpenValve(string) error
}

type PumpController interface {
	Start() error
}

type MarkStore interface {
	ExecutionMarkExists(string) bool
	SaveExecutionMark(store.ExecutionMark) error
}

type Service struct {
	valve ValveController
	pump  PumpController
	marks MarkStore
	now   func() time.Time
	runs  int
	pits  map[string]*Pit
	cycles []DrainCycle
	levelThreshold float64
	state string
	inFlight map[string]bool
}

func NewService(valve ValveController, pump PumpController, marks MarkStore) *Service {
	return &Service{
		valve:          valve,
		pump:           pump,
		marks:          marks,
		now:            time.Now,
		runs:           0,
		pits:           make(map[string]*Pit),
		cycles:         make([]DrainCycle, 0),
		levelThreshold: 0.5,
		state:          "idle",
		inFlight:       make(map[string]bool),
	}
}

func (s *Service) RunCount() int {
	return s.runs
}

func (s *Service) SetState(value string) {
	s.state = value
}

func (s *Service) State() string {
	return s.state
}

func (s *Service) MarkInFlight(key string) {
	s.inFlight[key] = true
}

func (s *Service) IsInFlight(key string) bool {
	return s.inFlight[key]
}

func (s *Service) ClearInFlight(key string) {
	delete(s.inFlight, key)
}
