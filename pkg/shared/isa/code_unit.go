package isa

import "fmt"

type CodeUnit struct {
	Constants    []Value
	Instructions []Instruction
}

func EmptyCodeUnit() CodeUnit {
	return CodeUnit{
		Constants:    []Value{},
		Instructions: []Instruction{},
	}
}

func (c *CodeUnit) ReadConstant(address ConstantAddress) (Value, error) {
	if address >= ConstantAddress(len(c.Constants)) {
		return nil, fmt.Errorf("invalid constant address: %d", address)
	}
	return c.Constants[address], nil
}
