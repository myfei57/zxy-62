package cabin

import "errors"

func (r *Registry) Rename(id string, name string) error {
	item, ok := r.items[id]
	if !ok {
		return errors.New("cabin not found")
	}
	item.Name = name
	return nil
}

func (r *Registry) AssignZone(id string, zone string) error {
	item, ok := r.items[id]
	if !ok {
		return errors.New("cabin not found")
	}
	item.Zone = zone
	return nil
}

func (r *Registry) MoveCabin(id string, zone string, partition string) error {
	item, ok := r.items[id]
	if !ok {
		return errors.New("cabin not found")
	}
	item.Zone = zone
	item.Partition = partition
	return nil
}
