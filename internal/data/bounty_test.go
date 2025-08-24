package data

import (
	"testing"
)

func TestBerries_MarshalJSON(t *testing.T) {
	tests := []struct {
		berries  Berries
		expected string
	}{
		{100, `"100 berries"`},
		{1500000, `"1.5M berries"`},
		{30000000, `"30M berries"`},
		{1500000000, `"1.5B berries"`},
		{5000000000, `"5B berries"`},
	}

	for _, tt := range tests {
		result, err := tt.berries.MarshalJSON()
		if err != nil {
			t.Errorf("MarshalJSON() error = %v", err)
		}
		if string(result) != tt.expected {
			t.Errorf("MarshalJSON(%d) = %s, want %s", tt.berries, string(result), tt.expected)
		}
	}
}

func TestBerries_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected Berries
	}{
		{`"100 berries"`, 100},
		{`"1.5M berries"`, 1500000},
		{`"30M berries"`, 30000000},
		{`"1.5B berries"`, 1500000000},
		{`"5B berries"`, 5000000000},
	}

	for _, tt := range tests {
		var b Berries
		err := b.UnmarshalJSON([]byte(tt.input))
		if err != nil {
			t.Errorf("UnmarshalJSON(%s) error = %v", tt.input, err)
		}
		if b != tt.expected {
			t.Errorf("UnmarshalJSON(%s) = %v, want %v", tt.input, b, tt.expected)
		}
	}
}
