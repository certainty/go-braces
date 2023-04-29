package compiler

import (
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
)

type MultiError struct {
	errors []error
}

func (m *MultiError) Error() string {
	errorMessages := make([]string, len(m.errors))
	for i, e := range m.errors {
		errorMessages[i] = e.Error()
	}
	return strings.Join(errorMessages, "\n")
}

func (m *MultiError) Add(err error) {
	if err != nil {
		m.errors = append(m.errors, err)
	}
}

func (m *MultiError) Empty() bool {
	return len(m.errors) == 0
}

type ErrorHandler struct {
	multiError MultiError
}

func (h *ErrorHandler) HandleError(err error, pos *lexer.Position) (string, error) {
	h.multiError.Add(fmt.Errorf("%s: %w", pos, err))
	return "", nil
}

type ErrorWithID interface {
	ID() string
}

////
// Reader errors
///

const (
	ReaderErrorInvalidCharacter = "0x0001"
)

type ReaderError struct {
	id          string
	Message     string
	Pos         lexer.Position
	FailedSlice string
	Origin      error
}

func (e ReaderError) Error() string {
	return fmt.Sprintf("Reader error: %s", e.Message)
}

func (e ReaderError) ID() string {
	return e.id
}

////
// ParserError errors
///

const (
	ParserErrorInvalidExpression = "0x0050"
)

type ParserError struct {
	id      string
	Message string
	// TODO add location
	// TODO add information about the expression that failed and where it orginated from (maybe macro expansion)
}

func (e ParserError) Error() string {
	return fmt.Sprintf("Parser error: %s", e.Message)
}

func (e ParserError) ID() string {
	return e.id
}

type ErrorDetail struct {
	Type            string
	underlyingError error
}

// Combination of all errors
type CompilerError struct {
	details []ErrorDetail
}

func (e CompilerError) Error() string {
	return "Compiler error"
}
