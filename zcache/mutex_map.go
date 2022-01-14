package zcache

import (
	"fmt"
	"sync"
)

func NewMutexMap() *MutexMap {
	return &MutexMap{
		mutexes: make(map[string]*sync.Mutex),
	}
}

type MutexMap struct {
	mut     sync.Mutex
	mutexes map[string]*sync.Mutex
}

func (m *MutexMap) Lock(id string) sync.Locker {
	m.mut.Lock()
	defer m.mut.Unlock()
	mut, exists := m.mutexes[id]
	if !exists {
		mut = &sync.Mutex{}
		m.mutexes[id] = mut
	}
	mut.Lock()
	return mut
}

func (m *MutexMap) Unlock(id string) {
	m.mut.Lock()
	defer m.mut.Unlock()
	mut, exists := m.mutexes[id]
	if !exists {
		panic(fmt.Errorf("mutex with id %s must exist before calling unlock", id))
	}
	mut.Unlock()
}
