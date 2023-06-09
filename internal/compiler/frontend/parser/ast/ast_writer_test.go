package ast

import (
	"testing"
)

func TestASTWriter_Write(t *testing.T) {
	writer := NewASTWriter()

	ast := &AST{
		Nodes: []Node{
			&BinaryExpression{
				Left:     &LiteralExpression{Value: 2},
				Operator: BinOpAdd,
				Right:    &LiteralExpression{Value: 3},
			},
			&UnaryExpression{
				Operator: UnaryOpNeg,
				Operand:  &LiteralExpression{Value: 4},
			},
			&BinaryExpression{
				Left:     &LiteralExpression{Value: true},
				Operator: BinOpAnd,
				Right: &UnaryExpression{
					Operator: UnaryOpNot,
					Operand:  &LiteralExpression{Value: false},
				},
			},
		},
	}

	expected := "((+ 2 3) (- 4) (and true (not false)))"
	actual := writer.Write(ast)

	if actual != expected {
		t.Errorf("Expected output to be %q, but got %q", expected, actual)
	}
}
