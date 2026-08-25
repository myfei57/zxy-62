package env

import "tunnelnet/internal/cabin"

type Filterer interface {
	Filter(float64) float64
}

type Alarm struct {
	SensorID string  `json:"sensor_id"`
	Value    float64 `json:"value"`
}

type Service struct {
	filter         Filterer
	partition      *cabin.PartitionTable
	sensorPartition map[string]string
	readings       map[string][]Reading
	sensorOrder    []string
	alarmThreshold float64
	alarms         []Alarm
}

func NewService(filter Filterer, partition *cabin.PartitionTable, threshold float64) *Service {
	return &Service{
		filter:          filter,
		partition:       partition,
		sensorPartition: make(map[string]string),
		readings:        make(map[string][]Reading),
		sensorOrder:     make([]string, 0),
		alarmThreshold:  threshold,
		alarms:          make([]Alarm, 0),
	}
}

func (s *Service) Alarms() []Alarm {
	return s.alarms
}
