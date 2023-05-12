package reader

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
)

type Reader struct {
	instrumentation compiler_introspection.Instrumentation
}

func NewReader(instrumentation compiler_introspection.Instrumentation) *Reader {
	return &Reader{
		instrumentation: instrumentation,
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
	r.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseRead)
	defer r.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseRead)

	parser := NewParser(r.instrumentation)
	ast, errors := parser.Parse(input)

	if len(errors) > 0 {
		return nil, ReaderError{Details: errors}
	} else {
		return ast, nil
	}
}
