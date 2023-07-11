package isa

import "github.com/certainty/go-braces/internal/isa/arity"

type Value interface{}

type Bool bool
type Byte byte
type Int int64
type UInt uint64
type Float float64
type Char rune
type String string
type Vector []Value
type Map map[Value]Value
type Closure struct {
	Function Function
	UpValues []*Value
}

type Function struct {
	Label Label
	Arity arity.Arity
	Code  CodeUnit
}

var _ Value = Bool(false)
var _ Value = Byte(0)
var _ Value = Int(0)
var _ Value = UInt(0)
var _ Value = Float(0.0)
var _ Value = Char(0)
var _ Value = String("")
var _ Value = (*Closure)(nil)
var _ Value = (*Function)(nil)
var _ Value = Vector{}
var _ Value = Map{}

func (v Char) String() string {
	return "CharValue"
}

func (v String) String() string {
	return "StringValue"
}

func (v Byte) String() string {
	return "ByteValue"
}

func (v Bool) String() string {
	return "BoolValue"
}

func (Int) String() string {
	return "IntegerValue"
}

func (v Float) String() string {
	return "FloatValue"
}

func (v Vector) String() string {
	return "VectorValue"
}

func (v Map) String() string {
	return "MapValue"
}

func (p *Function) String() string {
	return "ProcedureValue"
}

func (c *Closure) String() string {
	return "ClosureValue"
}
