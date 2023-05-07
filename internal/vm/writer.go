package vm

import (
	"log"
	"reflect"

	"github.com/certainty/go-braces/internal/isa"
)

type Writer struct {
	internedStrings *InternedStringTable
}

func NewWriter(internedStrings *InternedStringTable) *Writer {
	return &Writer{
		internedStrings: internedStrings,
	}
}

func (w *Writer) Write(v isa.Value) string {
	log.Printf("v: %s", reflect.TypeOf(v))
	log.Printf("v: %s", v.Inspect())
	switch value := v.(type) {
	case isa.BoolValue:
		if value {
			return "#t"
		} else {
			return "#f"
		}
	default:
		panic("CompilerBug: unknown value")
	}
}
