package transformer

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/parser"
	"github.com/certainty/go-braces/internal/introspection"
)

type CPSTransformer struct {
	introspectionAPI introspection.API
}

func NewCpsTransformer(introspectionAPI introspection.API) *CPSTransformer {
	return &CPSTransformer{introspectionAPI: introspectionAPI}
}

func (c *CPSTransformer) Transform(coreAst *parser.CoreAST) (*parser.CoreAST, error) {
	return nil, nil
}
