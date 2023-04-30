package reader

import (
	"github.com/certainty/go-braces/internal/compiler/location"
	"github.com/certainty/go-braces/internal/introspection"
)

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

type ReaderError struct {
	Details []ReadError
}

func (e ReaderError) Error() string {
	return "ReaderError"
}

func (r *Reader) Read(input location.Input) (*DatumAST, error) {
	ast, errors := r.parser.Parse(input)
	if errors != nil && len(errors) > 0 {
		return nil, ReaderError{Details: errors}
	} else {
		return ast, nil
	}
}
