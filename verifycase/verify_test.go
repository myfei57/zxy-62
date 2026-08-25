package verifycase

import (
	"testing"

	"tunnelnet/internal/fire"
)

type suppressorRecorder struct {
	order *[]string
}

func (s suppressorRecorder) Activate(zone string) {
	*s.order = append(*s.order, zone)
}

type noopDamper struct{}

func (noopDamper) Save(string, string, any) error { return nil }

type noopVent struct{}

func (noopVent) Stop() {}

type noopDeduper struct{}

func (noopDeduper) Record(string) bool { return true }

func TestUtFireZoneOrder(t *testing.T) {
	order := []string{}
	service := fire.NewService(noopDamper{}, noopVent{}, suppressorRecorder{&order}, noopDeduper{}, []string{"zone-a", "zone-b", "zone-c"})
	service.ActivateZones("zone-b")
	if len(order) == 0 || order[0] != "zone-b" {
		t.Fatalf("source zone must activate first: %v", order)
	}
}
