package dict

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	value []int
	sum   int
)

func BenchmarkLock(b *testing.B) {
	mutex := NewRWMutex(1, 10)

	concurrent := 1000
	times := 10000

	wg := &sync.WaitGroup{}

	b.ReportAllocs()
	b.ResetTimer()

	for k := 0; k < b.N; k++ {
		value = make([]int, concurrent)

		for i := 0; i < concurrent; i++ {
			wg.Add(1)

			go func(key int) {
				for j := 0; j < times; j++ {
					mutex.Lock(key)
					value[key]++
					mutex.UnLock(key)
				}

				wg.Done()
			}(i)
		}

		wg.Wait()

		for k, v := range value {
			require.EqualValuesf(b, times, v, `元素[%d]`, k)
			sum += v
		}
	}
}
