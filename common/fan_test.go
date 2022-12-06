package common

import (
	"testing"
)

func BenchmarkFan(b *testing.B) {
	type args struct {
		done chan struct{}
	}
	for i := 0; i < b.N; i++ {
		done := make(chan struct{})
		pd := Producer[int](done, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
		pc := Processor[int, int](done, pd, func(r int) int {
			return r
		})
		Consumer[int](done, pc)
	}
}
