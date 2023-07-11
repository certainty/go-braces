package codegen

import (
	"fmt"
	"log"

	"github.com/certainty/go-braces/internal/compiler/frontend/ir"
	"github.com/certainty/go-braces/internal/compiler/frontend/lexer"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/isa"
	"github.com/certainty/go-braces/internal/isa/arity"
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
	module                        *isa.AssemblyModule
	functionAddresses             map[isa.Label]isa.Address
	generalPurposeRegisterOffset  uint8
	currentGeneralPurposeRegister uint32
}

func NewCodegenerator(instrumentation compiler_introspection.Instrumentation) *Codegenerator {
	mod := isa.NewAssemblyModule(isa.NewAssemblyMeta("main", isa.AssemblyTypeExecutable))

	return &Codegenerator{
		instrumentation:               instrumentation,
		registerAccu:                  isa.Register(isa.REG_SP_ACCU),
		generalPurposeRegisterOffset:  16,
		currentGeneralPurposeRegister: 16,
		functionAddresses:             make(map[isa.Label]isa.Address, 0),
		module:                        mod,
	}
}

func (c *Codegenerator) NextRegister() isa.Register {
	c.currentGeneralPurposeRegister += 1
	return isa.Register(c.currentGeneralPurposeRegister)
}

func (c *Codegenerator) GenerateModule(irModule *ir.Module) (*isa.AssemblyModule, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseCodegen)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseCodegen)

	for _, fun := range irModule.Functions {
		if err := c.emitFunction(fun); err != nil {
			return nil, err
		}
	}

	c.module.EntryPoint = 0 // FIXME: set to main
	return c.module, nil
}

// we should return the function address to fill the jump table
func (c *Codegenerator) emitFunction(fun *ir.Function) error {
	codeBuilder := newCodeUnitBuilder(c.instrumentation)
	for _, block := range fun.Blocks {
		if err := c.emitBlock(block, codeBuilder); err != nil {
			return err
		}
	}
	c.addFunction(isa.Label(fun.Name), arity.Exactly(0), *codeBuilder.BuildCodeUnit())
	return nil
}

func (c *Codegenerator) addFunction(label isa.Label, theArity arity.Arity, code isa.CodeUnit) {
	c.module.Functions = append(c.module.Functions, isa.Function{
		Label: label,
		Arity: theArity,
		Code:  code,
	})
	c.functionAddresses[label] = isa.Address(len(c.module.Functions) - 1)
}

func (c *Codegenerator) emitBlock(block *ir.BasicBlock, builder *CodeUnitBuilder) error {
	for _, instruction := range block.Instructions {
		log.Printf("emitBlock: %s", instruction)
		switch inst := instruction.(type) {
		case ir.ReturnInstruction:
			builder.AddInstruction(isa.InstRet(c.findRegister(inst.Register)))
		case ir.SimpleInstruction:
			log.Printf("emitSimple: %T", inst)
			if err := c.emitSimpleInstruction(inst, builder); err != nil {
				return fmt.Errorf("emitSimpleInstruction: %w", err)
			}
		case ir.Literal:
			value, err := c.convertValue(inst)
			if err != nil {
				return fmt.Errorf("emitBlock: %w", err)
			}
			address := builder.AddConstant(&value)
			builder.AddInstruction(isa.InstLoad(address, c.registerAccu))
		default:
			return fmt.Errorf("unknown instruction type: %T", instruction)
		}
	}
	return nil
}

func (c *Codegenerator) emitSimpleInstruction(inst ir.SimpleInstruction, builder *CodeUnitBuilder) error {
	switch inst.Operation {
	case ir.Add:
		left := inst.Operands[0].(ir.Register)
		right := inst.Operands[1]

		switch right := right.(type) {
		case ir.Literal:
			switch v := right.(type) {
			case uint8:
			case int:
				if v <= 255 {
					builder.AddInstruction(isa.InstAddI(isa.Register(left), isa.ImmediateValue(v), c.findRegister(inst.Register)))
				} else {
					builder.AddInstruction(isa.InstAdd(isa.Register(left), isa.Register(v), c.findRegister(inst.Register)))
				}
			}
		case ir.Register:
			builder.AddInstruction(isa.InstAdd(isa.Register(left), isa.Register(right), c.findRegister(inst.Register)))
		}

	case ir.Sub:
		builder.AddInstruction(isa.InstSub(c.findRegister(inst.Register), c.registerAccu, c.registerAccu))
	default:
		panic("unknown instruction")
	}
	return nil
}

func (c *Codegenerator) findRegister(reg ir.Register) isa.Register {
	return isa.Register(reg)
}

func (c *Codegenerator) convertValue(v interface{}) (isa.Value, error) {
	switch v := v.(type) {
	case bool:
		return isa.Bool(v), nil
	case lexer.CodePoint:
		return isa.Char(v.Char), nil
	case int:
		return isa.Int(v), nil
	case uint:
		return isa.UInt(v), nil
	case float32, float64:
		return isa.Float(v.(float64)), nil
	case string:
		return isa.String(v), nil

	default:
		return nil, fmt.Errorf("unknown value type: %T", v)
	}
}
