package dict

import "testing"

func TestNewSliceFromSliceGeneric(t *testing.T) {
	var (
		data = []testSlice{
			{
				A: 1,
				B: `a`,
			},
			{
				A: 1,
				B: `b`,
			},
			{
				A: 2,
				B: `a`,
			},
			{
				A: 3,
				B: "c",
			},
		}
	)

	result := NewSliceFromSliceGeneric(data, func(slice testSlice) (int, string) {
		return slice.A, slice.B
	})

	result.ForRange(func(key int, value []string) {
		t.Log(key, value)
	})
}

type testSlice struct {
	A int
	B string
}
