package main

import (
	"fmt"
	"sync"

	"github.com/onodera-punpun/oshirase"
)

type notifies struct {
	Notifies []oshirase.Notify
	NMU      sync.RWMutex
}

func newNotifies() *notifies {
	return &notifies{}
}

func (ns *notifies) add(n *oshirase.Notify) {
	ns.NMU.Lock()
	defer ns.NMU.Unlock()

	ns.Notifies = append(ns.Notifies, *n)
}

func (ns *notifies) delete(id uint32) error {
	idx, err := ns.findByID(id)
	if err != nil {
		return err
	}

	ns.NMU.Lock()
	defer ns.NMU.Unlock()
	ns.Notifies = append(ns.Notifies[:idx], ns.Notifies[idx+1:]...)

	return nil
}

func (ns *notifies) findByID(id uint32) (index int, err error) {
	ns.NMU.RLock()
	defer ns.NMU.RUnlock()

	for i, n := range ns.Notifies {
		if n.ID == id {
			return i, nil
		}
	}

	return -1, fmt.Errorf("id %d dosesn't exist", id)
}
