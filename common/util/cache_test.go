package util

import (
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
				cacheType: LFU,
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
