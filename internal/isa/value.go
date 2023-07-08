package isa

// Values are not actullay part of the ISA normally
// and this file will go away eventually, but in the
// beginning it's very helpful to take some shortcuts
// during development. The compiler doesn't need to care too
// much about layout and encoding of the values in the constant pool
// for example.

type Value interface{}

type CharValue rune

func (v CharValue) String() string {
	return "CharValue"
}

type StringValue string

func (v StringValue) String() string {
	return "StringValue"
}

type Int8Value int8
type Int16Value int16
type Int32Value int32
type Int64Value int64

type IntegerValue int64

func (v IntegerValue) String() string {
	return "IntegerValue"
}

type Uint8Value uint8
type Uint16Value uint16
type Uint32Value uint32
type Uint64Value uint64

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
