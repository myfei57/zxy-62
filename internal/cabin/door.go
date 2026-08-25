package cabin

import "errors"

func (r *Registry) ReleaseDoor(id string) error {
	item, ok := r.items[id]
	if !ok {
		return errors.New("cabin not found")
	}
	item.DoorOpen = true
	return nil
}

func (r *Registry) SealDoor(id string) error {
	item, ok := r.items[id]
	if !ok {
		return errors.New("cabin not found")
	}
	item.DoorOpen = false
	return nil
}

func (r *Registry) DoorState(id string) (bool, error) {
	item, ok := r.items[id]
	if !ok {
		return false, errors.New("cabin not found")
	}
	return item.DoorOpen, nil
}
