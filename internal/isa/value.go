package isa

// Values are not actullay part of the ISA normally
// and this file will go away eventually, but in the
// beginning it's very helpful to take some shortcuts
// during development. The compiler doesn't need to care too
// much about layout and encoding of the values in the constant pool
// for example.

type Value interface{}

type BoolValue bool

type ProcedureValue struct {
	Code CodeUnit
	// more to come later: like arity
}

type ClosureValue struct {
	Procedure ProcedureValue
	UpValues  []*Value
}
