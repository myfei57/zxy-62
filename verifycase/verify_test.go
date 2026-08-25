package verifycase

import (
	"testing"

	"tunnelnet/internal/cabin"
	"tunnelnet/internal/cable"
	"tunnelnet/internal/env"
)

func TestUtCabinMappingFresh(t *testing.T) {
	table := cabin.NewPartitionTable()
	table.RePartition("part-1", "east")
	service := env.NewService(cable.NewService(1), table, 75.0)
	service.BindSensorAt("sensor-1", "part-1")
	if got := service.RouteAlarm("sensor-1"); got != "east" {
		t.Fatalf("unexpected route: %s", got)
	}
	table.RePartition("part-1", "west")
	if got := service.RouteAlarm("sensor-1"); got != "west" {
		t.Fatalf("alarm still routes to stale cabin: %s", got)
	}
}
