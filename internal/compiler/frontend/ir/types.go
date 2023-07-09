package ir

import (
	"fmt"
	"strings"
)

type Type interface {
	Name() string
}

type Bool struct{}
type Int8 struct{}
type Int16 struct{}
type Int32 struct{}
type Int64 struct{}
type UInt8 struct{}
type UInt16 struct{}
type UInt32 struct{}
type UInt64 struct{}
type Float32 struct{}
type String struct{}
type Char struct{}
type Table struct {
	KeyType, ValueType Type
}
type Ptr struct {
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
	BoolType    = Bool{}
	Int8Type    = Int8{}
	Int16Type   = Int16{}
	Int32Type   = Int32{}
	Int64Type   = Int64{}
	UInt8Type   = UInt8{}
	UInt16Type  = UInt16{}
	UInt32Type  = UInt32{}
	UInt64Type  = UInt64{}
	Float32Type = Float32{}
	CharType    = Char{}
	StringType  = String{}
)

func (Int8) Name() string    { return "int8" }
func (Int16) Name() string   { return "int16" }
func (Int32) Name() string   { return "int32" }
func (Int64) Name() string   { return "int64" }
func (UInt8) Name() string   { return "uint8" }
func (UInt16) Name() string  { return "uint16" }
func (UInt32) Name() string  { return "uint32" }
func (UInt64) Name() string  { return "uint64" }
func (Float32) Name() string { return "float32" }
func (Char) Name() string    { return "char" }
func (String) Name() string  { return "string" }

func (t Table) Name() string {
	return fmt.Sprintf("table[%s]%s", t.KeyType.Name(), t.ValueType.Name())
}

func (p Ptr) Name() string {
	return fmt.Sprintf("*%s", p.BaseType.Name())
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
