package isa

import (
	"fmt"
)

// Values are not actullay part of the ISA normally
// and this file will go away eventually, but in the
// beginning it's very helpful to take some shortcuts
// during development. The compiler doesn't need to care too
// much about layout and encoding of the values in the constant pool
// for example.

type Value interface {
	Inspect() string
}

type BoolValue bool

func (b BoolValue) Inspect() string {
	return "BoolValue"
}

type ProcedureValue struct {
	Code CodeUnit
	// more to come later: like arity
}

func (p *ProcedureValue) Inspect() string {
	return "ProcedureValue"
}

type ClosureValue struct {
	Procedure ProcedureValue
	UpValues  []*Value
}

func (c *ClosureValue) Inspect() string {
	return "ClosureValue"
}

func ValueFromDatum(datum Datum) (Value, error) {
	switch datum := datum.(type) {
	case DatumBool:
		// at this point location information is stripped
		return BoolValue(datum.Value), nil
	default:
		return nil, fmt.Errorf("Can't convert datum to value")
	}
}
