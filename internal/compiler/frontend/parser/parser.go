package parser

import (
	"fmt"

	ast "github.com/certainty/go-braces/internal/compiler/frontend/ast/hl"
	"github.com/certainty/go-braces/internal/compiler/frontend/lexer"
	"github.com/certainty/go-braces/internal/compiler/frontend/token"
	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
)

type Precedence uint8
type Associativity uint8

const (
	PREC_NONE     Precedence = iota
	PREC_SET                 // set x expr
	PREC_OR                  // ||
	PREC_AND                 // &&
	PREC_EQ                  // == !=
	PREC_COMP                // < > <= >=
	PREC_TERM                // + -
	PREC_FACTOR              // * / % /
	PREC_EXPONENT            // **
	PREC_UNARY               // ! -
	PREC_CALL                //  ()
	PREC_PRIMARY
)

const (
	ASSOC_NONE Associativity = iota
	ASSOC_LEFT
	ASSOC_RIGHT
)

func precedenceFor(tokenType token.Type) (Precedence, Associativity) {
	switch tokenType {
	case token.MUL, token.DIV, token.REM:
		return PREC_FACTOR, ASSOC_LEFT
	case token.POW:
		return PREC_EXPONENT, ASSOC_RIGHT
	case token.ADD, token.SUB:
		return PREC_TERM, ASSOC_LEFT
	case token.EQ, token.NEQ, token.LT, token.LTE, token.GT, token.GTE:
		return PREC_COMP, ASSOC_LEFT
	case token.LAND:
		return PREC_AND, ASSOC_LEFT
	case token.LOR:
		return PREC_OR, ASSOC_LEFT
	case token.LPAREN:
		return PREC_NONE, ASSOC_LEFT
	default:
		return PREC_NONE, ASSOC_NONE
	}
}

type Parser struct {
	instrumentation compiler_introspection.Instrumentation
	scanner         *lexer.Scanner
	previousToken   *token.Token
	currentToken    *token.Token
	astBuilder      *ast.Builder
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

func (p *Parser) Reset() {
	p.errors = []ParseError{}
	p.panicMode = false
	p.hadError = false
	p.previousToken = nil
	p.currentToken = nil
}

func (p *Parser) Parse(input *input.Input) (*ast.Source, error) {
	p.Reset()
	p.astBuilder = ast.NewBuilder()
	p.scanner = lexer.New(input)

	p.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseParse)
	defer p.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseParse)

	return p.parseInput()
}

func (p *Parser) parseInput() (*ast.Source, error) {
	p.advance()
	for {
		p.parseDeclaration()
		if p.currentToken.IsEOF() {
			break
		}
	}

	source := ast.Source{
		Statements: p.astBuilder.Result(),
	}

	if p.hadError {
		return &source, ParseErrors{Errors: p.errors}
	} else {
		return &source, nil
	}
}

func (p *Parser) advance() {
	p.previousToken = p.currentToken
	for {
		nextToken := p.scanner.NextToken()
		p.currentToken = &nextToken
		if !nextToken.IsIllegal() {
			break
		}

		p.errorAtCurrent(ParseErrorLexerError, "invalid token")
	}
}

func (p *Parser) consume(tokenType token.Type, message string) {
	if p.currentToken.Type == tokenType {
		p.advance()
		return
	}
	p.errorAtCurrent(ParseErrorIdUnexpectedToken, message)
}

// Try to recover to next synchrnization point.
// These are:
// * statement boundaries
// * new blocks
// * function boundaries
func (p *Parser) synchronize() {
	p.panicMode = false

	for !p.currentToken.IsEOF() {
		if p.previousToken.Type == token.SEMICOLON {
			return
		}

		switch p.currentToken.Type {
		case token.RBRACE, token.PROC, token.IF, token.FOR:
			return
		default:
			p.advance()
		}
	}
}

func (p *Parser) parseDeclaration() {
	switch p.currentToken.Type {
	case token.PROC:
		p.parseProcedureDeclaration()
	case token.EOF:
		p.advance()
		return
	default:
		p.errorAtCurrent(ParseErrorIdUnexpectedToken, "expected declaration")
	}

	if p.hadError {
		p.synchronize()
	}
}

func (p *Parser) parseProcedureDeclaration() {
	p.consume(token.PROC, "expected proc")
	location := p.currentToken.Location
	procName := p.parseIdentifier()
	params := p.parseArguments()

	var result *ast.TypeSpec = nil
	if p.match(token.COLON) {
		r := p.parseTypeSpec()
		result = &r
	}

	body := p.parseBlock()
	procedure := p.astBuilder.NewProcDecl(location, procName, params, result, body)
	p.astBuilder.AddNode(procedure)
}

func (p *Parser) parseArguments() []ast.Field {
	args := []ast.Field{}
	p.consume(token.LPAREN, "expected '('")
	for {
		if p.check(token.RPAREN) {
			break
		}
		argName := p.parseIdentifier()
		argType := p.parseTypeSpec()
		args = append(args, p.astBuilder.NewField(argName, &argType))
		if !p.match(token.COMMA) {
			break
		}
	}
	p.consume(token.RPAREN, "expected ')'")
	return args
}

func (p *Parser) parseTypeSpec() ast.TypeSpec {
	location := p.currentToken.Location
	p.consume(token.COLON, "expected ':'")
	return p.astBuilder.NewTypeSpec(location, p.parseIdentifier())
}

func (p *Parser) parseIdentifier() ast.Identifier {
	p.consume(token.IDENTIFIER, "expected identifier")
	return p.astBuilder.NewIdentifier(
		p.previousToken.Location,
		string(p.previousToken.Text),
	)
}

func (p *Parser) parseBlock() ast.BlockExpr {
	// TODO: track scope
	p.consume(token.LBRACE, "expected '{'")
	location := p.previousToken.Location
	statements := []ast.Statement{}

	for !p.match(token.RBRACE) {
		statements = append(statements, p.parseBlockStatment())
	}
	return p.astBuilder.NewBlockExpr(location, statements)
}

func (p *Parser) parseBlockStatment() ast.Statement {
	if p.match(token.LBRACE) {
		return p.astBuilder.NewExprStatement(p.parseBlock())
	}

	return p.astBuilder.NewExprStatement(p.parseExpression())
}

func (p *Parser) parseExpression() ast.Expression {
	return p.parseBinaryExpressions(PREC_SET)
}

// precedence climbing algorithm to make sure we treat precedence and associativity correctly
func (p *Parser) parseBinaryExpressions(minPrecedence Precedence) ast.Expression {
	left := p.parseUnaryExpression()
	var right ast.Expression

	for {
		precedence, assoc := precedenceFor(p.currentToken.Type)
		if precedence < minPrecedence {
			break
		}
		tok := p.currentToken
		p.advance()

		// now we climb the precedence ladder
		if assoc == ASSOC_LEFT {
			right = p.parseBinaryExpressions(precedence + 1)
		} else {
			right = p.parseBinaryExpressions(precedence)
		}

		left = p.astBuilder.NewBinOp(tok.Location, ast.TokenToBinaryOp(*tok), left, right)
	}
	return left
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	if p.match(lexer.TOKEN_BANG) || p.match(lexer.TOKEN_MINUS) || p.match(lexer.TOKEN_PLUS) {
		tok := p.previousToken
		right := p.parseUnaryExpression()
		return p.astBuilder.NewUnaryOp(tok.Location, ast.TokenToUnaryOp(*tok), right)
	} else if p.match(lexer.TOKEN_LPAREN) {
		expr := p.parseBinaryExpressions(PREC_SET)
		p.consume(lexer.TOKEN_RPAREN, "expected closing parenthesis")
		return expr
	} else if p.match(lexer.TOKEN_IDENTIFIER) {
		fmt.Printf("Parsing identifier - not yes supported\n")
		return nil
	} else if p.matchLiteral() {
		return p.astBuilder.NewLiteralExpression(*p.previousToken, p.previousToken.Location)
	} else if p.match(lexer.TOKEN_EOF) {
		p.errorAtCurrent(ParseErrorIdUnexpectedEOF, "unexpected end of input", nil)
		return nil
	}
	p.errorAtCurrent(ParseErrorIdUnexpectedToken, "expected unary expression", nil)
	return nil
}

func (p *Parser) matchLiteral() bool {
	literals := []lexer.TokenType{
		lexer.TOKEN_INTEGER,
		lexer.TOKEN_FLOAT,
		lexer.TOKEN_STRING,
		lexer.TOKEN_TRUE,
		lexer.TOKEN_FALSE,
		lexer.TOKEN_CHARACTER,
	}

	for _, literal := range literals {
		if p.match(literal) {
			return true
		}
	}
	return false
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

func (p *Parser) errorAtCurrent(id ParseErrorId, message string) {
	p.errorAt(*p.currentToken, id, message)
}

// func (p *Parser) errorAtPrevious(id ParseErrorId, message string, cause error) {
// 	p.errorAt(*p.previousToken, id, message, cause)
// }

func (p *Parser) errorAt(token token.Token, id ParseErrorId, message string) {
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
		Cause:   nil,
	})
}
