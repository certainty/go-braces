package parser

import (
	"fmt"

	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
)

type ParseErrorId int

const (
	ParseErrorGeneric ParseErrorId = iota
	ParseErrorIdUnexpectedEOF
	ParseErrorIdUnexpectedToken
	ParseErrorLexerError
	ParseErrorNotImplemented
)

type ParseError struct {
	Where   token.Location
	What    string
	Cause   error
	Context string
	Id      ParseErrorId
}

type ParseErrors struct {
	Errors []ParseError
}

func (e ParseErrors) Error() string {
	return fmt.Sprintf("parse errors: %v", e.Errors)
}
