package client

import (
	"encoding/json"
	"testing"
)

func BenchmarkCandy(b *testing.B) {
	i64Client := NewCandyDBClient[int64]()
	key, _ := json.Marshal("key")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		i64Client.Add(key, int64(i))
		i64Client.Find(key)
	}
}
