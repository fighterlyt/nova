package dict

// revive:disable:confusing-naming  // 泛型不同类型，会误报

type Tuple[A comparable, B comparable, C any] struct {
	data *Dict[A, *Dict[B, C]]
	newC func() C
}

func NewTuple[A comparable, B comparable, C any](_ A, _ B, c func() C, capacity int) *Tuple[A, B, C] {
	return &Tuple[A, B, C]{
		data: NewDict[A, *Dict[B, C]](capacity),
		newC: c,
	}
}

func (t *Tuple[A, B, C]) Set(a A, b B, c C) {
	t.setWithLock(a, b, c, true)
}

func (t *Tuple[A, B, C]) setWithLock(a A, b B, c C, needLock bool) {
	bData, exist := t.data.getWithLock(a, needLock)

	if !exist {
		bData = NewDict[B, C](10)
	}

	bData.setWithLock(b, c, needLock)

	t.data.setWithLock(a, bData, needLock)
}

func (t *Tuple[A, B, C]) Get(a A, b B) C {
	var (
		bData, exist = t.data.Get(a)
		cValue       C
	)

	if !exist {
		return t.newC()
	}

	cValue, _ = bData.Get(b)

	return cValue
}

func (t *Tuple[A, B, C]) SetDynamic(a A, b B, c func(C) C) {
	var (
		bData, exist = t.data.Get(a)
		cValue       C
	)

	if !exist {
		bData = NewDict[B, C](10)
		cValue = t.newC()
	} else {
		cValue, exist = bData.Get(b)

		if !exist {
			cValue = t.newC()
		}
	}

	cValue = c(cValue)

	bData.Set(b, cValue)

	t.data.Set(a, bData)
}

func (t *Tuple[A, B, C]) ForRange(f func(a A, b B, c C)) {
	t.data.ForRange(func(a A, value *Dict[B, C]) {
		value.ForRange(func(b B, c C) {
			f(a, b, c)
		})
	})
}

func (t *Tuple[A, B, C]) Delete(a A, b B) bool {
	bData, exist := t.data.Get(a)

	if !exist {
		return false
	}

	_, exist = bData.Get(b)

	if !exist {
		return false
	}

	return bData.DeleteWithCheck(b)
}

func (t *Tuple[A, B, C]) DeleteNonExistA(a ...A) {
	t.data.DeleteNonExist(a...)
}

func (t *Tuple[A, B, C]) DeleteAll() {
	t.data.DeleteAll()
}

func (t *Tuple[A, B, C]) Len1(a A) int {
	aData, exist := t.data.Get(a)

	if !exist {
		return 0
	}

	return aData.Len()
}

func (t *Tuple[A, B, C]) SetNX(a A, b B, c C) {
	bData, exist := t.data.Get(a)

	if !exist {
		bData = NewDict[B, C](10)
	}

	bData.SetNX(b, c)

	t.data.SetNX(a, bData)
}

func (t *Tuple[A, B, C]) Get1(a A) []B {
	aData, exist := t.data.Get(a)

	if !exist {
		return nil
	}

	result := make([]B, 0, aData.Len())

	aData.ForRange(func(key B, _ C) {
		result = append(result, key)
	})

	return result
}

func (t *Tuple[A, B, C]) Reset(data map[A]map[B]C) {
	t.data.lock.Lock()
	defer t.data.lock.Unlock()

	t.data.deleteAllWithLock(false)

	for a := range data {
		for b := range data[a] {
			t.setWithLock(a, b, data[a][b], false)
		}
	}
}
