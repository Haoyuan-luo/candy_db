package common

import (
	"sync"
	"testing"
)

func BenchmarkFan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		done := make(chan struct{})
		pd := Producer[int](done, testData)
		pc := Processor[int, int](done, pd, func(r int) int {
			return r
		})
		Consumer[int](done, pc)
	}
}

func BenchmarkFanWithSp(b *testing.B) {
	wg := &sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		sp := NewSplitter(testData)
		for {
			got, ok := sp.Get()
			if !ok {
				break
			}
			wg.Add(1)
			go func() {
				done := make(chan struct{})
				pd := Producer[int](done, got)
				pc := Processor[int, int](done, pd, func(r int) int {
					return r
				})
				Consumer[int](done, pc)
				wg.Done()
			}()
		}
	}
	wg.Wait()
}
