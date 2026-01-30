package dict

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTuple(t *testing.T) {
	tuple := NewTuple(``, ``, func() int64 {
		return 0
	}, 10)

	require.Zero(t, tuple.Get(`a`, `b`))

	tuple.Set(`a`, `b`, 1)

	require.EqualValues(t, 1, tuple.Get(`a`, `b`))

	tuple.SetDynamic(`a`, `b`, func(i int64) int64 {
		return i + 1
	})

	require.EqualValues(t, 2, tuple.Get(`a`, `b`))

	tuple.ForRange(func(a string, b string, c int64) {
		t.Logf(`a[%s],b[%s],c[%d]`, a, b, c)
	})
}

func TestTuple_SetNx(t *testing.T) {
	tuple := NewTuple(``, ``, func() int64 {
		return 0
	}, 10)

	require.Zero(t, tuple.Get(`a`, `b`))

	tuple.Set(`a`, `b`, 1)

	require.EqualValues(t, 1, tuple.Get(`a`, `b`))

	tuple.SetNX(`a`, `b`, 2)

	require.EqualValues(t, 1, tuple.Get(`a`, `b`))

	tuple.Set(`a`, `b`, 2)

	require.EqualValues(t, 2, tuple.Get(`a`, `b`))
}
