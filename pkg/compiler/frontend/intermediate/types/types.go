package types

import (
	"fmt"
	"strings"
)

type Type interface {
	String() string
	tpe()
}

type Bool struct{}
type Byte struct{}
type Int struct{}
type UInt struct{}
type Float struct{}
type Rational struct{}
type Complex complex128
type String struct{}
type Char struct{}
type Void struct{}

type U8Vec struct{}

type Struct struct {
	Name   string
	Fields map[string]Type
}

type Ref struct {
	BaseType Type
}

type Chan struct {
	ElementType Type
}

type Procedure struct {
	Params  []Type
	Results []Type
}

var (
	BoolType   = Bool{}
	ByteType   = Byte{}
	IntType    = Int{}
	FloatType  = Float{}
	CharType   = Char{}
	StringType = String{}
	VoidType   = Void{}
)

func (Struct) tpe()    {}
func (Ref) tpe()       {}
func (Chan) tpe()      {}
func (Procedure) tpe() {}
func (Bool) tpe()      {}
func (Byte) tpe()      {}
func (Int) tpe()       {}
func (UInt) tpe()      {}
func (Float) tpe()     {}
func (Rational) tpe()  {}
func (Complex) tpe()   {}
func (String) tpe()    {}
func (Char) tpe()      {}
func (Void) tpe()      {}
func (U8Vec) tpe()     {}

func (Int) String() string      { return "int" }
func (UInt) String() string     { return "uint" }
func (Float) String() string    { return "float" }
func (Rational) String() string { return "rational" }
func (Complex) String() string  { return "complex" }
func (Bool) String() string     { return "bool" }
func (Char) String() string     { return "char" }
func (String) String() string   { return "string" }
func (Byte) String() string     { return "byte" }
func (Void) String() string     { return "void" }
func (U8Vec) String() string    { return "u8vec" }

func (p Ref) String() string {
	return fmt.Sprintf("*%s", p.BaseType.String())
}

func (c Chan) String() string {
	return fmt.Sprintf("chan[%s]", c.ElementType.String())
}

func (p Procedure) String() string {
	params := make([]string, len(p.Params))

	for idx, p := range p.Params {
		params[idx] = p.String()
	}

	results := make([]string, len(p.Results))
	for idx, r := range p.Results {
		results[idx] = r.String()
	}
	return fmt.Sprintf("proc(%s) -> (%s)", strings.Join(params, ", "), strings.Join(results, ", "))
}

func (s Struct) String() string {
	var fields []string
	for name, t := range s.Fields {
		fields = append(fields, fmt.Sprintf("%s: %s", name, t.String()))
	}
	return fmt.Sprintf("struct{%s}", fields)
}
