package env

import "sort"

type Reading struct {
	SensorID string  `json:"sensor_id"`
	Value    float64 `json:"value"`
	At       string  `json:"at"`
}

func (s *Service) BindSensorAt(sensorID string, partitionID string) {
	s.sensorPartition[sensorID] = partitionID
	if !containsString(s.sensorOrder, sensorID) {
		s.sensorOrder = append(s.sensorOrder, sensorID)
		sort.Strings(s.sensorOrder)
	}
}

func (s *Service) Sample(sensorID string, value float64, at string) bool {
	s.readings[sensorID] = append(s.readings[sensorID], Reading{SensorID: sensorID, Value: value, At: at})
	return s.TemperatureAlarm(sensorID, value)
}

func (s *Service) History(sensorID string) []Reading {
	return append([]Reading{}, s.readings[sensorID]...)
}

func (s *Service) Sensors() []string {
	return append([]string{}, s.sensorOrder...)
}

func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
