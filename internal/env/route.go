package env

func (s *Service) RouteAlarm(sensorID string) string {
	partitionID := s.sensorPartition[sensorID]
	return s.partition.CabinOf(partitionID)
}
