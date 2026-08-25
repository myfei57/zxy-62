package verifycase

import (
	"testing"

	"tunnelnet/internal/fire"
)

type dedupSet struct {
	seen map[string]bool
}

func (d *dedupSet) Record(key string) bool {
	if d.seen[key] {
		return false
	}
	d.seen[key] = true
	return true
}

type noopDamper struct{}

func (noopDamper) Save(string, string, any) error { return nil }

type noopVent struct{}

func (noopVent) Stop() {}

type noopSuppressor struct{}

func (noopSuppressor) Activate(string) {}

func TestUtFireDedupSingle(t *testing.T) {
	service := fire.NewService(noopDamper{}, noopVent{}, noopSuppressor{}, &dedupSet{seen: map[string]bool{}}, []string{"zone-a"})
	service.RaiseAlarm("fire-event-1")
	service.RaiseAlarm("fire-event-1")
	if service.AlarmCount() != 1 {
		t.Fatalf("alarm count = %d, want 1", service.AlarmCount())
	}
}
