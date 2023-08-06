package types

import (
	"fmt"
	"log"

	"github.com/certainty/go-braces/pkg/compiler/frontend/astutils"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/ast"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
)

type TypeUniverse struct {
	ProcTypes       map[string]Procedure
	ExpressionTypes map[astutils.NodeId]Type
}

func NewTypeUniverse() *TypeUniverse {
	return &TypeUniverse{
		ProcTypes:       make(map[string]Procedure),
		ExpressionTypes: make(map[astutils.NodeId]Type),
	}
}

type Checker struct {
	instrumentation compiler_introspection.Instrumentation
	typeUniverse    *TypeUniverse
}

func NewChecker(Instrumentation compiler_introspection.Instrumentation) Checker {
	return Checker{instrumentation: Instrumentation, typeUniverse: NewTypeUniverse()}
}

func (t *Checker) assignType(node ast.Node, tpe Type) Type {
	t.typeUniverse.ExpressionTypes[node.ID()] = tpe
	return tpe
}

func (t Checker) Check(ast *ast.Source) (*TypeUniverse, error) {
	t.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseTypeCheck)
	defer t.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseTypeCheck)

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
	case ast.ProcDecl:
		return t.typeCheckProcDecl(&node)
	case ast.BasicLitExpr:
		return t.typeCheckLiteral(&node)
	case ast.BinaryExpr:
		return t.typeCheckBinaryExpression(&node)
	case ast.BlockExpr:
		return t.typeCheckBlockExpr(&node)
	case ast.ExprStmt:
		return t.typeCheck(node.Expr)
	default:
		log.Printf("Unknown node type %T", node)
		return UnknownType, nil
	}
}

func (t *Checker) typeCheckProcDecl(node *ast.ProcDecl) (Type, error) {
	paramTypes := make([]Type, len(node.Type.Params))
	resultTypes := []Type{}

	for _, param := range node.Type.Params {
		paramType, err := t.typeFromName(param.Type.Name)
		if err != nil {
			return nil, err
		}
		t.typeUniverse.ExpressionTypes[param.ID()] = paramType
		paramTypes = append(paramTypes, paramType)
	}

	if node.Type.Result != nil {
		resultType, err := t.typeFromName(node.Type.Result.Name)
		if err != nil {
			return nil, err
		}
		t.typeUniverse.ExpressionTypes[node.Type.Result.ID()] = resultType
		resultTypes = append(resultTypes, resultType)
	} else {
		resultTypes = append(resultTypes, UnitType)
	}

	procType := Procedure{
		Params:  paramTypes,
		Results: resultTypes,
	}

	t.typeUniverse.ProcTypes[node.Name.Name] = procType
	t.typeUniverse.ExpressionTypes[node.ID()] = procType

	blockType, err := t.typeCheck(node.Body)
	if err != nil {
		return nil, err
	}
	if blockType != procType.Results[0] {
		return nil, fmt.Errorf("type mismatch: %v != %v", blockType, procType.Results[0])
	}

	return procType, nil
}

func (t *Checker) typeCheckBlockExpr(node *ast.BlockExpr) (Type, error) {
	var blockReturnType Type = UnitType
	var err error

	for _, stmt := range node.Statements {
		blockReturnType, err = t.typeCheck(stmt)
		if err != nil {
			return nil, err
		}
	}
	return blockReturnType, nil
}

func (t *Checker) typeFromName(typeIdentifier ast.Identifier) (Type, error) {
	switch typeIdentifier.Name {
	case "Int":
		return IntType, nil
	case "Float":
		return FloatType, nil
	case "UInt":
		return UIntType, nil
	case "Bool":
		return BoolType, nil
	case "String":
		return StringType, nil
	case "Char":
		return CharType, nil
	case "Byte":
		return ByteType, nil
	case "Unit":
		return UnitType, nil
	default:
		return nil, fmt.Errorf("unknown type %v", typeIdentifier.Name)
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
