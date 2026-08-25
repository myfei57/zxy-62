package verifycase

import (
	"testing"

	"tunnelnet/internal/drain"
	"tunnelnet/internal/store"
)

type recordMarkStore struct {
	seen map[string]bool
}

func (m *recordMarkStore) ExecutionMarkExists(key string) bool {
	return m.seen[key]
}

func (m *recordMarkStore) SaveExecutionMark(mark store.ExecutionMark) error {
	m.seen[mark.Key] = true
	return nil
}

type noopValve struct{}

func (noopValve) OpenValve(string) error { return nil }

type noopPump struct{}

func (noopPump) Start() error { return nil }

func TestUtDrainNoDuplicate(t *testing.T) {
	marks := &recordMarkStore{seen: map[string]bool{}}
	service := drain.NewService(noopValve{}, noopPump{}, marks)
	runs := 0
	run := func() error {
		runs++
		return nil
	}
	service.Execute("cmd-1", run)
	service.Execute("cmd-1", run)
	if runs != 1 {
		t.Fatalf("command ran %d times, want exactly once", runs)
	}
}
