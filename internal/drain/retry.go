package drain

import "tunnelnet/internal/store"

func (s *Service) Execute(commandKey string, run func() error) error {
	if s.IsInFlight(commandKey) || s.marks.ExecutionMarkExists(commandKey) {
		return nil
	}
	s.MarkInFlight(commandKey)
	if err := run(); err != nil {
		return err
	}
	s.runs++
	if err := s.marks.SaveExecutionMark(store.ExecutionMark{Key: commandKey, ExecutedAt: s.now().UTC().Format("2006-01-02T15:04:05Z")}); err != nil {
		return err
	}
	s.ClearInFlight(commandKey)
	return nil
}
