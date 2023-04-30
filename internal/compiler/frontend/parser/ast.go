package parser

type CoreAST struct {
	Expressions []SchemeExpression
}

func NewCoreAST() *CoreAST {
	return &CoreAST{
		Expressions: []SchemeExpression{},
	}
}

func (ast *CoreAST) AddExpression(expression SchemeExpression) {
	ast.Expressions = append(ast.Expressions, expression)
}
