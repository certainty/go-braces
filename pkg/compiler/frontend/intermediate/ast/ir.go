package ir

//TODO: this does too much as we're already lowering into something resembling SSA
// We will introduce an IR that uses the original order of operations
// Later we'll turn it into SSA (CFG)

import (
	"fmt"
	hlast "github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/ast"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	hltypes "github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/types"
	"github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/types"
	"log"
)

type Label string
type Register uint

type Callable interface {
	aCallable()
}

type Module struct {
	Name       Label
	Source     token.Origin
	Procedures []Procedure
}

func (m Module) String() string {
	writer := NewIRWriter(&m)
	return writer.Write()
}

type Procedure struct {
	tpe    types.Type
	Name   Label
	Args   []Argument
	Blocks []*BasicBlock
}

func (Procedure) aCallable() {}

type Argument struct {
	tpe      types.Type
	Register Register
}

type BasicBlock struct {
	Label        Label
	Instructions []Instruction
}

func (b BasicBlock) IsEmpty() bool {
	return len(b.Instructions) == 0
}

func (b BasicBlock) LastInstruction() Instruction {
	if b.IsEmpty() {
		return nil
	}

	return b.Instructions[len(b.Instructions)-1]
}

type Instruction interface {
	Type() types.Type
}

// %register = op tpe operand1, operand2, ...
type SimpleInstruction struct {
	tpe       types.Type
	Register  Register
	Operation Operation
	Operands  []Operand
}

func (i SimpleInstruction) Type() types.Type {
	return i.tpe
}

var _ Instruction = (*SimpleInstruction)(nil)

type ReturnInstruction struct {
	tpe      types.Type
	Register Register
}

func (i ReturnInstruction) Type() types.Type {
	return i.tpe
}

var _ Instruction = (*ReturnInstruction)(nil)

type AssignmentInstruction struct {
	tpe      types.Type
	Register Register
	Operand  Operand
}

func (i AssignmentInstruction) Type() types.Type {
	return i.tpe
}

var _ Instruction = (*AssignmentInstruction)(nil)

type Operation uint8
type Operand interface{}
type Literal interface{}

var _ Operand = Literal(nil)
var _ Operand = Label("")
var _ Operand = Register(0)

const (
	Add Operation = iota
	Sub
	Mul
	Div
	Or
	And
	Xor
	Neg
)

type RegisterAllocator struct {
	count     int
	registers map[string]Register
}

func NewRegisterAllocator() *RegisterAllocator {
	return &RegisterAllocator{
		count: 0,
	}
}

func (r RegisterAllocator) Last() Register {
	return Register(r.count)
}

func (r *RegisterAllocator) Next(variableName string) Register {
	r.count++
	if variableName != "" {
		r.registers[variableName] = Register(r.count)
	}
	return Register(r.count)
}

type IrBuilder struct {
	typeUniverse hltypes.TypeUniverse
	Module       *Module
}

func NewBuilder(origin token.Origin, tpeUniverse hltypes.TypeUniverse) *IrBuilder {
	return &IrBuilder{
		Module:       CreateModule("", origin),
		typeUniverse: tpeUniverse,
	}
}

func LowerToIR(origin token.Origin, theAst *hlast.Source, tpeUniverse hltypes.TypeUniverse) (*Module, error) {
	builder := NewBuilder(origin, tpeUniverse)

	if err := builder.lower(theAst); err != nil {
		return nil, err
	}
	return builder.Module, nil
}

func (b *IrBuilder) blockBuilder(label string, registers *RegisterAllocator) *BlockBuilder {
	return NewBlockBuilder(Label(label), registers, b)
}

func (b *IrBuilder) lower(theAst *hlast.Source) error {
	for _, node := range theAst.Declarations {
		switch node := node.(type) {
		case hlast.ProcDecl:
			proc, err := b.lowerProcedure(node)
			if err != nil {
				return err
			}
			b.Module.Procedures = append(b.Module.Procedures, proc)
		default:
			return fmt.Errorf("unexpected node type: %T", node)
		}
	}
	return nil
}

func (b *IrBuilder) lowerStatement(node hlast.Statement, blockBuilder *BlockBuilder) (Register, error) {
	switch node := node.(type) {
	case hlast.ExprStmt:
		reg, _, err := b.lowerExpression(node.Expr, blockBuilder)
		if err != nil {
			return 0, err
		}
		return reg, nil
	default:
		return 0, fmt.Errorf("unexpected node type: %T", node)
	}
}

func (b *IrBuilder) lowerExpression(node hlast.Expression, blockBuilder *BlockBuilder) (Register, types.Type, error) {
	switch node := node.(type) {
	case hlast.BasicLitExpr:
		exprType, err := b.typeOf(node)
		if err != nil {
			return 0, nil, err
		}
		loweredType, err := b.lowerType(exprType)
		if err != nil {
			return 0, nil, err
		}
		return blockBuilder.OpLit(loweredType, node.Value), loweredType, nil
	case hlast.BinaryExpr:
		return b.lowerBinaryExpression(blockBuilder, node)
	default:
		return 0, nil, fmt.Errorf("unexpected node type: %T", node)
	}

}

func (b *IrBuilder) lowerProcedure(decl hlast.ProcDecl) (Procedure, error) {
	var err error
	funType, err := b.typeOf(decl)
	if err != nil {
		return Procedure{}, err
	}
	loweredType, err := b.lowerType(funType)
	if err != nil {
		return Procedure{}, err
	}
	proc := CreateProcedure(loweredType, Label(decl.Name.Name))
	procRegisters := NewRegisterAllocator()

	blockBuilder := b.blockBuilder("entry", procRegisters)
	for _, stmt := range decl.Body.Statements {
		b.lowerStatement(stmt, blockBuilder)
	}

	if blockBuilder.IsEmpty() {
		reg := blockBuilder.OpLit(types.UnitType, Literal(nil))
		blockBuilder.OpRet(types.UnitType, reg)
	} else {
		lastInst := blockBuilder.LastInstruction()
		blockBuilder.OpRet(lastInst.Type(), procRegisters.Last())
	}

	proc.Blocks = append(proc.Blocks, blockBuilder.Block)
	return proc, nil
}

func (b *IrBuilder) lowerBinaryExpression(builder *BlockBuilder, expr hlast.BinaryExpr) (Register, types.Type, error) {
	exprType, err := b.typeOf(expr)
	if err != nil {
		return 0, nil, err
	}
	loweredType, err := b.lowerType(exprType)
	if err != nil {
		return 0, nil, err
	}

	leftReg, leftTpe, err := b.lowerExpression(expr.Left, builder)
	if err != nil {
		return 0, nil, err
	}
	rightReg, _, err := b.lowerExpression(expr.Right, builder)
	if err != nil {
		return 0, nil, err
	}

	switch expr.Op.Type {
	case token.ADD:
		return builder.OpAdd(loweredType, leftReg, rightReg), leftTpe, nil
	case token.MUL:
		return builder.OpMul(loweredType, leftReg, rightReg), leftTpe, nil
	default:
		return 0, nil, fmt.Errorf("unexpected binary operator: %v", expr.Op)
	}
}

func (b *IrBuilder) typeOf(node hlast.Node) (hltypes.Type, error) {
	log.Printf("typeuniverse: %v", b.typeUniverse)
	tpe, ok := b.typeUniverse[node.ID()]
	if !ok {
		return nil, fmt.Errorf("no type found for node: %T", node)
	}
	return tpe, nil
}

func (b *IrBuilder) lowerType(t hltypes.Type) (types.Type, error) {
	switch t.(type) {
	case hltypes.Byte:
		return types.ByteType, nil
	case hltypes.Int:
		return types.IntType, nil
	case hltypes.Float:
		return types.FloatType, nil
	case hltypes.Bool:
		return types.BoolType, nil
	case hltypes.String:
		return types.StringType, nil
	default:
		return nil, fmt.Errorf("unexpected type: %s", t)
	}
}
