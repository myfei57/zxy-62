package verifycase

import (
	"testing"

	"tunnelnet/internal/cabin"
	"tunnelnet/internal/patrol"
)

func TestUtPatrolRouteCurrent(t *testing.T) {
	registry := cabin.NewRegistry()
	registry.Upsert(cabin.Cabin{ID: "cabin-1", Name: "east"})
	service := patrol.NewService(registry)
	registry.Upsert(cabin.Cabin{ID: "cabin-2", Name: "west"})
	service.RebuildRoute()
	ids := service.RouteIDs()
	if !contains(ids, "cabin-2") {
		t.Fatalf("route does not include the new cabin: %v", ids)
	}
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
