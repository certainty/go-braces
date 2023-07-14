package ast

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/lexer"
	"github.com/certainty/go-braces/internal/compiler/location"
)

type (
	NodeId uint64

	AST struct {
		Nodes []Node
	}

	Node interface {
		ID() NodeId
		Location() location.Location
	}

	TypedNode interface{ TypeDecl() TypeDecl }

	Declaration interface {
		Node
		declNode()
	}

	BadDeclaration struct {
		id NodeId
	}

	Expression interface {
		Node
		exprNode()
	}

	BadExpression struct {
		id NodeId
	}

	Statement interface {
		Node
		stmtNode()
	}

	BadStatement struct {
		id NodeId
	}

	Identifier struct {
		id       NodeId
		location location.Location
		Label    string
	}

	TypeDecl struct {
		id       NodeId
		location location.Location
		Name     Identifier
	}

	PackageDecl struct {
		id       NodeId
		location location.Location
		Name     Identifier
	}

	CallableDecl struct {
		id       NodeId
		location location.Location
		// the declared type
		TpeDecl TypeDecl

		IsProcedure bool
		Name        Identifier
		Arguments   []ArgumentDecl
		Body        Block
	}

	ArgumentDecl struct {
		id       NodeId
		location location.Location
		Name     Identifier
		TpeDecl  TypeDecl
	}

	Block struct {
		id       NodeId
		location location.Location
		Code     []Node
	}

	LiteralExpression struct {
		id       NodeId
		location location.Location
		Token    lexer.Token
		Value    interface{}
	}

	UnaryExpression struct {
		id       NodeId
		location location.Location
		Type     TypeDecl
		Operator UnaryOperator
		Operand  Expression
	}

	BinaryExpression struct {
		id       NodeId
		location location.Location
		Left     Expression
		Right    Expression
		Operator BinaryOperator
		Type     TypeDecl
	}
)

func (BadExpression) exprNode()       {}
func (Identifier) exprNode()          {}
func (Block) exprNode()               {}
func (UnaryExpression) exprNode()     {}
func (b BinaryExpression) exprNode()  {}
func (l LiteralExpression) exprNode() {}

func (BadStatement) stmtNode() {}

func (BadDeclaration) declNode() {}
func (CallableDecl) declNode()   {}
func (ArgumentDecl) declNode()   {}
func (PackageDecl) declNode()    {}

func (e BadExpression) ID() NodeId     { return e.id }
func (e Identifier) ID() NodeId        { return e.id }
func (e Block) ID() NodeId             { return e.id }
func (e UnaryExpression) ID() NodeId   { return e.id }
func (e BinaryExpression) ID() NodeId  { return e.id }
func (e LiteralExpression) ID() NodeId { return e.id }
func (e BadStatement) ID() NodeId      { return e.id }
func (e BadDeclaration) ID() NodeId    { return e.id }
func (e CallableDecl) ID() NodeId      { return e.id }
func (e ArgumentDecl) ID() NodeId      { return e.id }
func (e PackageDecl) ID() NodeId       { return e.id }
func (e TypeDecl) ID() NodeId          { return e.id }

func (id Identifier) Location() location.Location       { return id.location }
func (d TypeDecl) Location() location.Location          { return d.location }
func (d PackageDecl) Location() location.Location       { return d.location }
func (def ArgumentDecl) Location() location.Location    { return def.location }
func (b Block) Location() location.Location             { return b.location }
func (d CallableDecl) Location() location.Location      { return d.location }
func (l LiteralExpression) Location() location.Location { return l.location }
func (u UnaryExpression) Location() location.Location   { return u.location }
func (b BinaryExpression) Location() location.Location  { return b.location }

func (d CallableDecl) TypeDecl() TypeDecl { return d.TpeDecl }
func (d ArgumentDecl) TypeDecl() TypeDecl { return d.TpeDecl }

func New() *AST {
	return &AST{
		Nodes: []Node{},
	}
}

func (ast *AST) ASTring() string {
	return NewASTWriter().Write(ast)
}
