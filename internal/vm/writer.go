package vm

import (
	"fmt"
	"github.com/certainty/go-braces/internal/isa"
	"log"
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
	switch value := v.(type) {
	case isa.Char:
		return w.writeChar(value)
	case isa.UInt8, isa.UInt16, isa.UInt32, isa.UInt64:
		return fmt.Sprintf("u%d", value)
	case isa.Int, isa.Int8, isa.Int16, isa.Int32, isa.Int64:
		return fmt.Sprintf("%d", value)
	case isa.Float:
		return fmt.Sprintf("%f", value)
	case isa.String:
		return fmt.Sprintf("%q", string(value))
	default:
		log.Panicf("unhandled value type: %T", value)
		return ""
	}
}

func (w *Writer) writeChar(value isa.Char) string {
	switch value {
	case '\n':
		return "#\\newline"
	case '\r':
		return "#\\return"
	case 8:
		return "#\\backspace"
	case 20:
		return "#\\space"
	case 127:
		return "#\\delete"
	case 27:
		return "#\\escape"
	case 0:
		return "#\\null"
	case 7:
		return "#\\alarm"
	default:
		return fmt.Sprintf("#\\%c", value)
	}
}
