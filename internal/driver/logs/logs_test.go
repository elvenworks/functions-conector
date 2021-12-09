package logs

import "testing"

func TestInit(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"Success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init()
		})
	}
}
