package store

func (s *Store) Export(kind string) (map[string]any, error) {
	ids, err := s.List(kind)
	if err != nil {
		return nil, err
	}
	out := make(map[string]any)
	for _, id := range ids {
		var value map[string]any
		if err := s.Load(kind, id, &value); err != nil {
			return nil, err
		}
		out[id] = value
	}
	return out, nil
}

func (s *Store) ExportAll() (map[string]map[string]any, error) {
	out := make(map[string]map[string]any)
	for _, kind := range s.Kinds() {
		records, err := s.Export(kind)
		if err != nil {
			return nil, err
		}
		out[kind] = records
	}
	return out, nil
}
