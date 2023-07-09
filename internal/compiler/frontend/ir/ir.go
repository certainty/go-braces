package ir

import (
	"fmt"
	"github.com/certainty/go-braces/internal/compiler/frontend/parser/ast"
	"github.com/certainty/go-braces/internal/compiler/frontend/types"
	"github.com/certainty/go-braces/internal/compiler/location"
)

type Label string
type Register string

type Module struct {
	Name      Label
	Source    location.Origin
	Functions []*Function
}

type Function struct {
	tpe    Type
	Name   Label
	Blocks []*BasicBlock
}

type BasicBlock struct {
	Label        Label
	Instructions []Instruction
}

type Instruction interface{}

// %register = op tpe operand1, operand2, ...
type SimpleInstruction struct {
	tpe       Type
	Register  Register
	Operation Operation
	Operand   Operand
}

type ReturnInstruction struct {
	tpe      Type
	Register Register
}

type Operation uint8
type Operand interface{}
type Constant interface{}

var _ Operand = Constant(nil)
var _ Operand = Label("")
var _ Operand = Register("")

const (
	Ret Operation = iota
	Add
	Sub
	Mul
	Div
	Or
	And
	Xor
	Neg
)

func LowerToIR(origin location.Origin, theAst *ast.AST, moduleName Label) (*Module, error) {
	mod := CreateModule(moduleName, origin)

	for _, node := range theAst.Nodes {
		switch node := node.(type) {
		case ast.FunctionDecl:
			fun, err := LowerFunction(node)
			if err != nil {
				return nil, err
			}
			mod.Functions = append(mod.Functions, fun)
		default:
			return nil, fmt.Errorf("unexpected node type: %T", node)
		}
	}

	return mod, nil
}

func LowerFunction(decl ast.FunctionDecl) (*Function, error) {
	tpe, err := lowerType(decl.Type)
	if err != nil {
		return nil, err
	}
	fun := CreateFunction(tpe, Label(decl.Name.ID))
	return fun, nil
}

func lowerType(t types.Type) (Type, error) {
	switch t.(type) {
	case types.Int:
		// for now
		return Int64Type, nil
	case types.Float:
		return Float32Type, nil
	case types.Bool:
		return UInt8Type, nil
	case types.String:
		return StringType, nil
	default:
		return nil, fmt.Errorf("unexpected type: %s", t)
	}
}
