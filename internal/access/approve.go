package access

import (
	"time"

	"github.com/google/uuid"

	"tunnelnet/internal/store"
)

func (s *Service) Approve(cabinID string, requester string) (store.Approval, error) {
	approval := store.Approval{
		ID:        uuid.NewString(),
		CabinID:   cabinID,
		Requester: requester,
		GrantedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := s.approvals.SaveApproval(approval); err != nil {
		return store.Approval{}, err
	}
	if err := s.doors.ReleaseDoor(cabinID); err != nil {
		return store.Approval{}, err
	}
	return approval, nil
}
