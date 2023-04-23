package bytecode

import (
	"fmt"

	"github.com/certainty/go-braces/vm/internal/language/value"
)

type Address = uint32
type ConstAddress = uint32

type Chunk struct {
	Instructions      []Instruction      // The sequence of instructions to execute
	Constants         []value.Value      // The constant values used by the instructions
	IntrospectionInfo *IntrospectionInfo // Auxiliary information for introspection
	DebugInfo         *DebugInfo         // Auxiliary information for debugging
}

func (c *Chunk) InstructionAt(address Address) (*Instruction, error) {
	if address >= uint32(len(c.Instructions)) {
		return nil, fmt.Errorf("no instruction at address %d", address)
	}

	return &c.Instructions[address], nil
}

func (c *Chunk) ConstantAt(address ConstAddress) (*value.Value, error) {
	if address >= uint32(len(c.Constants)) {
		return nil, fmt.Errorf("no constant at address %d", address)
	}

	return &c.Constants[address], nil
}

func (c *Chunk) IntrospectionInfoAt(address Address) (*SourceInformation, error) {
	return c.IntrospectionInfo.SourceInformationAt(address)
}

func (c *Chunk) Size() uint32 {
	return uint32(len(c.Instructions))
}
