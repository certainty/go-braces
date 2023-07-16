package codegen

import (
	"fmt"
	"log"

	"github.com/certainty/go-braces/internal/compiler/frontend/ir"
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

func (c *CodeUnitBuilder) AddConstant(constant isa.Value) isa.ConstantAddress {
	c.constants = append(c.constants, constant)
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

	for _, proc := range irModule.Procedures {
		if err := c.emitProcedure(&proc); err != nil {
			return nil, err
		}
	}

	c.module.EntryPoint = 0 // FIXME: set to main
	return c.module, nil
}

// we should return the function address to fill the jump table
func (c *Codegenerator) emitProcedure(proc *ir.Procedure) error {
	codeBuilder := newCodeUnitBuilder(c.instrumentation)
	for _, block := range proc.Blocks {
		if err := c.emitBlock(block, codeBuilder); err != nil {
			return err
		}
	}
	c.addFunction(isa.Label(proc.Name), arity.Exactly(0), *codeBuilder.BuildCodeUnit())
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
		switch inst := instruction.(type) {
		case ir.ReturnInstruction:
			builder.AddInstruction(isa.InstRet(c.findRegister(inst.Register)))

		case ir.SimpleInstruction:
			if err := c.emitSimpleInstruction(inst, builder); err != nil {
				return fmt.Errorf("emitSimpleInstruction: %w", err)
			}
		case ir.AssignmentInstruction:
			addr := builder.AddConstant(isa.Value(inst.Operand))
			reg := c.findRegister(inst.Register)
			builder.AddInstruction(isa.InstLoad(reg, addr))
		case ir.Literal:
			value, err := c.convertValue(inst)
			if err != nil {
				return fmt.Errorf("emitBlock: %w", err)
			}
			address := builder.AddConstant(value)
			builder.AddInstruction(isa.InstLoad(c.registerAccu, address))
		default:
			return fmt.Errorf("unknown instruction type: %T", instruction)
		}
	}
	return nil
}

func (c *Codegenerator) emitSimpleInstruction(inst ir.SimpleInstruction, builder *CodeUnitBuilder) error {
	switch inst.Operation {
	case ir.Add, ir.Sub, ir.Mul:
		return c.emitArithmeticInstruction(inst, builder)
	default:
		log.Printf("unknown instruction: %v", inst.Operation)
		panic("unknown instruction")
	}
}

// TODO: add support for immediate values
func (c *Codegenerator) emitArithmeticInstruction(inst ir.SimpleInstruction, builder *CodeUnitBuilder) error {
	if len(inst.Operands) != 2 {
		return fmt.Errorf("expected 2 operands, got %d", len(inst.Operands))
	}

	var left, right isa.Register

	if c.isImmediate(inst.Operands[0]) {
		addr := builder.AddConstant(isa.Value(inst.Operands[0].(ir.Literal)))
		left = c.NextRegister()
		builder.AddInstruction(isa.InstLoad(left, addr))
	} else {
		left = c.findRegister(inst.Operands[0].(ir.Register))
	}

	if c.isImmediate(inst.Operands[1]) {
		addr := builder.AddConstant(isa.Value(inst.Operands[1].(ir.Literal)))
		right = c.NextRegister()
		builder.AddInstruction(isa.InstLoad(right, addr))
	} else {
		right = c.findRegister(inst.Operands[1].(ir.Register))
	}

	switch inst.Operation {
	case ir.Add:
		builder.AddInstruction(isa.InstAdd(c.findRegister(inst.Register), left, right))
	case ir.Mul:
		builder.AddInstruction(isa.InstMul(c.findRegister(inst.Register), left, right))
	case ir.Sub:
		builder.AddInstruction(isa.InstSub(c.findRegister(inst.Register), left, right))
	}
	return nil
}

func (c *Codegenerator) isImmediate(op ir.Operand) bool {
	if op, ok := op.(ir.Literal); ok {
		if v, ok := op.(int); ok {
			if v <= 255 {
				return true
			}
		}
	}
	return false
}

func (c *Codegenerator) findRegister(reg ir.Register) isa.Register {
	return isa.Register(reg)
}

func (c *Codegenerator) convertValue(v interface{}) (isa.Value, error) {
	switch v := v.(type) {
	case bool:
		return isa.Bool(v), nil
	case rune:
		return isa.Char(v), nil
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
