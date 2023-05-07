package reader

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/introspection"
)

type Reader struct {
	introspectionAPI introspection.API
}

func NewReader(introspectionAPI introspection.API) *Reader {
	return &Reader{
		introspectionAPI: introspectionAPI,
	}
}

type ReaderError struct {
	Details []ReadError
}

func (e ReaderError) Error() string {
	details := ""
	for _, detail := range e.Details {
		details += detail.Error() + "\n"
	}
	return fmt.Sprintf("ReaderError: %s", details)
}

func (r Reader) Read(input *input.Input) (*DatumAST, error) {
	parser := NewParser(r.introspectionAPI)
	ast, errors := parser.Parse(input)

	if len(errors) > 0 {
		return nil, ReaderError{Details: errors}
	} else {
		return ast, nil
	}
}
