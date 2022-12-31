package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAppendService(t *testing.T) {
	t.Run("parallel append", func(t *testing.T) {
		service := NewAppendImpl[int](100)
		for i := 0; i < 100; i++ {
			tuple := service.Ready()
			go func(j int) {
				tuple(j)
			}(i)
		}
		service.Close()
		assert.Equal(t, 100, len(service.Get()))
	})
}
