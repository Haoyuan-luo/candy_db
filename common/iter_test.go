package common

import "testing"

var test = func() (re []int) {
	for i := 0; i < 5000; i++ {
		re = append(re, i)
	}
	return
}()

func BenchmarkGS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Filter(test, func(i int) bool {
			return i%2 == 0
		})
	}
}

func BenchmarkGSV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FilterV2(test, func(i int) bool {
			return i%2 == 0
		})
	}
}
