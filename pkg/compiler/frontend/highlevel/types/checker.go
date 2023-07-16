package types

import (
	"fmt"
	"log"

	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/ast"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
)

type TypeUniverse map[ast.NodeId]Type

type Checker struct {
	instrumentation compiler_introspection.Instrumentation
	typeUniverse    map[ast.NodeId]Type
}

func NewChecker(Instrumentation compiler_introspection.Instrumentation) Checker {
	return Checker{instrumentation: Instrumentation, typeUniverse: make(map[ast.NodeId]Type)}
}

func (t *Checker) assignType(node ast.Node, tpe Type) Type {
	t.typeUniverse[node.ID()] = tpe
	return tpe
}

func (t Checker) Check(ast *ast.Source) (TypeUniverse, error) {
	t.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseTypeCheck)
	defer t.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseTypeCheck)

	log.Printf("Type checking AST: %v", ast.ASTString())
	log.Printf("Type universe: %v", t.typeUniverse)

	for _, node := range ast.Declarations {
		_, err := t.typeCheck(node)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("Type universe after check: %v", t.typeUniverse)
	return t.typeUniverse, nil
}

func (t Checker) typeCheck(node ast.Node) (Type, error) {
	switch node := node.(type) {
	case ast.BasicLitExpr:
		return t.typeCheckLiteral(&node)
	case ast.BinaryExpr:
		return t.typeCheckBinaryExpression(&node)
	default:
		return UnknownType, nil
	}
}

func (t *Checker) typeCheckLiteral(node *ast.BasicLitExpr) (Type, error) {
	switch node.Token.Type {
	case token.FIXNUM:
		return t.assignType(node, IntType), nil
	case token.FLONUM:
		return t.assignType(node, FloatType), nil
	case token.STRING:
		return t.assignType(node, StringType), nil
	case token.CHAR:
		return t.assignType(node, CharType), nil
	case token.BYTE:
		return t.assignType(node, ByteType), nil
	case token.BOOLEAN:
		return t.assignType(node, BoolType), nil
	default:
		return nil, fmt.Errorf("unknown literal type %T", (*node).Value)
	}
}

func (t *Checker) typeCheckBinaryExpression(node *ast.BinaryExpr) (Type, error) {
	leftType, err := t.typeCheck(node.Left)
	if err != nil {
		return nil, err
	}
	rightType, err := t.typeCheck(node.Right)
	if err != nil {
		return nil, err
	}

	// make sure left operand matches required type
	switch node.Op.Type {
	case token.ADD, token.SUB, token.MUL, token.DIV, token.REM, token.POW:
		if leftType != IntType && leftType != FloatType && leftType != UIntType {
			return nil, fmt.Errorf("operands must be numeric")
		}
	case token.LAND, token.LOR:
		if leftType != BoolType {
			return nil, fmt.Errorf("operands must be boolean")
		}
	}

	// now make sure right operand matches left operand
	if leftType == rightType {
		return t.assignType(node, leftType), nil
	} else {
		return nil, fmt.Errorf("type mismatch: %v != %v", leftType, rightType)
	}
}
