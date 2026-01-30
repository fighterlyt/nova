package pools

import "sync"

type Pool[T any] struct {
	data *sync.Pool
}

func NewPool[T any](initCapacity int, newFunc func() T) *Pool[T] {
	result := &Pool[T]{
		data: &sync.Pool{
			New: func() any {
				return newFunc()
			},
		},
	}

	for i := 0; i < initCapacity; i++ {
		result.data.Put(newFunc())
	}

	return result
}

func (p *Pool[T]) Get() T {
	return p.data.Get().(T)
}

func (p *Pool[T]) Put(t T) {
	p.data.Put(t)
}
