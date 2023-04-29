package vm

import "github.com/certainty/go-braces/internal/isa/value"

type Writer struct {
	internedStrings *InternedStringTable
}

func NewWriter(internedStrings *InternedStringTable) *Writer {
	return &Writer{
		internedStrings: internedStrings,
	}
}

func (w *Writer) Write(v value.Value) string {
	switch value := v.(type) {
	case value.BoolValue:
		if value {
			return "#t"
		} else {
			return "#f"
		}
	default:
		panic("CompilerBug: unknown value")
	}
}
