package env

var routeCache = map[string]string{}

func (s *Service) RouteAlarm(sensorID string) string {
	if cabin, ok := routeCache[sensorID]; ok {
		return cabin
	}
	partitionID := s.sensorPartition[sensorID]
	cabin := s.partition.CabinOf(partitionID)
	routeCache[sensorID] = cabin
	return cabin
}
