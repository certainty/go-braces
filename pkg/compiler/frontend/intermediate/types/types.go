package ir

import (
	"fmt"
	"strings"
)

type Type interface {
	Name() string
}

type Bool struct{}
type Byte struct{}
type Int struct{}
type UInt struct{}
type Float struct{}
type Rational struct{}
type Complex struct{}
type String struct{}
type Char struct{}
type Unit struct{}

type Vector struct {
	elementType Type
}

type ByteVector struct{}

type Tuple struct {
	elementTypes []Type
}

type Place struct {
	BaseType Type
}

type SetablePlace struct {
	BaseType Type
}

type Chan struct {
	ElementType Type
}

type Func struct {
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
	UnitType   = Unit{}
)

func (Int) Name() string      { return "int" }
func (UInt) Name() string     { return "uint" }
func (Float) Name() string    { return "float" }
func (Rational) Name() string { return "rational" }
func (Complex) Name() string  { return "complex" }
func (Bool) Name() string     { return "bool" }
func (Char) Name() string     { return "char" }
func (String) Name() string   { return "string" }
func (Byte) Name() string     { return "byte" }
func (Unit) Name() string     { return "unit" }

func (p Place) Name() string {
	return fmt.Sprintf("*%s", p.BaseType.Name())
}

func (p SetablePlace) Name() string {
	return fmt.Sprintf("mut *%s", p.BaseType.Name())
}

func (c Chan) Name() string {
	return fmt.Sprintf("chan[%s]", c.ElementType.Name())
}

func (f Func) Name() string {
	params := make([]string, len(f.Params))
	for _, p := range f.Params {
		params = append(params, p.Name())
	}
	results := make([]string, len(f.Results))
	for _, r := range f.Results {
		results = append(results, r.Name())
	}
	return fmt.Sprintf("func(%s) -> (%s)", strings.Join(params, ", "), strings.Join(results, ", "))
}

func (v Vector) Name() string {
	return fmt.Sprintf("vector[%s]", v.elementType.Name())
}

func (t Tuple) Name() string {
	elements := make([]string, len(t.elementTypes))
	for _, e := range t.elementTypes {
		elements = append(elements, e.Name())
	}
	return fmt.Sprintf("(%s)", strings.Join(elements, ", "))
}
