package common

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestNewBloomFilter(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNewBloomFilter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBloomFilter()
			key1 := []byte(strconv.FormatInt(int64(1234567), 10))
			key2 := []byte(strconv.FormatInt(int64(2345678), 10))
			key3 := []byte(strconv.FormatInt(int64(3456789), 10))
			got.account = 3
			filter := got.GetBloomArray(key1, key2, key3)
			assert.Equal(t, filter, []int{1, 1, 1})
		})
	}
}
