package dict

import "sync"

type RWMutex[T comparable] struct {
	*Dict[T, *sync.RWMutex]
}

func NewRWMutex[T comparable](_ T, capacity int) *RWMutex[T] {
	return &RWMutex[T]{
		Dict: NewDict[T, *sync.RWMutex](capacity),
	}
}

func (m RWMutex[T]) Lock(key T) {
	m.get(key).Lock()
}

func (m RWMutex[T]) get(key T) *sync.RWMutex {
	lock, _ := m.Get(key)

	if lock == nil {
		lock = &sync.RWMutex{}

		if !m.SetNX(key, lock) {
			lock, _ = m.Get(key)
		}
	}

	return lock
}

func (m RWMutex[T]) UnLock(key T) {
	m.get(key).Unlock()
}
