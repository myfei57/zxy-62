package audit

type MarkSet struct {
	marks map[string]bool
}

func NewMarkSet() *MarkSet {
	return &MarkSet{marks: make(map[string]bool)}
}

func (m *MarkSet) Record(key string) bool {
	if m.marks[key] {
		return false
	}
	m.marks[key] = true
	return true
}
