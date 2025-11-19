package config_test

import (
	"loan_system/internal/pkg/config"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		// TODO: Add test cases.
		{
			name: "load config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Load()
		})
	}
}
