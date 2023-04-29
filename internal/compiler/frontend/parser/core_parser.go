package parser

import "github.com/certainty/go-braces/internal/introspection"

type CoreParser struct {
	introspectionAPI introspection.API
}

func NewCoreParser(introspectionAPI introspection.API) *CoreParser {
	return &CoreParser{introspectionAPI: introspectionAPI}
}
