package parser

import "fmt"

type CoreAST struct {
	Expressions []SchemeExpression
}

func NewCoreAST() *CoreAST {
	return &CoreAST{
		Expressions: []SchemeExpression{},
	}
}

func (ast *CoreAST) String() string {
	return fmt.Sprintf("CoreAST %s ", ast.Expressions)
}

func (ast *CoreAST) AddExpression(expression SchemeExpression) {
	ast.Expressions = append(ast.Expressions, expression)
}
