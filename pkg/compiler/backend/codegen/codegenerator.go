package codegen

import (
	"fmt"

	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	ir "github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/ast"
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

func newCodeUnitBuilder(instrumentation compiler_introspection.Instrumentation, registers uint32) *CodeUnitBuilder {
	return &CodeUnitBuilder{
		instrumentation: instrumentation,
		instructions:    make([]isa.Instruction, 0),
		constants:       make([]isa.Value, 0),
		registers:       registers,
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
	registers                     map[string]isa.Register
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
		registers:                     make(map[string]isa.Register, 0),
	}
}

func (c *Codegenerator) GenerateModule(ssaModule *ir.Module) (*isa.AssemblyModule, error) {
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

func (c *Codegenerator) emitDeclaration(decl ir.Declaration) error {
	switch decl := decl.(type) {
	case *ir.ProcDecl:
		return c.emitProcedure(decl)
	default:
		return fmt.Errorf("unknown declaration type %T", decl)
	}
}

// we should return the function address to fill the jump table
func (c *Codegenerator) emitProcedure(procDecl *ir.ProcDecl) error {
	codeBuilder := newCodeUnitBuilder(c.instrumentation, c.currentGeneralPurposeRegister)

	if procDecl.SSABlocks == nil {
		return fmt.Errorf("procedure %s has no SSA blocks", procDecl.Name.Value)
	}

	for _, block := range procDecl.SSABlocks {
		if err := c.emitBlock(block, codeBuilder); err != nil {
			return err
		}
	}

	c.addFunction(isa.Label(procDecl.Name.Value), arity.Exactly(0), *codeBuilder.BuildCodeUnit())
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
	for _, statement := range block.Statements {
		switch stmt := statement.(type) {
		case *ir.ReturnStmt:
			err := c.emitReturn(stmt, builder)
			if err != nil {
				return nil
			}
		case *ir.AssignStmt:
			err := c.emitAssignment(stmt, builder)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown instruction type: %T", statement)
		}
	}
	return nil
}

func (c *Codegenerator) emitReturn(stmt *ir.ReturnStmt, builder *CodeUnitBuilder) error {
	switch v := stmt.Value.(type) {
	case *ir.AtomicLitExpr:
		reg, error := c.emitLiteral(v, builder)
		if error != nil {
			return error
		}
		builder.AddInstruction(isa.InstRet(reg))
		return nil
	case *ir.Variable:
		builder.AddInstruction(isa.InstRet(c.findRegister(v)))
		return nil
	default:
		return fmt.Errorf("unknown return type: %T", stmt.Value)
	}
}

func (c *Codegenerator) emitAssignment(stmt *ir.AssignStmt, builder *CodeUnitBuilder) error {
	reg, err := c.emitExpression(stmt.Expr, builder)
	if err != nil {
		return err
	}
	target := c.findRegister(stmt.Variable)
	builder.AddInstruction(isa.InstStore(target, reg))

	return nil
}

func (c *Codegenerator) emitExpression(expr ir.Expression, builder *CodeUnitBuilder) (isa.Register, error) {
	switch expr := expr.(type) {
	case *ir.AtomicLitExpr:
		return c.emitLiteral(expr, builder)
	case *ir.BinaryExpr:
		return c.emitBinaryExpression(expr, builder)
	default:
		return 0, fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (c *Codegenerator) emitLiteral(expr *ir.AtomicLitExpr, builder *CodeUnitBuilder) (isa.Register, error) {
	value := expr.Value.LitValue

	if value == nil {
		return 0, fmt.Errorf("nil return value")
	}

	convertedValue, err := c.convertValue(value)
	if err != nil {
		return 0, err
	}
	addr := builder.AddConstant(convertedValue)
	reg := builder.NextRegister()
	builder.AddInstruction(isa.InstLoad(reg, addr))
	return reg, nil
}

func (c *Codegenerator) emitBinaryExpression(expr *ir.BinaryExpr, builder *CodeUnitBuilder) (isa.Register, error) {
	left, err := c.emitBinaryExprOperand(expr.Left, builder)
	if err != nil {
		return 0, err
	}
	right, err := c.emitBinaryExprOperand(expr.Right, builder)
	if err != nil {
		return 0, err
	}
	reg := builder.NextRegister()

	switch expr.Op.Type {
	case token.ADD:
		builder.AddInstruction(isa.InstAdd(reg, left, right))
	case token.SUB:
		builder.AddInstruction(isa.InstSub(reg, left, right))
	case token.MUL:
		builder.AddInstruction(isa.InstMul(reg, left, right))
	default:
		return 0, fmt.Errorf("unknown binary expression type: %T", expr)
	}

	return reg, nil
}

func (c *Codegenerator) emitBinaryExprOperand(expr ir.Expression, builder *CodeUnitBuilder) (isa.Register, error) {
	switch expr := expr.(type) {
	case *ir.AtomicLitExpr:
		return c.emitLiteral(expr, builder)
	case *ir.Variable:
		return c.findRegister(expr), nil
	default:
		return 0, fmt.Errorf("unknown binary expression operand type: %T", expr)
	}
}

func (c *Codegenerator) findRegister(v *ir.Variable) isa.Register {
	varKey := variableKey(*v)
	reg, ok := c.registers[varKey]

	if !ok {
		reg = c.NextRegister()
		c.registers[varKey] = reg
		return reg
	} else {
		return reg
	}
}

func variableKey(v ir.Variable) string {
	return fmt.Sprintf("%s%d", v.Name, v.Version)
}

func (c *Codegenerator) NextRegister() isa.Register {
	c.currentGeneralPurposeRegister++
	return isa.Register(c.currentGeneralPurposeRegister)
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
