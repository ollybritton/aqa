package evaluator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModuleName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"collatz", "collatz"},
		{"collatz.aqa", "collatz"},
		{"collatz-code.aqa", "collatz_code"},
		{"~/example/collatz-code.aqa", "collatz_code"},
		{"~/example/", "example"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, pathToModuleName(tt.input))
	}
}
