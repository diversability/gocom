package goroutineid

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
)

func TestGoID(t *testing.T) {
	var s sync.Map
	var w sync.WaitGroup
	var number int32 = 0
	for i := 0; i < 10000; i++ {
		w.Add(1)
		go func() {
			if _, ok := s.Load(GetGoID()); ok {
				t.Fatalf("fuck: %d", GetGoID())
			}

			atomic.AddInt32(&number, 1)
			fmt.Printf("new goid: %d %+v\n", GetGoID(), G())
			s.Store(GetGoID(), 1)
			w.Done()
		}()
	}

	w.Wait()

	var out []int64
	s.Range(func(k, v interface{}) bool {
		out = append(out, k.(int64))
		return true
	})

	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })

	for i := 0; i < len(out); i++ {
		fmt.Printf("%d\n", out[i])
	}

	for i := 0; i < len(out) - 1; i++ {
		if out[i] + 1 != out[i+1] {
			fmt.Printf("miss: %d\n", out[i] + 1)
		}
	}

	fmt.Printf("circle: %d out: %d\n", number, len(out))
}
