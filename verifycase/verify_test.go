package verifycase

import (
	"testing"

	"tunnelnet/internal/fire"
)

type damperRecorder struct {
	order *[]string
}

func (d damperRecorder) Save(string, string, any) error {
	*d.order = append(*d.order, "damper")
	return nil
}

type ventRecorder struct {
	order *[]string
}

func (v ventRecorder) Stop() {
	*v.order = append(*v.order, "vent")
}

type noopSuppressor struct{}

func (noopSuppressor) Activate(string) {}

type noopDeduper struct{}

func (noopDeduper) Record(string) bool { return true }

func TestUtVentFireOrder(t *testing.T) {
	order := []string{}
	service := fire.NewService(damperRecorder{&order}, ventRecorder{&order}, noopSuppressor{}, noopDeduper{}, []string{"zone-1"})
	if err := service.HandleSmoke("zone-1"); err != nil {
		t.Fatalf("handle smoke: %v", err)
	}
	if len(order) != 2 || order[0] != "damper" || order[1] != "vent" {
		t.Fatalf("unexpected order: %v", order)
	}
}
