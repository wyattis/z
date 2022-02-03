package zsync

import "sync"

func NewLockerMap(capacity int) *LockerMap {
	return &LockerMap{
		mutexes: make(map[string]*sync.Mutex, capacity),
	}
}

// A LockerMap manages a map of sync.Mutex for safely accessing resources that
// are distinguished by a key
type LockerMap struct {
	mutexes map[string]*sync.Mutex
}

// Safely read/write a resource using a function
func (m *LockerMap) With(key string, cb func()) {
	mut := m.Lock(key)
	defer mut.Unlock()
	cb()
}

// Lock a resource. Calls sync.Mutex.Lock() for the mutex defined for this key.
func (m *LockerMap) Lock(key string) sync.Locker {
	mut, exists := m.mutexes[key]
	if !exists {
		mut = &sync.Mutex{}
		m.mutexes[key] = mut
	}
	mut.Lock()
	return mut
}

// Unlock a resource that was locked. Calls sync.Mutex.Unlock() for the mutex
// defined for this key.
func (m *LockerMap) Unlock(key string) {
	mut, exists := m.mutexes[key]
	if !exists {
		panic("calling unlock on a mutex that doesn't exist")
	}
	mut.Unlock()
}

func NewRWLockerMap(capacity int) *RWLockerMap {
	return &RWLockerMap{
		mutexes: make(map[string]*sync.RWMutex, capacity),
	}
}

// Create a map of sync.RWMutex
type RWLockerMap struct {
	mutexes map[string]*sync.RWMutex
}

// Safely read/write a resource using a function
func (m *RWLockerMap) With(key string, cb func()) {
	mut := m.Lock(key)
	defer mut.Unlock()
	cb()
}

// Safely read from a resource using a function
func (m *RWLockerMap) RWith(key string, cb func()) {
	mut := m.RLock(key)
	defer mut.Unlock()
	cb()
}

// Lock a resource for reading. Calls sync.RWMutex.Lock on the mutex for this
// key.
func (m *RWLockerMap) Lock(key string) *sync.RWMutex {
	mut, exists := m.mutexes[key]
	if !exists {
		mut = &sync.RWMutex{}
		m.mutexes[key] = mut
	}
	mut.Lock()
	return mut
}

// Unlock a resource that was locked for reading a writing. Calls
// sync.RWMutex.Unlock on the mutex for this key.
func (m *RWLockerMap) Unlock(key string) {
	mut, exists := m.mutexes[key]
	if !exists {
		panic("calling unlock on a mutex that doesn't exist")
	}
	mut.Unlock()
}

// Lock a resource for reading. Calls sync.RWMutex.RLock on the mutex for this
// key.
func (m *RWLockerMap) RLock(key string) *sync.RWMutex {
	mut, exists := m.mutexes[key]
	if !exists {
		mut = &sync.RWMutex{}
		m.mutexes[key] = mut
	}
	mut.RLock()
	return mut
}

// Unlock a resource that was locked for reading. Calls sync.RWMutex.Unlock on
// the mutex for this key.
func (m *RWLockerMap) RUnlock(key string) {
	mut, exists := m.mutexes[key]
	if !exists {
		panic("calling unlock on a mutex that doesn't exist")
	}
	mut.RUnlock()
}
