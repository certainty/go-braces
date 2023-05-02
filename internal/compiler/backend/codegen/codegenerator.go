package codegen

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/frontend/ir"
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/certainty/go-braces/internal/isa"
)

type ConstantAddress uint64

type CodeUnitBuilder struct {
	// refine representation of constants later to allow deduplication
	constants        []isa.Value
	instructions     []isa.Instruction
	introspectionAPI introspection.API
}

func newCodeUnitBuilder(introspectionAPI introspection.API) *CodeUnitBuilder {
	return &CodeUnitBuilder{
		introspectionAPI: introspectionAPI,
		instructions:     make([]isa.Instruction, 0),
		constants:        make([]isa.Value, 0),
	}
}

func (c *CodeUnitBuilder) BuildCodeUnit() *isa.CodeUnit {
	return &isa.CodeUnit{
		Constants:    c.constants,
		Instructions: c.instructions,
	}
}

func (c *CodeUnitBuilder) AddConstant(constant *isa.Value) ConstantAddress {
	c.constants = append(c.constants, *constant)
	return ConstantAddress(len(c.constants) - 1)
}

func (c *CodeUnitBuilder) AddInstruction(instruction isa.Instruction) {
	c.instructions = append(c.instructions, instruction)
}

type Codegenerator struct {
	introspectionAPI              introspection.API
	registerAccu                  isa.Register
	generalPurposeRegisterOffset  uint8
	currentGeneralPurposeRegister uint32
}

func NewCodegenerator(introspectionAPI introspection.API) *Codegenerator {
	return &Codegenerator{
		introspectionAPI:              introspectionAPI,
		registerAccu:                  isa.Register(0),
		generalPurposeRegisterOffset:  16,
		currentGeneralPurposeRegister: 16,
	}
}

func (c *Codegenerator) NextRegister() isa.Register {
	c.currentGeneralPurposeRegister += 1
	return isa.Register(c.currentGeneralPurposeRegister)
}

func (c *Codegenerator) GenerateModule(intermediate *ir.IR) (*isa.AssemblyModule, error) {
	codeBuilder := newCodeUnitBuilder(c.introspectionAPI)

	for _, block := range intermediate.Blocks {
		if err := c.emitBlock(block, codeBuilder); err != nil {
			return nil, err
		}
	}
	codeBuilder.AddInstruction(isa.InstHalt(c.registerAccu))

	module := isa.NewAssemblyModule(
		isa.NewAssemblyMeta("", isa.AssemblyTypeExecutable),
		codeBuilder.BuildCodeUnit(),
		[]isa.ClosureValue{},
		nil,
		nil,
	)

	return module, nil
}

func (c *Codegenerator) emitBlock(block *ir.IRBlock, builder *CodeUnitBuilder) error {
	for _, instruction := range block.Instructions {
		switch instruction.(type) {
		case ir.IRConstant:
			// include register to use in IR?
			builder.AddInstruction(isa.InstTrue(c.registerAccu))
		default:
			return fmt.Errorf("unknown instruction type: %T", instruction)
		}
	}
	return nil
}
