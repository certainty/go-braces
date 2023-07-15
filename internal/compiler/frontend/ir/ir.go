package ir

import (
	"fmt"
	"log"

	"github.com/certainty/go-braces/internal/compiler/frontend/ast/hl"
	"github.com/certainty/go-braces/internal/compiler/frontend/token"
	"github.com/certainty/go-braces/internal/compiler/frontend/types"
)

type Label string
type Register uint

type Module struct {
	Name      Label
	Source    token.Origin
	Functions []*Function
}

func (m Module) String() string {
	writer := NewIRWriter(&m)
	return writer.Write()
}

type Function struct {
	tpe    Type
	Name   Label
	Args   []Argument
	Blocks []*BasicBlock
}

type Argument struct {
	tpe      Type
	Register Register
}

type BasicBlock struct {
	Label        Label
	Instructions []Instruction
}

type Instruction interface {
	Type() Type
}

// %register = op tpe operand1, operand2, ...
type SimpleInstruction struct {
	tpe       Type
	Register  Register
	Operation Operation
	Operands  []Operand
}

func (i SimpleInstruction) Type() Type {
	return i.tpe
}

var _ Instruction = (*SimpleInstruction)(nil)

type ReturnInstruction struct {
	tpe      Type
	Register Register
}

func (i ReturnInstruction) Type() Type {
	return i.tpe
}

var _ Instruction = (*ReturnInstruction)(nil)

type AssignmentInstruction struct {
	tpe      Type
	Register Register
	Operand  Operand
}

func (i AssignmentInstruction) Type() Type {
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

func (r *RegisterAllocator) Next(variableName string) Register {
	r.count++
	if variableName != "" {
		r.registers[variableName] = Register(r.count)
	}
	return Register(r.count)
}

type IrBuilder struct {
	typeUniverse types.TypeUniverse
	Module       *Module
}

func NewBuilder(origin token.Origin, tpeUniverse types.TypeUniverse) *IrBuilder {
	return &IrBuilder{
		Module:       CreateModule("", origin),
		typeUniverse: tpeUniverse,
	}
}

func LowerToIR(origin token.Origin, theAst *ast.Source, tpeUniverse types.TypeUniverse) (*Module, error) {
	builder := NewBuilder(origin, tpeUniverse)

	if err := builder.lower(theAst); err != nil {
		return nil, err
	}
	return builder.Module, nil
}

func (b *IrBuilder) blockBuilder(label string, registers *RegisterAllocator) *BlockBuilder {
	return NewBlockBuilder(Label(label), registers, b)
}

func (b *IrBuilder) lower(theAst *ast.Source) error {
	for _, node := range theAst.Statements {
		switch node := node.(type) {
		case ast.Bl:
			fun, err := b.lowerFunction(node)
			if err != nil {
				return err
			}
			b.Module.Functions = append(b.Module.Functions, fun)
		default:
			return fmt.Errorf("unexpected node type: %T", node)
		}
	}
	return nil
}

func (b *IrBuilder) lowerBody(node ast.Node, blockBuilder *BlockBuilder) (Register, Type, error) {
	switch node := node.(type) {
	case ast.BasicLitExpr:
		exprType, err := b.typeOf(node)
		if err != nil {
			return 0, nil, err
		}
		loweredType, err := b.lowerType(exprType)
		if err != nil {
			return 0, nil, err
		}
		return blockBuilder.OpLit(loweredType, node.Value), loweredType, nil
	case ast.BinaryExpr:
		return b.lowerBinaryExpression(blockBuilder, node)
	default:
		return 0, nil, fmt.Errorf("unexpected node type: %T", node)
	}
}

func (b *IrBuilder) lowerFunction(decl ast.CallableDecl) (*Function, error) {
	var err error
	funType, err := b.typeOf(decl)
	if err != nil {
		return nil, err
	}
	loweredType, err := b.lowerType(funType)
	fun := CreateFunction(loweredType, Label(decl.Name.Label))
	funRegisters := NewRegisterAllocator()

	for _, arg := range decl.Arguments {
		declType, err := b.typeOf(decl)
		tpe, err := b.lowerType(declType)
		if err != nil {
			return nil, err
		}

		register := funRegisters.Next(arg.Name.Label)
		fun.Args = append(fun.Args, Argument{tpe: tpe, Register: register})
	}

	blockBuilder := b.blockBuilder("entry", funRegisters)

	var returnRegister Register
	var tpe Type

	for _, stmt := range decl.Body.Code {
		returnRegister, tpe, err = b.lowerBody(stmt, blockBuilder)
		if err != nil {
			return nil, err
		}
	}

	blockBuilder.OpRet(tpe, returnRegister)

	fun.Blocks = append(fun.Blocks, blockBuilder.Block)
	return fun, nil
}

func (b *IrBuilder) lowerBinaryExpression(builder *BlockBuilder, expr ast.BinaryExpression) (Register, Type, error) {
	exprType, err := b.typeOf(expr)
	if err != nil {
		return 0, nil, err
	}
	loweredType, err := b.lowerType(exprType)
	if err != nil {
		return 0, nil, err
	}

	leftReg, leftTpe, err := b.lowerBody(expr.Left, builder)
	if err != nil {
		return 0, nil, err
	}
	rightReg, _, err := b.lowerBody(expr.Right, builder)
	if err != nil {
		return 0, nil, err
	}

	switch expr.Operator {
	case ast.BinOpAdd:
		return builder.OpAdd(loweredType, leftReg, rightReg), leftTpe, nil
	case ast.BinOpMul:
		return builder.OpMul(loweredType, leftReg, rightReg), leftTpe, nil
	default:
		return 0, nil, fmt.Errorf("unexpected binary operator: %v", expr.Operator)
	}
}

func (b *IrBuilder) typeOf(node ast.Node) (types.Type, error) {
	log.Printf("typeuniverse: %v", b.typeUniverse)
	tpe, ok := b.typeUniverse[node.ID()]
	if !ok {
		return nil, fmt.Errorf("no type found for node: %T", node)
	}
	return tpe, nil
}

func (b *IrBuilder) lowerType(t types.Type) (Type, error) {
	switch t.(type) {
	case types.Byte:
		return ByteType, nil
	case types.Int:
		return IntType, nil
	case types.Float:
		return FloatType, nil
	case types.Bool:
		return BoolType, nil
	case types.String:
		return StringType, nil
	default:
		return nil, fmt.Errorf("unexpected type: %s", t)
	}
}
