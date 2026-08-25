package cabin

import "errors"

func (r *Registry) OpenValve(id string) error {
	item, ok := r.items[id]
	if !ok {
		return errors.New("cabin not found")
	}
	item.ValveOpen = true
	return nil
}

func (r *Registry) CloseValve(id string) error {
	item, ok := r.items[id]
	if !ok {
		return errors.New("cabin not found")
	}
	item.ValveOpen = false
	return nil
}

func (r *Registry) ValveState(id string) (bool, error) {
	item, ok := r.items[id]
	if !ok {
		return false, errors.New("cabin not found")
	}
	return item.ValveOpen, nil
}
