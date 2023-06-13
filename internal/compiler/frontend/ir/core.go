package ir

import (
	"fmt"
	"log"

	"github.com/certainty/go-braces/internal/compiler/frontend/parser/ast"
	"github.com/certainty/go-braces/internal/compiler/location"
)

type CoreAST struct {
	Nodes []CoreNode
}

type CoreNode interface {
	Location() location.Location
}

type CoreConstant struct {
	location location.Location
	Value    interface{}
}

var _ CoreNode = (*CoreConstant)(nil)

func NewConstant(value interface{}, location location.Location) CoreConstant {
	return CoreConstant{Value: value, location: location}
}

func (c CoreConstant) Location() location.Location {
	return c.location
}

type Call struct {
	location location.Location
	Operator Callable
	Operands []CoreNode
}

var _ CoreNode = (*Call)(nil)

func NewCall(location location.Location, operator Callable, operands ...CoreNode) Call {
	return Call{Operator: operator, Operands: operands, location: location}
}

func (c Call) Location() location.Location {
	return c.location
}

type Callable interface {
	CoreNode
}

type PrimitiveOp uint8

const (
	PrimitiveOpAdd PrimitiveOp = iota
	PrimitiveOpSub
	PrimitiveOpMul
	PrimitiveOpDiv
	PrimitiveOpPow
	PrimitiveOpMod
	PrimitiveOpAnd
	PrimitiveOpOr
	PrimitiveOpNot
	PrimitiveOpNeg
)

type Primitive struct {
	location location.Location
	Op       PrimitiveOp
}

var _ Callable = (*Primitive)(nil)
var _ CoreNode = (*Primitive)(nil)

func NewPrimitive(op PrimitiveOp, location location.Location) Primitive {
	return Primitive{Op: op, location: location}
}

func (p Primitive) Location() location.Location {
	return p.location
}

type JunctionOp uint8

const (
	JunctionOpAnd JunctionOp = iota
	JunctionOpOr
)

type Junction struct {
	location location.Location
	Junctor  JunctionOp
	Left     CoreNode
	Right    CoreNode
}

func NewJunction(junctor JunctionOp, left CoreNode, right CoreNode, location location.Location) Junction {
	return Junction{Junctor: junctor, Left: left, Right: right, location: location}
}

func (j Junction) Location() location.Location {
	return j.location
}

func NewCoreAST() CoreAST {
	return CoreAST{Nodes: make([]CoreNode, 0)}
}

// desugar the AST into a core AST
func LowerToCore(theAST *ast.AST) (*CoreAST, error) {
	coreAST := NewCoreAST()
	coreASTWriter := NewCoreASTWriter()

	for _, expression := range theAST.Nodes {
		coreNode, err := lowerNode(expression)
		if err != nil {
			return nil, err
		}
		coreAST.Nodes = append(coreAST.Nodes, coreNode)
	}

	log.Printf("CORE: %s", coreASTWriter.Write(coreAST)) // TODO: build a writer for the coreAST
	return &coreAST, nil
}

func lowerNode(node ast.Node) (CoreNode, error) {
	switch node := node.(type) {
	case ast.LiteralExpression:
		return NewConstant(node.Value, node.Location()), nil
	case ast.UnaryExpression:
		return nil, fmt.Errorf("unhandled expression type %T", node)
	case ast.BinaryExpression:
		left, err := lowerNode(node.Left)
		if err != nil {
			return nil, err
		}
		right, err := lowerNode(node.Right)
		if err != nil {
			return nil, err
		}
		// TODO: distinguish between booleans and other binary expression in source AST
		switch node.Operator {
		case ast.BinOpAnd:
			return NewJunction(JunctionOpAnd, left, right, node.Location()), nil
		case ast.BinOpOr:
			return NewJunction(JunctionOpOr, left, right, node.Location()), nil
		default:
			return NewCall(node.Location(), NewPrimitive(callableFromOperator(node.Operator), node.Location()), left, right), nil
		}
	default:
		return nil, fmt.Errorf("unhandled expression type %T", node)
	}
}

func callableFromOperator(operator ast.BinaryOperator) PrimitiveOp {
	switch operator {
	case ast.BinOpAdd:
		return PrimitiveOpAdd
	case ast.BinOpSub:
		return PrimitiveOpSub
	case ast.BinOpMul:
		return PrimitiveOpMul
	case ast.BinOpDiv:
		return PrimitiveOpDiv
	case ast.BinOpPow:
		return PrimitiveOpPow
	case ast.BinOpMod:
		return PrimitiveOpMod
	case ast.BinOpAnd:
		return PrimitiveOpAnd
	case ast.BinOpOr:
		return PrimitiveOpOr
	default:
		panic(fmt.Sprintf("unhandled operator %v", operator))
	}
}
