package store

import (
	"os"
	"path/filepath"
)

func (s *Store) FileCount() int {
	total := 0
	for _, kind := range s.Kinds() {
		n, err := s.Count(kind)
		if err == nil {
			total += n
		}
	}
	return total
}

func (s *Store) TotalBytes() int64 {
	var total int64
	_ = filepath.Walk(s.Root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && filepath.Ext(path) == ".json" {
			total += info.Size()
		}
		return nil
	})
	return total
}
