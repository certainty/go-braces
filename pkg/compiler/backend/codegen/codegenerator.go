package codegen

import (
	"fmt"

	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	"github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/ssa"
	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
	"github.com/certainty/go-braces/pkg/shared/isa"
	"github.com/certainty/go-braces/pkg/shared/isa/arity"
)

type CodeUnitBuilder struct {
	// refine representation of constants later to allow deduplication
	constants       []isa.Value
	instructions    []isa.Instruction
	instrumentation compiler_introspection.Instrumentation
	registers       uint32
}

func newCodeUnitBuilder(instrumentation compiler_introspection.Instrumentation) *CodeUnitBuilder {
	return &CodeUnitBuilder{
		instrumentation: instrumentation,
		instructions:    make([]isa.Instruction, 0),
		constants:       make([]isa.Value, 0),
		registers:       0,
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

func (c *CodeUnitBuilder) NextRegister() isa.Register {
	c.registers++
	return isa.Register(c.registers)
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

func (c *Codegenerator) GenerateModule(ssaModule *ssa.Module) (*isa.AssemblyModule, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseCodegen)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseCodegen)

	for _, decl := range ssaModule.Declarations {

		if err := c.emitDeclaration(decl); err != nil {
			return nil, err
		}
	}

	c.module.EntryPoint = 0 // FIXME: set to main
	return c.module, nil
}

func (c *Codegenerator) emitDeclaration(decl ssa.Declaration) error {
	switch decl := decl.(type) {
	case ssa.ProcDecl:
		return c.emitProcedure(&decl)
	default:
		return fmt.Errorf("unknown declaration type %T", decl)
	}
}

// we should return the function address to fill the jump table
func (c *Codegenerator) emitProcedure(procDecl *ssa.ProcDecl) error {
	codeBuilder := newCodeUnitBuilder(c.instrumentation)
	for _, block := range procDecl.Blocks {
		if err := c.emitBlock(block, codeBuilder); err != nil {
			return err
		}
	}
	c.addFunction(isa.Label(procDecl.Name), arity.Exactly(0), *codeBuilder.BuildCodeUnit())
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

func (c *Codegenerator) emitBlock(block *ssa.BasicBlock, builder *CodeUnitBuilder) error {
	for _, statement := range block.Statements {
		switch stmt := statement.(type) {
		case ssa.ReturnStmt:
			c.emitReturn(stmt, builder)
		case ssa.SetStmt:
			c.emitAssignment(stmt, builder)
		default:
			return fmt.Errorf("unknown instruction type: %T", statement)
		}
	}
	return nil
}

func (c *Codegenerator) emitReturn(stmt ssa.ReturnStmt, builder *CodeUnitBuilder) error {
	switch v := stmt.Value.(type) {
	case ssa.AtomicLitExpr:
		reg, error := c.emitLiteral(v, builder)
		if error != nil {
			return error
		}
		builder.AddInstruction(isa.InstRet(*reg))
		return nil
	case ssa.Variable:
		builder.AddInstruction(isa.InstRet(c.findRegister(v)))
		return nil
	default:
		return fmt.Errorf("unknown return type: %T", stmt.Value)
	}
}

func (c *Codegenerator) emitAssignment(stmt ssa.SetStmt, builder *CodeUnitBuilder) {
	addr := builder.AddConstant(isa.Value(stmt.Value))
	reg := c.findRegister(stmt.Variable)
	builder.AddInstruction(isa.InstLoad(reg, addr))
}

func (c *Codegenerator) emitExpression(expr ssa.Expression, builder *CodeUnitBuilder) (*isa.Register, error) {
	switch expr := expr.(type) {
	case ssa.AtomicLitExpr:
		return c.emitLiteral(expr, builder)
	case ssa.BinaryExpr:
		return c.emitBinaryExpression(expr, builder)
	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (c *Codegenerator) emitLiteral(expr ssa.AtomicLitExpr, builder *CodeUnitBuilder) (*isa.Register, error) {
	value := expr.IrExpr.HlExpr.Token.LitValue
	if value == nil {
		return nil, fmt.Errorf("nil return value")
	}

	convertedValue, err := c.convertValue(value)
	if err != nil {
		return nil, err
	}
	addr := builder.AddConstant(convertedValue)
	reg := builder.NextRegister()
	builder.AddInstruction(isa.InstLoad(reg, addr))
	return &reg, nil
}

func (c *Codegenerator) emitBinaryExpression(expr ssa.BinaryExpr, builder *CodeUnitBuilder) (*isa.Register, error) {
	left := c.findRegister(expr.Left)
	right := c.findRegister(expr.Right)
	reg := builder.NextRegister()

	switch expr.IrExpr.Op.Type {
	case token.ADD:
		builder.AddInstruction(isa.InstAdd(reg, left, right))
	case token.SUB:
		builder.AddInstruction(isa.InstSub(reg, left, right))
	case token.MUL:
		builder.AddInstruction(isa.InstMul(reg, left, right))
	default:
		return nil, fmt.Errorf("unknown binary expression type: %T", expr)
	}

	return &reg, nil
}

func (c *Codegenerator) findRegister(v ssa.Variable) isa.Register {
	return isa.Register(v.Version)
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
