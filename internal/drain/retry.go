package drain

import "tunnelnet/internal/store"

func (s *Service) Execute(commandKey string, run func() error) error {
	if err := run(); err != nil {
		return err
	}
	s.runs++
	_ = s.marks.SaveExecutionMark(store.ExecutionMark{Key: commandKey, ExecutedAt: s.now().UTC().Format("2006-01-02T15:04:05Z")})
	return nil
}
