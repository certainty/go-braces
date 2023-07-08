package ir

import (
	"fmt"
	"github.com/certainty/go-braces/internal/compiler/frontend/parser/ast"
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

func LowerToIR(origin location.Origin, ast *ast.AST, moduleName Label) (*Module, error) {
	mod := CreateModule(moduleName, origin)
	for _, node := range ast.Nodes {
		switch node.(type) {

		}
	}

	return nil, nil
}
