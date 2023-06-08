package parser

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/frontend/lexer"
	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/isa"
)

type Parser struct {
	instrumentation compiler_introspection.Instrumentation
	scanner         *lexer.Scanner
	previousToken   *lexer.Token
	currentToken    *lexer.Token
	ast             *AST
	errors          []ParseError
	panicMode       bool
	hadError        bool
}

func NewParser(instrumentation compiler_introspection.Instrumentation) *Parser {
	return &Parser{
		instrumentation: instrumentation,
		scanner:         nil,
	}
}

// the parser follows "crafting interpreters" and uses a recursive descent parser using pratt's top down operator precedence
func (p *Parser) Parse(input *input.Input) (*AST, error) {
	// reset the parser state
	p.errors = []ParseError{}
	p.panicMode = false
	p.hadError = false
	p.previousToken = nil
	p.currentToken = nil
	p.ast = NewAST()
	p.scanner = lexer.New(input.Buffer, input.Origin)

	p.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseParse)
	defer p.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseParse)

	return p.parseInput()
}

func (p *Parser) parseInput() (*AST, error) {
	p.advance()

	expr := p.parseExpression()
	p.ast.AddExpression(expr)

	p.consume(lexer.TOKEN_EOF, "expected end of input")

	if p.hadError {
		return nil, ParseErrors{Errors: p.errors}
	} else {
		return p.ast, nil
	}
}

func (p *Parser) advance() {
	p.previousToken = p.currentToken
	for {
		nextToken, err := p.scanner.NextToken()
		p.currentToken = &nextToken
		if err == nil {
			break
		}
		p.errorAtCurrent(ParseErrorLexerError, err.Error(), err)
	}
}

func (p *Parser) consume(tokenType lexer.TokenType, message string) {
	if p.currentToken.Type == tokenType {
		p.advance()
		return
	}
	p.errorAtCurrent(ParseErrorIdUnexpectedToken, message, nil)
}

func (p *Parser) parseExpression() Expression {
	return p.parseLiteral()
}

func (p *Parser) parseLiteral() Expression {
	var value isa.Value
	switch p.currentToken.Type {
	case lexer.TOKEN_TRUE:
		value = isa.BoolValue(true)
	case lexer.TOKEN_FALSE:
		value = isa.BoolValue(false)
	case lexer.TOKEN_CHARACTER:
		value = isa.CharValue(p.currentToken.Value.(rune))
	case lexer.TOKEN_INTEGER:
		value = isa.IntegerValue(p.currentToken.Value.(int64))
	case lexer.TOKEN_FLOAT:
		value = isa.FloatValue(p.currentToken.Value.(float64))
	case lexer.TOKEN_STRING:
		value = isa.StringValue(p.currentToken.Value.(string))
	default:
		p.errorAtCurrent(ParseErrorIdUnexpectedToken, "expected literal", nil)
		return nil
	}
	p.advance()
	return NewLiteralExpression(value, p.previousToken.Location)
}

func (p *Parser) compilerBug(message string) {
	panic(fmt.Sprintf("compiler bug: %s", message))
}

func (p *Parser) errorAtCurrent(id ParseErrorId, message string, cause error) {
	p.errorAt(*p.currentToken, id, message, cause)
}

func (p *Parser) errorAtPrevious(id ParseErrorId, message string, cause error) {
	p.errorAt(*p.previousToken, id, message, cause)
}

func (p *Parser) errorAt(token lexer.Token, id ParseErrorId, message string, cause error) {
	if p.panicMode {
		return
	}

	p.panicMode = true
	p.hadError = true

	// todo add more details to the error message including the line and column
	// and the code that is being parsed
	p.errors = append(p.errors, ParseError{
		Id:      id,
		Where:   token.Location,
		What:    message,
		Context: fmt.Sprintf("Token: %s at line %d (%d..%d)", token, token.Location.Line, token.Location.StartOffset, token.Location.EndOffset),
		Cause:   cause,
	})
}
