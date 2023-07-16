package vm_test

import (
	"testing"

	"github.com/certainty/go-braces/pkg/shared/isa"
	"github.com/certainty/go-braces/pkg/vm"
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
			input:    isa.Bool(true),
			expected: "true",
		},
		{
			name:     "Write false BoolValue",
			input:    isa.Bool(false),
			expected: "false",
		},
		{
			name:     "Write #\\c CharValue",
			input:    isa.Char('c'),
			expected: "#\\c",
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
