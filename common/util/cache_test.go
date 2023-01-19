package util

import (
	"fmt"
	"testing"
)

func TestNewCacheService(t *testing.T) {
	type args struct {
		cacheType string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestNewCacheService",
			args: args{
				cacheType: Lfu,
			},
		},
	}
	for _, tt := range tests {
		cache := NewCacheService(tt.args.cacheType)
		cache.Set("key1", []byte("value"))
		cache.Get("key1")
		cache.Get("key1")
		cache.Set("key2", []byte("value"))
		cache.Get("key2")
		cache.Set("key3", []byte("value"))
		cache.Get("key3")
		cache.Set("key4", []byte("value"))
		cache.Get("key4")
		cache.Get("key4")
		cache.Get("key4")
		cache.Set("key5", []byte("value"))
		cache.Get("key5")
		cache.Set("key6", []byte("value"))
		cache.Get("key6")
		cache.Set("key7", []byte("value"))
		cache.Get("key7")
		cache.Set("key8", []byte("value"))
		cache.Get("key8")
		cache.Get("key8")
		cache.Get("key8")
		cache.Get("key8")
		cache.Set("key9", []byte("value"))
		cache.Get("key9")
		cache.Set("key10", []byte("value"))
		cache.Get("key10")
	}
}

func Test_newCmkRow(t *testing.T) {
	type args struct {
		numContainer int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_newCmkRow",
			args: args{
				numContainer: 2,
			},
		},
	}
	h := Hash()
	k64 := h.CalcHash([]byte("key"))

	for _, tt := range tests {
		r := newCmkRow(tt.args.numContainer)
		r.increment(uint64(k64 % uint64(tt.args.numContainer)))
		r.increment(uint64(k64 % uint64(tt.args.numContainer)))
		r.increment(uint64(k64 % uint64(tt.args.numContainer)))
		x := int64(r.get(uint64(k64 % uint64(tt.args.numContainer))))
		fmt.Println(x)
	}
}
