package store

type ExecutionMark struct {
	Key        string `json:"key"`
	ExecutedAt string `json:"executed_at"`
}

func (s *Store) SaveExecutionMark(mark ExecutionMark) error {
	return s.Save("mark", mark.Key, mark)
}

func (s *Store) ExecutionMarkExists(key string) bool {
	return s.Exists("mark", key)
}
