package common

import (
	"strconv"
	"testing"
)

func TestNewSkipList(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNewSkipList",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSkipList()
			for i := 0; i < 100000; i++ {
				got.AddNode(Container{
					key:   []byte(strconv.FormatInt(int64(i), 10)),
					input: i,
				})
			}
			one := &Container{
				key: []byte(strconv.FormatInt(int64(10), 10)),
			}
			two := &Container{
				key: []byte(strconv.FormatInt(int64(100), 10)),
			}
			three := &Container{
				key: []byte(strconv.FormatInt(int64(1000), 10)),
			}
			got.FindNode(one, two, three)
		})
	}
}
