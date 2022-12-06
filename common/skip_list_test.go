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
				got.AddNode(NewNode([]byte(strconv.FormatInt(int64(i), 10))).SetData(i))
			}
			ret := got.FindNode([]byte(strconv.FormatInt(int64(10000), 10)), []byte(strconv.FormatInt(int64(700), 10)), []byte(strconv.FormatInt(int64(300), 10)))
			if ret == nil {
				t.Errorf("TestNewSkipList() = %v, want %v", ret, 300)
			}
		})
	}
}
