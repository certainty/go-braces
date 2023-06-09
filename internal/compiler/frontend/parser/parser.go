package parser

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/frontend/lexer"
	"github.com/certainty/go-braces/internal/compiler/frontend/parser/ast"
	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
)

type Precedence uint8

const (
	PREC_NONE Precedence = iota
	PREC_SET
	PREC_OR
	PREC_AND
	PREC_EQUALITY   // == !=
	PREC_COMPARISON // < > <= >=
	PREC_TERM       // + -
	PREC_FACTOR     // * / % /
	PREC_EXPONENT   // ^
	PREC_UNARY      // ! -
	PREC_CALL       //  ()
	PREC_PRIMARY
)

func precedenceFor(tokenType lexer.TokenType) Precedence {
	switch tokenType {
	case lexer.TOKEN_STAR, lexer.TOKEN_SLASH, lexer.TOKEN_MOD:
		return PREC_FACTOR
	case lexer.TOKEN_CARET:
		return PREC_EXPONENT
	case lexer.TOKEN_PLUS, lexer.TOKEN_MINUS:
		return PREC_TERM
	case lexer.TOKEN_EQUAL_EQUAL, lexer.TOKEN_BANG_EQUAL, lexer.TOKEN_LT, lexer.TOKEN_LT_EQUAL, lexer.TOKEN_GT, lexer.TOKEN_GT_EQUAL:
		return PREC_COMPARISON
	case lexer.TOKEN_AMPERSAND_AMPERSAND:
		return PREC_AND
	case lexer.TOKEN_PIPE_PIPE:
		return PREC_OR
	case lexer.TOKEN_LPAREN:
		return PREC_NONE // is this correct?
	default:
		return PREC_NONE
	}
}

type Parser struct {
	instrumentation compiler_introspection.Instrumentation
	scanner         *lexer.Scanner
	previousToken   *lexer.Token
	currentToken    *lexer.Token
	ast             *ast.AST
	errors          []ParseError
	panicMode       bool
	hadError        bool
}

func NewParser(instrumentation compiler_introspection.Instrumentation) *Parser {
	p := &Parser{
		instrumentation: instrumentation,
		scanner:         nil,
	}
	return p
}

// the parser follows "crafting interpreters" and uses a recursive descent parser using pratt's top down operator precedence
func (p *Parser) Parse(input *input.Input) (*ast.AST, error) {
	// reset the parser state
	p.errors = []ParseError{}
	p.panicMode = false
	p.hadError = false
	p.previousToken = nil
	p.currentToken = nil
	p.ast = ast.New()
	p.scanner = lexer.New(input)

	// now start the parsing process
	p.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseParse)
	defer p.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseParse)

	return p.parseInput()
}

func (p *Parser) parseInput() (*ast.AST, error) {
	p.advance()
	expr := p.parseExpression()
	p.ast.AddExpression(expr)

	p.advance()
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

func (p *Parser) parseExpression() ast.Expression {
	return p.parseBinaryExpressions(PREC_SET)
}

// precedence climbing algorithm to make sure we treat precedence and associativity correctly
func (p *Parser) parseBinaryExpressions(minPrecedence Precedence) ast.Expression {
	left := p.parseUnaryExpression()
	for {
		precedence := precedenceFor(p.currentToken.Type)
		if precedence < minPrecedence {
			break
		}
		tok := p.currentToken
		p.advance()
		// now we climb the precedence ladder
		right := p.parseBinaryExpressions(precedence + 1)
		left = ast.BinOp(tok.Location, ast.TokenToBinaryOp(*tok), left, right)
	}
	return left
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	if p.match(lexer.TOKEN_BANG, lexer.TOKEN_MINUS, lexer.TOKEN_PLUS) {
		tok := p.previousToken
		right := p.parseUnaryExpression()
		return ast.UnaryOp(tok.Location, ast.TokenToUnaryOp(*tok), right)
	} else if p.match(lexer.TOKEN_LPAREN) {
		expr := p.parseBinaryExpressions(PREC_SET)
		p.consume(lexer.TOKEN_RPAREN, "expected closing parenthesis")
		return expr
	} else if p.match(lexer.TOKEN_IDENTIFIER) {
		fmt.Printf("Parsing identifier - not yes supported\n")
		return nil
	} else if p.match(lexer.TOKEN_INTEGER) {
		return ast.NewLiteralExpression(*p.previousToken, p.previousToken.Location)
	} else if p.match(lexer.TOKEN_FLOAT) {
		return ast.NewLiteralExpression(*p.previousToken, p.previousToken.Location)
	} else if p.match(lexer.TOKEN_EOF) {
		p.errorAtCurrent(ParseErrorIdUnexpectedEOF, "unexpected end of input", nil)
		return nil
	}
	p.errorAtCurrent(ParseErrorIdUnexpectedToken, "expected unary expression", nil)
	return nil
}

func (p *Parser) match(tokenTypes ...lexer.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType lexer.TokenType) bool {
	return p.currentToken.Type == tokenType
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
