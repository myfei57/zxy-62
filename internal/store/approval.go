package store

type Approval struct {
	ID        string `json:"id"`
	CabinID   string `json:"cabin_id"`
	Requester string `json:"requester"`
	GrantedAt string `json:"granted_at"`
}

func (s *Store) SaveApproval(approval Approval) error {
	return s.Save("approval", approval.ID, approval)
}

func (s *Store) LoadApproval(id string) (Approval, error) {
	var approval Approval
	err := s.Load("approval", id, &approval)
	return approval, err
}

func (s *Store) ApprovalExists(id string) bool {
	return s.Exists("approval", id)
}

func (s *Store) ListApprovals() ([]Approval, error) {
	ids, err := s.List("approval")
	if err != nil {
		return nil, err
	}
	out := make([]Approval, 0, len(ids))
	for _, id := range ids {
		approval, err := s.LoadApproval(id)
		if err != nil {
			return nil, err
		}
		out = append(out, approval)
	}
	return out, nil
}
