package util

import (
	"testing"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestLogger",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := Logger()
			logger.Log(INFO, "test logger tool")
			logger.Cost()
		})
	}
}
