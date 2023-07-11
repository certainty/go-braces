package parser

import (
	"fmt"
	"log"

	"github.com/certainty/go-braces/internal/compiler/frontend/lexer"
	"github.com/certainty/go-braces/internal/compiler/frontend/parser/ast"
	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
)

type Precedence uint8
type Associativity uint8

const (
	PREC_NONE       Precedence = iota
	PREC_SET                   // set x expr
	PREC_OR                    // ||
	PREC_AND                   // &&
	PREC_EQUALITY              // == !=
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_FACTOR                // * / % /
	PREC_EXPONENT              // **
	PREC_UNARY                 // ! -
	PREC_CALL                  //  ()
	PREC_PRIMARY
)

const (
	ASSOC_NONE Associativity = iota
	ASSOC_LEFT
	ASSOC_RIGHT
)

func precedenceFor(tokenType lexer.TokenType) (Precedence, Associativity) {
	switch tokenType {
	case lexer.TOKEN_STAR, lexer.TOKEN_SLASH, lexer.TOKEN_MOD:
		return PREC_FACTOR, ASSOC_LEFT
	case lexer.TOKEN_POWER:
		return PREC_EXPONENT, ASSOC_RIGHT
	case lexer.TOKEN_PLUS, lexer.TOKEN_MINUS:
		return PREC_TERM, ASSOC_LEFT
	case lexer.TOKEN_EQUAL_EQUAL, lexer.TOKEN_BANG_EQUAL, lexer.TOKEN_LT, lexer.TOKEN_LT_EQUAL, lexer.TOKEN_GT, lexer.TOKEN_GT_EQUAL:
		return PREC_COMPARISON, ASSOC_LEFT
	case lexer.TOKEN_AMPERSAND_AMPERSAND:
		return PREC_AND, ASSOC_LEFT
	case lexer.TOKEN_PIPE_PIPE:
		return PREC_OR, ASSOC_LEFT
	case lexer.TOKEN_LPAREN:
		return PREC_NONE, ASSOC_LEFT
	default:
		return PREC_NONE, ASSOC_NONE
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

	println("parsing input: ", input.Source())

	// now start the parsing process
	p.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseParse)
	defer p.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseParse)

	return p.parseInput()
}

func (p *Parser) parseInput() (*ast.AST, error) {
	p.advance()
	p.parsePackageDeclaration()
	if p.hadError {
		log.Print("Had error, synchronizing")
		p.synchronize()
	}

	for {
		log.Print("Parsing next declaration")
		p.parseDeclaration()
		if p.currentToken.Type == lexer.TOKEN_EOF {
			break
		}
	}

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
	log.Printf("DEBUG consume current: %v prev: %v", p.currentToken, p.previousToken)
	if p.currentToken.Type == tokenType {
		p.advance()
		return
	}
	p.errorAtCurrent(ParseErrorIdUnexpectedToken, message, nil)
}

// Try to recover to next synchrnization point.
// These are:
// * statement boundaries
// * new blocks
// * function boundaries
func (p *Parser) synchronize() {
	p.panicMode = false

	for p.currentToken.Type != lexer.TOKEN_EOF {
		if p.previousToken.Type == lexer.TOKEN_SEMICOLON {
			return
		}

		switch p.currentToken.Type {
		case lexer.TOKEN_RBRACE, lexer.TOKEN_FUN, lexer.TOKEN_PROC, lexer.TOKEN_IF, lexer.TOKEN_FOR:
			return
		default:
			p.advance()
		}
	}
}

func (p *Parser) parseDeclaration() {
	switch p.currentToken.Type {
	case lexer.TOKEN_FUN:
		p.parseFunctionDeclaration()
	case lexer.TOKEN_PROC:
		p.parseProcedureDeclaration()
	case lexer.TOKEN_EOF:
		p.advance()
		return
	default:
		p.errorAtCurrent(ParseErrorIdUnexpectedToken, "expected declaration", nil)
	}

	if p.hadError {
		p.synchronize()
	}
}

func (p *Parser) parsePackageDeclaration() {
	log.Print("Parsing package declaration")

	p.consume(lexer.TOKEN_PACKAGE, "expected package")
	packageLocation := p.previousToken.Location

	packageName := p.parseIdentifier()
	packageDecl := ast.NewPackageDecl(packageName, packageLocation)

	log.Printf("Package name: %s", packageName.ID)
	log.Printf("Package location: %d %d", packageLocation.Line, packageLocation.StartOffset)

	p.ast.Nodes = append(p.ast.Nodes, packageDecl)
}

func (p *Parser) parseFunctionDeclaration() {
	p.consume(lexer.TOKEN_FUN, "expected fun")
	location := p.currentToken.Location
	funcName := p.parseIdentifier()
	args := p.parseArguments()
	tpe := p.parseTypeDeclaration()
	body := p.parseBlock()
	function := ast.NewFunctionDecl(tpe, funcName, args, body, location)
	p.ast.Nodes = append(p.ast.Nodes, function)

}

func (p *Parser) parseProcedureDeclaration() {
	p.consume(lexer.TOKEN_PROC, "expected proc")
	location := p.currentToken.Location
	procName := p.parseIdentifier()
	args := p.parseArguments()
	log.Printf("Arguments: %v", args)
	// optional return type
	tpe := p.parseTypeDeclaration()
	log.Printf("Type: %v", tpe)
	body := p.parseBlock()
	log.Printf("Body: %v", body)
	procedure := ast.NewProcedureDecl(tpe, procName, args, body, location)
	p.ast.Nodes = append(p.ast.Nodes, procedure)
}

func (p *Parser) parseArguments() []ast.ArgumentDecl {
	args := []ast.ArgumentDecl{}
	p.consume(lexer.TOKEN_LPAREN, "expected '('")
	log.Printf("Parsing arguments")
	for {
		log.Printf("Token: %v", p.currentToken)
		if p.check(lexer.TOKEN_RPAREN) {
			break
		}
		location := p.currentToken.Location
		argName := p.parseIdentifier()
		argType := p.parseTypeDeclaration()
		args = append(args, ast.NewArgumentDecl(argName, argType, location))
		if !p.match(lexer.TOKEN_COMMA) {
			break
		}
	}
	p.consume(lexer.TOKEN_RPAREN, "expected ')'")
	return args
}

func (p *Parser) parseTypeDeclaration() ast.TypeDecl {
	location := p.currentToken.Location
	if p.check(lexer.TOKEN_COLON) {
		p.consume(lexer.TOKEN_COLON, "expected ':'")
		return ast.NewTypeDecl(p.parseIdentifier(), location)
	} else {
		return ast.NewTypeDecl(ast.NewIdentifier("void", location), location)
	}
}

func (p *Parser) parseIdentifier() ast.Identifier {
	p.consume(lexer.TOKEN_IDENTIFIER, "expected identifier")
	return ast.NewIdentifier(string(p.previousToken.Text), p.previousToken.Location)
}

func (p *Parser) parseBlock() ast.Block {
	// TODO: track scope
	p.consume(lexer.TOKEN_LBRACE, "expected '{'")
	location := p.previousToken.Location
	block := ast.NewBlock([]ast.Node{}, location)

	for !p.match(lexer.TOKEN_RBRACE) {
		block.Code = append(block.Code, p.parseBlockStatment())
	}

	return block
}

func (p *Parser) parseBlockStatment() ast.Node {
	if p.match(lexer.TOKEN_LBRACE) {
		return p.parseBlock()
	}

	return p.parseExpression()
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

		left = ast.BinOp(tok.Location, ast.TokenToBinaryOp(*tok), left, right)
	}
	return left
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	if p.match(lexer.TOKEN_BANG) || p.match(lexer.TOKEN_MINUS) || p.match(lexer.TOKEN_PLUS) {
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
	} else if p.matchLiteral() {
		return ast.NewLiteralExpression(*p.previousToken, p.previousToken.Location)
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

func (p *Parser) errorAtCurrent(id ParseErrorId, message string, cause error) {
	p.errorAt(*p.currentToken, id, message, cause)
}

// func (p *Parser) errorAtPrevious(id ParseErrorId, message string, cause error) {
// 	p.errorAt(*p.previousToken, id, message, cause)
// }

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
