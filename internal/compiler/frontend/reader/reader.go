package reader

import "github.com/certainty/go-braces/internal/introspection"

type Reader struct {
	introspectionAPI introspection.API
	parser           *Parser
}

func NewReader(introspectionAPI introspection.API) *Reader {
	return &Reader{
		introspectionAPI: introspectionAPI,
		parser:           NewParser(introspectionAPI),
	}
}

func (r *Reader) Read(source *[]byte) (*DatumAST, error) {
	return nil, nil
}
