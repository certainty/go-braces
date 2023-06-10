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

func NewCoreAST() CoreAST {
	return CoreAST{Nodes: make([]CoreNode, 0)}
}

// desugar the AST into a core AST
func LowerToCore(theAST *ast.AST) (*CoreAST, error) {
	coreAST := NewCoreAST()

	for _, expression := range theAST.Nodes {
		coreNode, err := lowerNode(expression)
		if err != nil {
			return nil, err
		}
		coreAST.Nodes = append(coreAST.Nodes, coreNode)
	}

	log.Printf("core %v", coreAST) // TODO: build a writer for the coreAST
	return &coreAST, nil
}

func lowerNode(node ast.Node) (CoreNode, error) {
	switch node := node.(type) {
	case ast.LiteralExpression:
		return NewConstant(node.Value, node.Location()), nil
	case ast.UnaryExpression:
		// convert to call
		return nil, fmt.Errorf("unhandled expression type %T", node)
	case ast.BinaryExpression:
		// convert to call
		return nil, fmt.Errorf("unhandled expression type %T", node)
	default:
		return nil, fmt.Errorf("unhandled expression type %T", node)
	}
}
