package typechecker

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/parser"
	"github.com/certainty/go-braces/internal/introspection"
)

type TypeChecker struct {
	introspectionAPI introspection.API
}

func NewTypeChecker(introspectionAPI introspection.API) *TypeChecker {
	return &TypeChecker{introspectionAPI: introspectionAPI}
}

func (t *TypeChecker) Check(ast *parser.CoreAST) error {
	return nil
}
