package dict

type Exist[T comparable] struct {
	*Dict[T, struct{}]
}

func NewExist[T comparable](initCapacity int) *Exist[T] {
	return &Exist[T]{
		Dict: NewDict[T, struct{}](initCapacity),
	}
}

func (e *Exist[T]) Exist(t T) bool {
	_, ok := e.Get(t)

	return ok
}

func (e *Exist[T]) Set(t T) {
	e.Dict.Set(t, struct{}{})
}

func (e *Exist[T]) Del(t T) {
	e.Dict.Delete(t)
}

func (e *Exist[T]) Clear() {
	e.DeleteAll()
}

type MultiKeyExist2[T comparable] struct {
	*Exist[T]
}

func NewMultiKeyDict[T comparable](initCapacity int) *MultiKeyExist2[[2]T] {
	return &MultiKeyExist2[[2]T]{
		Exist: NewExist[[2]T](initCapacity),
	}
}
