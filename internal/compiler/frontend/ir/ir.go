package ir

import (
	"fmt"
	"log"

	"github.com/certainty/go-braces/internal/compiler/frontend/parser/ast"
	"github.com/certainty/go-braces/internal/compiler/frontend/types"
	"github.com/certainty/go-braces/internal/compiler/location"
)

type Label string
type Register uint

type Module struct {
	Name      Label
	Source    location.Origin
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
	Module *Module
}

func NewBuilder(origin location.Origin) *IrBuilder {
	return &IrBuilder{
		Module: CreateModule("", origin),
	}
}

func LowerToIR(origin location.Origin, theAst *ast.AST) (*Module, error) {
	builder := NewBuilder(origin)
	if err := builder.lower(theAst); err != nil {
		return nil, err
	}
	return builder.Module, nil
}

func (b *IrBuilder) lower(theAst *ast.AST) error {
	for _, node := range theAst.Nodes {
		switch node := node.(type) {
		case ast.PackageDecl:
			b.Module.Name = Label(node.Name.ID)
		case ast.CallableDecl:
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

func (b *IrBuilder) lowerFunction(decl ast.CallableDecl) (*Function, error) {
	fun := CreateFunction(BoolType, Label(decl.Name.ID))
	registerAllocator := NewRegisterAllocator()

	// lower arguments
	for _, arg := range decl.Arguments {
		tpe, err := b.lowerType(translateType(arg.Type))
		if err != nil {
			return nil, err
		}
		register := registerAllocator.Next(arg.Name.ID)
		fun.Args = append(fun.Args, Argument{tpe: tpe, Register: register})
	}

	blockBuilder := NewBlockBuilder(Label("entry"), registerAllocator)

	for _, stmt := range decl.Body.Code {
		log.Printf("stmt: %T", stmt)
		switch stmt := stmt.(type) {
		case ast.BinaryExpression:
			b.lowerBinaryExpression(blockBuilder, stmt)
		default: // ignore
		}
	}

	if len(blockBuilder.Block.Instructions) > 0 {
		lastInstruction := blockBuilder.Block.Instructions[len(blockBuilder.Block.Instructions)-1]
		blockBuilder.OpRet(lastInstruction.Type(), lastInstruction.(SimpleInstruction).Register)
	}

	fun.Blocks = append(fun.Blocks, blockBuilder.Block)
	return fun, nil
}

func (b *IrBuilder) lowerBinaryExpression(builder *BlockBuilder, expr ast.BinaryExpression) {
	switch expr.Operator {
	case ast.BinOpAdd:
		leftOp := b.lowerOperand(builder, expr.Left)
		rightOp := b.lowerOperand(builder, expr.Right)
		builder.OpAdd(IntType, leftOp, rightOp)
	default:
		panic("not implemented")
	}
}

func (b *IrBuilder) lowerOperand(builder *BlockBuilder, expr ast.Expression) Operand {
	switch expr := expr.(type) {
	case ast.LiteralExpression:
		return Literal(expr.Value)
	default:
		log.Fatalf("unexpected expression type: %T", expr)
	}
	return nil
}

// this function will go away
func translateType(t ast.TypeDecl) types.Type {
	return types.Int{}
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
