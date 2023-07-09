package isa

type Value interface{}

type Char rune

func (v Char) String() string {
	return "CharValue"
}

type String string

func (v String) String() string {
	return "StringValue"
}

type Bool bool
type Int8 int8
type Int16 int16
type Int32 int32
type Int64 int64
type UInt8 uint8
type UInt16 uint16
type UInt32 uint32
type UInt64 uint64
type UInt uint
type Int int

func (Int) String() string {
	return "IntegerValue"
}

type Float float64

func (v Float) String() string {
	return "FloatValue"
}

type Arity interface{}
type AtLeast uint
type Exactly uint

type Function struct {
	Label string
	Code  CodeUnit
	Arity Arity
}

func (p *Function) String() string {
	return "ProcedureValue"
}

type Closure struct {
	Function Function
	UpValues []*Value
}

func (c *Closure) String() string {
	return "ClosureValue"
}
