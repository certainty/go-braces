package ir

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/parser"
	"github.com/certainty/go-braces/internal/introspection"
)

type Block struct{}

type SSA struct {
	blocks []Block
}

type SSATransformer struct {
	introspectionAPI introspection.API
}

func NewSSATransformer(introspectionAPI introspection.API) *SSATransformer {
	return &SSATransformer{introspectionAPI: introspectionAPI}
}

func (c *SSATransformer) Transform(coreAst *parser.CoreAST) (*SSA, error) {
	return nil, nil
}
