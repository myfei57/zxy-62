package access

import (
	"tunnelnet/internal/store"
)

type ApprovalStore interface {
	SaveApproval(store.Approval) error
	ApprovalExists(string) bool
}

type DoorController interface {
	ReleaseDoor(string) error
}

type Service struct {
	approvals ApprovalStore
	doors     DoorController
	requests  []Request
}

func NewService(approvals ApprovalStore, doors DoorController) *Service {
	return &Service{approvals: approvals, doors: doors, requests: make([]Request, 0)}
}

func (s *Service) Approved(id string) bool {
	return s.approvals.ApprovalExists(id)
}
