package dict

import (
	"sync"
)

type Slice[T comparable, V any] struct {
	*Dict[T, []V]
}

func NewSlice[T comparable, V any](initCapacity int) *Slice[T, V] {
	return &Slice[T, V]{
		Dict: NewDict[T, []V](initCapacity),
	}
}

func (d *Slice[T, V]) Append(key T, values ...V) {
	d.lock.Lock()
	defer d.lock.Unlock()

	data := d.data[key]

	if data == nil {
		data = []V{}
	}

	data = append(data, values...)

	d.data[key] = data
}

func (d *Slice[T, V]) ForRangeConcurrentWithKey(key T, process func(value V)) {
	if d == nil {
		return
	}

	d.lock.RLock()
	defer d.lock.RUnlock()

	wg := &sync.WaitGroup{}

	for _, item := range d.data[key] {
		wg.Add(1)

		go func(item V) {
			defer wg.Done()
			process(item)
		}(item)
	}

	wg.Wait()
}

func (d *Slice[T, V]) RemoveSliceElement(key T, match func(V) bool) {
	if d == nil {
		return
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	value := d.data[key]

	value = remove(value, match)

	d.data[key] = value
}

func remove[T any](slice []T, remove func(T) bool) []T {
	for i := 0; i < len(slice); i++ {
		if !remove(slice[i]) {
			continue
		}

		if i != len(slice)-1 {
			copy(slice[i:], slice[i+1:])
		}

		slice = slice[:len(slice)-1]
		i--
	}

	return slice
}

func NewSliceFromSlice[T comparable, V any](slice []V, get func(V) T) *Slice[T, V] {
	result := NewSlice[T, V](len(slice))

	for i := range slice {
		t := get(slice[i])

		result.Append(t, slice[i])
	}

	return result
}

func NewSliceFromSliceGeneric[T comparable, V, Item any](slice []Item, get func(Item) (T, V)) *Slice[T, V] {
	result := NewSlice[T, V](len(slice))

	for i := range slice {
		key, v := get(slice[i])

		result.Append(key, v)
	}

	return result
}
