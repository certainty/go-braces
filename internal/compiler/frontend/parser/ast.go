package parser

import "fmt"

type AST struct {
	Expressions []Expression
}

func New() *AST {
	return &AST{
		Expressions: []Expression{},
	}
}

func (ast *AST) String() string {
	return fmt.Sprintf("CoreAST %s ", ast.Expressions)
}

func (ast *AST) AddExpression(expression Expression) {
	ast.Expressions = append(ast.Expressions, expression)
}
