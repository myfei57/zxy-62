package cabin

type PartitionTable struct {
	mapping map[string]string
	revision int
}

func NewPartitionTable() *PartitionTable {
	return &PartitionTable{mapping: make(map[string]string)}
}

func (t *PartitionTable) RePartition(partitionID string, cabinID string) {
	t.mapping[partitionID] = cabinID
	t.revision++
}

func (t *PartitionTable) CabinOf(partitionID string) string {
	return t.mapping[partitionID]
}

func (t *PartitionTable) Revision() int {
	return t.revision
}
