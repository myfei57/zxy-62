package verifycase

import (
	"testing"

	"tunnelnet/internal/cabin"
	"tunnelnet/internal/quota"
)

func TestUtPowerBaselineFresh(t *testing.T) {
	registry := cabin.NewRegistry()
	registry.Upsert(cabin.Cabin{ID: "cabin-1", Name: "east"})
	book := cabin.NewPowerBook()
	book.SetBaseline("cabin-1", 100)
	service := quota.NewService(registry, book)
	registry.Upsert(cabin.Cabin{ID: "cabin-2", Name: "west"})
	service.RefreshBaseline(200)
	report := service.LoadReport()
	for _, row := range report {
		if row.CabinID == "cabin-2" && row.Watts == 200 {
			return
		}
	}
	t.Fatalf("new cabin missing from baseline report: %v", report)
}
