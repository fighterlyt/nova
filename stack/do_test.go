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
