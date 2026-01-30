package dict

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

// CachedDict 是一个带缓存的字典，支持所有的map支持类型,每次本地不存在时，调用load获取，否则直接取本地哦，可以通过CleanAll()/Clean()来清理缓存
type CachedDict[T comparable, V any] struct {
	lock *sync.RWMutex
	data map[T]V
	load func(ctx context.Context, key T) (info V, err error)
}

func NewCachedDict[T comparable, V any](load func(ctx context.Context, key T) (info V, err error), initCapacity int) *CachedDict[T, V] {
	return &CachedDict[T, V]{
		lock: &sync.RWMutex{},
		data: make(map[T]V, initCapacity),
		load: load,
	}
}

// revive:disable:confusing-naming
func (d CachedDict[T, V]) Get(ctx context.Context, key T) (info V, err error) {
	var (
		exist bool
	)

	d.lock.RLock()

	if info, exist = d.data[key]; exist {
		d.lock.RUnlock()

		return info, nil
	}

	d.lock.RUnlock()

	d.lock.Lock()
	defer d.lock.Unlock()

	if info, err = d.load(ctx, key); err != nil {
		return *new(V), errors.Wrap(err, `加载`)
	}

	d.data[key] = info

	return info, nil
}

// revive:enable:confusing-naming

func (d CachedDict[T, V]) CleanAll() {
	d.lock.Lock()
	defer d.lock.Unlock()

	for k := range d.data {
		delete(d.data, k)
	}
}

func (d CachedDict[T, V]) Clean(key T) {
	d.lock.Lock()
	defer d.lock.Unlock()

	delete(d.data, key)
}
