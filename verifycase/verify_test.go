package verifycase

import (
	"os"
	"path/filepath"
	"testing"

	"tunnelnet/internal/access"
	"tunnelnet/internal/cabin"
	"tunnelnet/internal/store"
)

func TestUtAccessDurableFirst(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "approval"), []byte("block"), 0o644); err != nil {
		t.Fatalf("block approval dir: %v", err)
	}
	storage, err := store.New(dir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	registry := cabin.NewRegistry()
	registry.Upsert(cabin.Cabin{ID: "cabin-1", Name: "east"})
	service := access.NewService(storage, registry)
	request := service.Submit("cabin-1", "worker")
	_, err = service.Grant(request.ID)
	if err == nil {
		t.Fatalf("expected the durable approval write to fail")
	}
	if open, _ := registry.DoorState("cabin-1"); open {
		t.Fatalf("door was released before the approval became durable")
	}
	if len(service.Pending()) != 1 {
		t.Fatalf("request should remain pending after the failed durable write")
	}
}
