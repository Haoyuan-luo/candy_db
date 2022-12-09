package util

import (
	"fmt"
	"strconv"
	"testing"
)

func BenchmarkGS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := MapChan(MapChan(Filter[int, intJudge](testJudge), func(t int) string {
			return strconv.FormatInt(int64(t), 10)
		}), func(t string) int {
			k, _ := strconv.ParseInt(t, 10, 64)
			return int(k)
		})
		for v := range res {
			fmt.Println(v)
		}
	}
}
