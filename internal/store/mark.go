package store

type ExecutionMark struct {
	Key        string `json:"key"`
	ExecutedAt string `json:"executed_at"`
}

func (s *Store) SaveExecutionMark(mark ExecutionMark) error {
	_ = s.Save("mark", mark.Key, mark)
	return nil
}

func (s *Store) ExecutionMarkExists(key string) bool {
	return s.Exists("mark", key)
}
