package bytecode

import "github.com/certainty/go-braces/vm/internal/language/value"

type ChunkBuilder struct {
	instructions      []Instruction
	constants         map[value.Value]ConstAddress
	constIndex        uint32
	sourceInformation []SourceInformation
}

func NewChunkBuilder() *ChunkBuilder {
	return &ChunkBuilder{
		instructions:      make([]Instruction, 0),
		constants:         make(map[value.Value]ConstAddress),
		constIndex:        0,
		sourceInformation: make([]SourceInformation, 0),
	}
}

func (cb *ChunkBuilder) AddInstruction(instruction Instruction, source SourceInformation) Address {
	cb.instructions = append(cb.instructions, instruction)
	cb.sourceInformation = append(cb.sourceInformation, source)
	return Address(len(cb.instructions) - 1)
}

func (cb *ChunkBuilder) AddConstant(constant value.Value) ConstAddress {
	if address, ok := cb.constants[constant]; ok {
		return address
	} else {
		cb.constants[constant] = cb.constIndex
		cb.constIndex++
		return cb.constIndex - 1
	}
}

func (cb *ChunkBuilder) Build() *Chunk {
	constants := make([]value.Value, len(cb.constants))

	for value, address := range cb.constants {
		constants[address] = value
	}

	return &Chunk{
		Instructions:      cb.instructions,
		Constants:         constants,
		IntrospectionInfo: &IntrospectionInfo{sourceInformation: cb.sourceInformation},
		DebugInfo:         nil,
	}
}
