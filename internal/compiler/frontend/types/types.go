package types

import (
	"fmt"
	"strings"
)

type Type interface {
	String() string
}

type Int struct{}
type UInt struct{}
type Byte struct{}
type Float struct{}
type Rational struct{}
type Complex complex128
type Bool struct{}
type String struct{}
type Char struct{}
type Unknown struct{}
type Void struct{}

type Struct struct {
	Name   string
	Fields map[string]Type
}

type Set struct {
	ElementType Type
}

type Map struct {
	KeyType, ValueType Type
}

type Vector struct {
	ElementType Type
}

type Ref struct {
	BaseType Type
}

type Chan struct {
	ElementType Type
}

type EnumVariant struct {
	Name string
	Type Type
}

type Enum struct {
	Name     string
	Variants []EnumVariant
}

type Function struct {
	Params  []Type
	Results []Type
}

type Procedure struct {
	Params  []Type
	Results []Type
}

var (
	IntType     = Int{}
	UIntType    = UInt{}
	FloatType   = Float{}
	VoidType    = Void{}
	BoolType    = Bool{}
	StringType  = String{}
	CharType    = Char{}
	UnknownType = Unknown{}
)

func (Unknown) String() string {
	return "unknown"
}

func (Int) String() string {
	return "int"
}

func (Byte) String() string {
	return "byte"
}

func (UInt) String() string {
	return "uint"
}

func (Float) String() string {
	return "float"
}

func (Bool) String() string {
	return "bool"
}

func (String) String() string {
	return "string"
}

func (Char) String() string {
	return "char"
}

func (m Map) String() string {
	return fmt.Sprintf("map[%s]%s", m.KeyType.String(), m.ValueType.String())
}

func (s Set) String() string {
	return fmt.Sprintf("set[%s]", s.ElementType.String())
}

func (v Vector) String() string {
	return fmt.Sprintf("vector[%s]", v.ElementType.String())
}

func (r Ref) String() string {
	return fmt.Sprintf("&%s", r.BaseType.String())
}

func (c Chan) String() string {
	return fmt.Sprintf("chan[%s]", c.ElementType.String())
}

func (s Struct) String() string {
	var fields []string
	for name, t := range s.Fields {
		fields = append(fields, fmt.Sprintf("%s: %s", name, t.String()))
	}
	return fmt.Sprintf("struct{%s}", fields)
}

func (e Enum) String() string {
	var variants []string
	for _, v := range e.Variants {
		variants = append(variants, v.String())
	}
	return fmt.Sprintf("enum{%s}", strings.Join(variants, ", "))
}

func (e EnumVariant) String() string {
	return fmt.Sprintf("%s(%s)", e.Name, e.Type.String())
}

func (p Procedure) String() string {
	var params []string
	for _, t := range p.Params {
		params = append(params, t.String())
	}
	var results []string
	for _, t := range p.Results {
		results = append(results, t.String())
	}
	return fmt.Sprintf("proc(%s) -> (%s)", strings.Join(params, ", "), strings.Join(results, ", "))
}

func (f Function) String() string {
	var params []string
	for _, t := range f.Params {
		params = append(params, t.String())
	}
	var results []string
	for _, t := range f.Results {
		results = append(results, t.String())
	}
	return fmt.Sprintf("func(%s) -> (%s)", strings.Join(params, ", "), strings.Join(results, ", "))
}
