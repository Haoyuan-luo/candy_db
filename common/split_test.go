package common

import (
	"testing"
)

func TestNewSplitter(t *testing.T) {
	type args struct {
		data []int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				data: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
		},
	}
	for _, tt := range tests {
		for sp := NewSplitter(tt.args.data); ; {
			got, ok := sp.Get()
			if !ok {
				break
			}
			l := Logger().SetField("Test split").Log(INFO, "%v", got)
			l.Cost()
		}
	}
}

var testData = func() (re []int) {
	for i := 0; i < 10000; i++ {
		re = append(re, i)
	}
	return
}()

func BenchmarkSplitter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for sp := NewSplitter(testData); ; {
			got, ok := sp.Get()
			if !ok {
				break
			}
			_ = got
		}
	}
}
