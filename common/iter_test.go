package common

import "testing"

func BenchmarkGS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Map(testData, func(t int) int {
			return t
		})
	}
}
