package store

import (
	"os"
	"strings"
)

func (s *Store) Delete(kind string, id string) error {
	return os.Remove(s.recordPath(kind, id))
}

func (s *Store) Count(kind string) (int, error) {
	ids, err := s.List(kind)
	if err != nil {
		return 0, err
	}
	return len(ids), nil
}

func (s *Store) Kinds() []string {
	entries, err := os.ReadDir(s.Root)
	if err != nil {
		return []string{}
	}
	out := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			out = append(out, entry.Name())
		}
	}
	return out
}

func (s *Store) Stats() map[string]int {
	out := make(map[string]int)
	for _, kind := range s.Kinds() {
		n, err := s.Count(kind)
		if err == nil {
			out[kind] = n
		}
	}
	return out
}
