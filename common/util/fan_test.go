package util

import (
	"sync"
	"testing"
)

func BenchmarkFan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		done := make(chan struct{})
		pd := Producer[int](done, testData)
		p1 := Processor[int, int](done, pd, func(r int) int {
			return r
		}, nil)
		p2 := Processor[int, int](done, pd, func(r int) int {
			return r
		}, nil)
		p3 := Processor[int, int](done, pd, func(r int) int {
			return r
		}, nil)
		Consumer[int](done, p1, p2, p3)
	}
}

func BenchmarkFanWithSp(b *testing.B) {
	wg := &sync.WaitGroup{}
	var batchChan []<-chan int
	for i := 0; i < b.N; i++ {
		sp := NewSplitter(testData)
		for {
			got, ok := sp.Get()
			if !ok {
				break
			}
			wg.Add(1)
			singleChan := make(<-chan int)
			go func() {
				done := make(chan struct{})
				pd := Producer[int](done, got)
				pc := Processor[int, int](done, pd, func(r int) int {
					return r
				}, nil)
				singleChan = Consumer[int](done, pc)
				wg.Done()
			}()
			batchChan = append(batchChan, singleChan)
		}
		done := make(chan struct{})
		res := Consumer[int](done, batchChan...)
		wg.Wait()
		<-res
	}
}
