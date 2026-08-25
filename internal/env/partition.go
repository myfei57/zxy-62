package env

func (s *Service) Repartition(partitionID string, cabinID string) {
	s.partition.RePartition(partitionID, cabinID)
}
