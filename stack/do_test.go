package stack

import "testing"

func TestStackFibonacci(t *testing.T) {
	var (
		slice = map[string][]string{
			``:   []string{`a`, `b`, `c`},
			`a`:  []string{`ab`, `ac`},
			`b`:  []string{`bc`},
			`c`:  []string{`ca`},
			`ab`: []string{`abc`, `abd`},
		}
	)

	path := Stack[string, string]([]string{``}, func(t string) []string {
		return slice[t]
	}, func(a string) string {
		return a
	}, 1)

	t.Log(path)
}

func TestStackSum2(t *testing.T) {
	var (
		sum int64
	)

	path := Stack[int64, int64]([]int64{100}, func(i int64) []int64 {
		sum += i

		if i > 1 {
			return []int64{i - 1}
		}

		return nil
	}, func(i int64) int64 {
		return i
	}, 10)

	t.Log(sum)
	t.Log(path)
}
