package vm

import (
	"fmt"
	"log"

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
	switch value := v.(type) {
	case isa.BoolValue:
		if value {
			return "true"
		} else {
			return "false"
		}
	case isa.CharValue:
		return w.writeChar(value)
	case isa.IntegerValue:
		return fmt.Sprintf("%d", value)
	case isa.FloatValue:
		return fmt.Sprintf("%f", value)
	case isa.StringValue:
		return fmt.Sprintf("%q", string(value))
	default:
		log.Panicf("unhandled value type: %T", value)
		return ""
	}
}

func (w *Writer) writeChar(value isa.CharValue) string {
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
