package vm_test

import (
	"testing"

	"github.com/certainty/go-braces/internal/isa"
	"github.com/certainty/go-braces/internal/vm"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	tests := []struct {
		name     string
		input    isa.Value
		expected string
	}{
		{
			name:     "Write true BoolValue",
			input:    isa.BoolValue(true),
			expected: "#t",
		},
		{
			name:     "Write false BoolValue",
			input:    isa.BoolValue(false),
			expected: "#f",
		},
	}

	emptyStringTable := vm.NewInternedStringTable()
	writer := vm.NewWriter(emptyStringTable)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := writer.Write(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
