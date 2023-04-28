package transformer

import "github.com/certainty/go-braces/internal/introspection"

type CPSTransformer struct {
	introspectionAPI *introspection.API
}

func NewCpsTransformer(introspectionAPI *introspection.API) *CPSTransformer {
	return &CPSTransformer{introspectionAPI: introspectionAPI}
}
