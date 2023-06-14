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

type Apply struct {
	location location.Location
	Operator Appicable
	Operands []CoreNode
}

var _ CoreNode = (*Apply)(nil)

func NewApply(location location.Location, operator Appicable, operands ...CoreNode) Apply {
	return Apply{Operator: operator, Operands: operands, location: location}
}

func (c Apply) Location() location.Location {
	return c.location
}

type Appicable interface {
	CoreNode
}

type PrimitiveOp uint8

type Primitive struct {
	location location.Location
	Op       PrimitiveOp
}

var _ Appicable = (*Primitive)(nil)
var _ CoreNode = (*Primitive)(nil)

func NewPrimitive(op PrimitiveOp, location location.Location) Primitive {
	return Primitive{Op: op, location: location}
}

func (p Primitive) Location() location.Location {
	return p.location
}

type LogicalOperator uint8

type LogicalConnective struct {
	location location.Location
	Operator LogicalOperator
	Left     CoreNode
	Right    CoreNode
}

func NewLogicalConnective(op LogicalOperator, left CoreNode, right CoreNode, location location.Location) LogicalConnective {
	return LogicalConnective{Operator: op, Left: left, Right: right, location: location}
}

func (j LogicalConnective) Location() location.Location {
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
			return NewLogicalConnective(LogicalOperator(ast.BinOpAdd), left, right, node.Location()), nil
		case ast.BinOpOr:
			return NewLogicalConnective(LogicalOperator(ast.BinOpOr), left, right, node.Location()), nil
		default:
			return NewApply(node.Location(), NewPrimitive(PrimitiveOp(node.Operator), node.Location()), left, right), nil
		}
	default:
		return nil, fmt.Errorf("unhandled expression type %T", node)
	}
}
