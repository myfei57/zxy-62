package access

import (
	"errors"

	"github.com/google/uuid"

	"tunnelnet/internal/store"
)

type Request struct {
	ID        string `json:"id"`
	CabinID   string `json:"cabin_id"`
	Requester string `json:"requester"`
	Status    string `json:"status"`
}

func (s *Service) Submit(cabinID string, requester string) Request {
	request := Request{
		ID:        uuid.NewString(),
		CabinID:   cabinID,
		Requester: requester,
		Status:    "pending",
	}
	s.requests = append(s.requests, request)
	return request
}

func (s *Service) Pending() []Request {
	out := make([]Request, 0)
	for _, request := range s.requests {
		if request.Status == "pending" {
			out = append(out, request)
		}
	}
	return out
}

func (s *Service) Grant(requestID string) (store.Approval, error) {
	index := -1
	for i, request := range s.requests {
		if request.ID == requestID {
			index = i
			break
		}
	}
	if index == -1 {
		return store.Approval{}, errors.New("request not found")
	}
	request := s.requests[index]
	approval, err := s.Approve(request.CabinID, request.Requester)
	if err != nil {
		return store.Approval{}, err
	}
	s.requests[index].Status = "granted"
	return approval, nil
}

func (s *Service) RequestCount() int {
	return len(s.requests)
}
