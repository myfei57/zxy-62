package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

type Store struct {
	Root string
}

func New(root string) (*Store, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	return &Store{Root: root}, nil
}

func (s *Store) recordPath(kind string, id string) string {
	return filepath.Join(s.Root, kind, id+".json")
}

func (s *Store) Save(kind string, id string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Join(s.Root, kind)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.recordPath(kind, id), data, 0o644)
}

func (s *Store) Load(kind string, id string, target any) error {
	data, err := os.ReadFile(s.recordPath(kind, id))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func (s *Store) Exists(kind string, id string) bool {
	_, err := os.Stat(s.recordPath(kind, id))
	return err == nil
}

func (s *Store) List(kind string) ([]string, error) {
	dir := filepath.Join(s.Root, kind)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	ids := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		ids = append(ids, entry.Name()[:len(entry.Name())-len(".json")])
	}
	sort.Strings(ids)
	return ids, nil
}
