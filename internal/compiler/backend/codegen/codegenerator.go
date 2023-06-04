package codegen

import (
	"fmt"
	"log"

	"github.com/certainty/go-braces/internal/compiler/frontend/ir"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/isa"
)

type CodeUnitBuilder struct {
	// refine representation of constants later to allow deduplication
	constants       []isa.Value
	instructions    []isa.Instruction
	instrumentation compiler_introspection.Instrumentation
}

func newCodeUnitBuilder(instrumentation compiler_introspection.Instrumentation) *CodeUnitBuilder {
	return &CodeUnitBuilder{
		instrumentation: instrumentation,
		instructions:    make([]isa.Instruction, 0),
		constants:       make([]isa.Value, 0),
	}
}

func (c *CodeUnitBuilder) BuildCodeUnit() *isa.CodeUnit {
	return &isa.CodeUnit{
		Constants:    c.constants,
		Instructions: c.instructions,
	}
}

func (c *CodeUnitBuilder) AddConstant(constant *isa.Value) isa.ConstantAddress {
	c.constants = append(c.constants, *constant)
	return isa.ConstantAddress(len(c.constants) - 1)
}

func (c *CodeUnitBuilder) AddInstruction(instruction isa.Instruction) {
	c.instructions = append(c.instructions, instruction)
}

type Codegenerator struct {
	instrumentation               compiler_introspection.Instrumentation
	registerAccu                  isa.Register
	generalPurposeRegisterOffset  uint8
	currentGeneralPurposeRegister uint32
}

func NewCodegenerator(instrumentation compiler_introspection.Instrumentation) *Codegenerator {
	return &Codegenerator{
		instrumentation:               instrumentation,
		registerAccu:                  isa.Register(isa.REG_SP_ACCU),
		generalPurposeRegisterOffset:  16,
		currentGeneralPurposeRegister: 16,
	}
}

func (c *Codegenerator) NextRegister() isa.Register {
	c.currentGeneralPurposeRegister += 1
	return isa.Register(c.currentGeneralPurposeRegister)
}

func (c *Codegenerator) GenerateModule(intermediate *ir.IR) (*isa.AssemblyModule, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseCodegen)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseCodegen)

	codeBuilder := newCodeUnitBuilder(c.instrumentation)

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
			value := instruction.(ir.IRConstant).Value
			switch value.(type) {
			case isa.BoolValue:
				if value == isa.BoolValue(true) {
					log.Printf("emitBlock: true")
					builder.AddInstruction(isa.InstTrue(c.registerAccu))
				} else {
					log.Printf("emitBlock: false")
					builder.AddInstruction(isa.InstFalse(c.registerAccu))
				}
			case isa.CharValue:
				log.Printf("emitBlock: Constant(char)")
				address := builder.AddConstant(&value)
				builder.AddInstruction(isa.InstConst(address, c.registerAccu))
			default:
				return fmt.Errorf("unknown constant type: %T", value)
			}
		default:
			return fmt.Errorf("unknown instruction type: %T", instruction)
		}
	}
	return nil
}
