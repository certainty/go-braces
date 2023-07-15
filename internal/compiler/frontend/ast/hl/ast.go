// Representation of the AST for the high level language.
// This package provides all the necessary datatypes and functionality to represent and work with
// the parse result of the high level parser.
//
// The AST allows to represent nodes, which didn't parse successfully. These are represented as special Bad* nodes.
// You can use the AST walker to traverse the AST and find all the Bad* nodes, i.e. for error reporting.
package ast

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/token"
)

type (
	NodeId uint64

	Node interface {
		ID() NodeId
		Location() token.Location
	}

	Expression interface {
		Node
		exprNode()
	}

	Declaration interface {
		Node
		declNode()
	}

	Statement interface {
		Node
		stmtNode()
	}
)

// //////////////////////////////////////////////////
// Expressions
// //////////////////////////////////////////////////
type (
	BadExpr struct {
		id NodeId
	}

	BasicLitExpr struct {
		id    NodeId
		Token token.Token
	}

	ParenExpr struct {
		id   NodeId
		Expr Expression
	}

	BlockExpr struct {
		id         NodeId
		location   token.Location
		Statements []Statement
	}

	UnaryExpr struct {
		id   NodeId
		Op   token.Token
		Expr Expression
	}

	BinaryExpr struct {
		id    NodeId
		Op    token.Token
		Left  Expression
		Right Expression
	}

	Identifier struct {
		id       NodeId
		location token.Location
		Name     string
	}
)

func (BadExpr) exprNode()      {}
func (BasicLitExpr) exprNode() {}
func (ParenExpr) exprNode()    {}
func (BlockExpr) exprNode()    {}
func (UnaryExpr) exprNode()    {}
func (BinaryExpr) exprNode()   {}
func (Identifier) exprNode()   {}

func (e BadExpr) ID() NodeId      { return e.id }
func (e BasicLitExpr) ID() NodeId { return e.id }
func (e ParenExpr) ID() NodeId    { return e.id }
func (e BlockExpr) ID() NodeId    { return e.id }
func (e UnaryExpr) ID() NodeId    { return e.id }
func (e BinaryExpr) ID() NodeId   { return e.id }
func (e Identifier) ID() NodeId   { return e.id }

func (e BadExpr) Location() token.Location      { return token.Location{} }
func (e BasicLitExpr) Location() token.Location { return e.Token.Location }
func (e ParenExpr) Location() token.Location    { return e.Expr.Location() }
func (e BlockExpr) Location() token.Location    { return e.location }
func (e UnaryExpr) Location() token.Location    { return e.Op.Location }
func (e BinaryExpr) Location() token.Location   { return e.Op.Location }
func (e Identifier) Location() token.Location   { return e.location }

// //////////////////////////////////////////////////
// statements
// //////////////////////////////////////////////////
type (
	BadStmt struct {
		id NodeId
	}

	ExprStmt struct {
		Expr Expression
	}
)

func (BadStmt) stmtNode()  {}
func (ExprStmt) stmtNode() {}

func (s BadStmt) ID() NodeId  { return s.id }
func (s ExprStmt) ID() NodeId { return s.Expr.ID() }

func (s BadStmt) Location() token.Location  { return token.Location{} }
func (s ExprStmt) Location() token.Location { return s.Expr.Location() }

// //////////////////////////////////////////////////
// declarations
// //////////////////////////////////////////////////
type (
	BadDecl struct {
		id NodeId
	}

	TypeSpec struct {
		id       NodeId
		location token.Location
		Name     Identifier
	}

	ProcDecl struct {
		id       NodeId
		location token.Location
		Name     Identifier
		Type     ProcType
		Body     BlockExpr
	}
)

func (BadDecl) declNode()  {}
func (TypeSpec) declNode() {}
func (ProcDecl) declNode() {}

func (d BadDecl) ID() NodeId  { return d.id }
func (d TypeSpec) ID() NodeId { return d.id }
func (d ProcDecl) ID() NodeId { return d.id }

func (d BadDecl) Location() token.Location  { return token.Location{} }
func (d TypeSpec) Location() token.Location { return d.location }
func (d ProcDecl) Location() token.Location { return d.location }

// //////////////////////////////////////////////////
// types
// //////////////////////////////////////////////////
type (
	Field struct {
		Id   Identifier
		Name Identifier
		Type *TypeSpec
	}

	ProcType struct {
		Params []Field
		Result *TypeSpec
	}
)

// //////////////////////////////////////////////////
// Source & Package
// //////////////////////////////////////////////////

type Source struct {
	Declarations []Declaration
	Statements   []Node
}
