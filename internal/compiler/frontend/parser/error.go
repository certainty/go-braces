package parser

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/location"
)

type ParseErrorId int

const (
	ParseErrorGeneric ParseErrorId = iota
	ParseErrorIdUnexpectedEOF
	ParseErrorIdUnexpectedToken
	ParseErrorLexerError
)

type ParseError struct {
	Where   location.Location
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
