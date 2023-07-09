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
type Float struct{}
type Bool struct{}
type String struct{}
type Char struct{}

type Struct struct {
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

type UserDefined struct {
	Name string
}

type Function struct {
	Params  []Type
	Results []Type
}

type Procedure struct {
	Params []Type
}

var (
	IntType    = Int{}
	UIntType   = UInt{}
	FloatType  = Float{}
	BoolType   = Bool{}
	StringType = String{}
	CharType   = Char{}
)

func (Int) String() string {
	return "int"
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

func (s Struct) String() string {
	var fields []string
	for name, t := range s.Fields {
		fields = append(fields, fmt.Sprintf("%s: %s", name, t.String()))
	}
	return fmt.Sprintf("struct{%s}", fields)
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

func (u UserDefined) String() string {
	return u.Name
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
