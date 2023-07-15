package types

import (
	"fmt"
	"log"

	"github.com/certainty/go-braces/internal/compiler/frontend/lexer"
	"github.com/certainty/go-braces/internal/compiler/frontend/parser/ast"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
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

func (t Checker) Check(ast *ast.AST) (TypeUniverse, error) {
	t.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseTypeCheck)
	defer t.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseTypeCheck)

	log.Printf("Type checking AST: %v", ast.ASTring())
	log.Printf("Type universe: %v", t.typeUniverse)

	for _, node := range ast.Nodes {
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
	case ast.LiteralExpression:
		return t.typeCheckLiteral(&node)
	case ast.BinaryExpression:
		return t.typeCheckBinaryExpression(&node)
	case ast.CallableDecl:
		return t.typeCheckCallableDecl(&node)
	default:
		return UnknownType, nil
	}
}

func (t Checker) typeCheckCallableDecl(node *ast.CallableDecl) (Type, error) {
	declaredType, err := t.typeForDecl(node.TpeDecl)
	if err != nil {
		return nil, err
	}
	var lastExprType Type
	for _, bodyNode := range node.Body.Code {
		lastExprType, err = t.typeCheck(bodyNode)
		if err != nil {
			return nil, err
		}
	}

	if lastExprType != declaredType {
		return nil, fmt.Errorf("declared type %v does not match body type %v", declaredType, lastExprType)
	}

	return t.assignType(node, declaredType), nil
}

func (t Checker) typeForDecl(tpeDecl ast.TypeDecl) (Type, error) {
	switch tpeDecl.Name.Label {
	case "int":
		return IntType, nil
	case "uint":
		return UIntType, nil
	case "float":
		return FloatType, nil
	default:
		return nil, fmt.Errorf("unknown type %v", tpeDecl.Name.Label)
		// handle other builtin and user-defined types
	}

}

func (t *Checker) typeCheckLiteral(node *ast.LiteralExpression) (Type, error) {
	if (*node).Value == nil {
		return nil, fmt.Errorf("literal has no value")
	} else {
		switch (*node).Value.(type) {
		case int:
			return t.assignType(node, IntType), nil
		case float64:
			return t.assignType(node, FloatType), nil
		case bool:
			return t.assignType(node, BoolType), nil
		case string:
			return t.assignType(node, StringType), nil
		case lexer.CodePoint:
			return t.assignType(node, CharType), nil
		default:
			return nil, fmt.Errorf("unknown literal type %T", (*node).Value)
		}
	}
}

func (t *Checker) typeCheckBinaryExpression(node *ast.BinaryExpression) (Type, error) {
	leftType, err := t.typeCheck(node.Left)
	if err != nil {
		return nil, err
	}
	rightType, err := t.typeCheck(node.Right)
	if err != nil {
		return nil, err
	}

	if node.IsNumeric() {
		if !(leftType == IntType || leftType == FloatType || leftType == UIntType) {
			return nil, fmt.Errorf("operand must be numeric")
		}
	}

	if node.IsBoolean() && leftType != BoolType {
		return nil, fmt.Errorf("operand must be boolean")
	}

	if leftType == rightType {
		return t.assignType(node, leftType), nil
	} else {
		return nil, fmt.Errorf("type mismatch: %v != %v", leftType, rightType)
	}
}
