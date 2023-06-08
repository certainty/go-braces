package isa

// Values are not actullay part of the ISA normally
// and this file will go away eventually, but in the
// beginning it's very helpful to take some shortcuts
// during development. The compiler doesn't need to care too
// much about layout and encoding of the values in the constant pool
// for example.

type Value interface{}

type BoolValue bool

func (b BoolValue) String() string {
	return "BoolValue"
}

type CharValue rune

func (v CharValue) String() string {
	return "CharValue"
}

type StringValue string

func (v StringValue) String() string {
	return "StringValue"
}

type IntegerValue int64

func (v IntegerValue) String() string {
	return "IntegerValue"
}

type FloatValue float64

func (v FloatValue) String() string {
	return "FloatValue"
}

type ProcedureValue struct {
	Code CodeUnit
	// more to come later: like arity
}

func (p *ProcedureValue) String() string {
	return "ProcedureValue"
}

type ClosureValue struct {
	Procedure ProcedureValue
	UpValues  []*Value
}

func (c *ClosureValue) String() string {
	return "ClosureValue"
}
