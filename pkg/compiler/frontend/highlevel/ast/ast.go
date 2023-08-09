// Representation of the AST for the high level language.
// This package provides all the necessary datatypes and functionality to represent and work with
// the parse result of the high level parser.
//
// The AST allows to represent nodes, which didn't parse successfully. These are represented as special Bad* nodes.
// You can use the AST walker to traverse the AST and find all the Bad* nodes, i.e. for error reporting.
package ast

import (
	"github.com/certainty/go-braces/pkg/compiler/frontend/astutils"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
)

type (
	Node interface {
		astutils.NodeIded
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
		id       astutils.NodeId
		location token.Location
	}

	BasicLitExpr struct {
		id    astutils.NodeId
		Token token.Token
	}

	ParenExpr struct {
		id   astutils.NodeId
		Expr Expression
	}

	BlockExpr struct {
		id         astutils.NodeId
		location   token.Location
		Statements []Statement
	}

	UnaryExpr struct {
		id   astutils.NodeId
		Op   token.Token
		Expr Expression
	}

	BinaryExpr struct {
		id    astutils.NodeId
		Op    token.Token
		Left  Expression
		Right Expression
	}

	Identifier struct {
		id       astutils.NodeId
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

func (e BadExpr) ID() astutils.NodeId      { return e.id }
func (e BasicLitExpr) ID() astutils.NodeId { return e.id }
func (e ParenExpr) ID() astutils.NodeId    { return e.id }
func (e BlockExpr) ID() astutils.NodeId    { return e.id }
func (e UnaryExpr) ID() astutils.NodeId    { return e.id }
func (e BinaryExpr) ID() astutils.NodeId   { return e.id }
func (e Identifier) ID() astutils.NodeId   { return e.id }

func (e BadExpr) Location() token.Location      { return e.location }
func (e BasicLitExpr) Location() token.Location { return e.Token.Location }
func (e ParenExpr) Location() token.Location    { return e.Expr.Location() }
func (e BlockExpr) Location() token.Location    { return e.location }
func (e UnaryExpr) Location() token.Location    { return e.Op.Location }
func (e BinaryExpr) Location() token.Location   { return e.Op.Location }
func (e Identifier) Location() token.Location   { return e.location }

func (e BasicLitExpr) Value() interface{} {
	return e.Token.LitValue
}

// //////////////////////////////////////////////////
// statements
// //////////////////////////////////////////////////
type (
	BadStmt struct {
		id       astutils.NodeId
		location token.Location
	}

	ExprStmt struct {
		Expr Expression
	}
)

func (BadStmt) stmtNode()  {}
func (ExprStmt) stmtNode() {}

func (s BadStmt) ID() astutils.NodeId  { return s.id }
func (s ExprStmt) ID() astutils.NodeId { return s.Expr.ID() }

func (s BadStmt) Location() token.Location  { return s.location }
func (s ExprStmt) Location() token.Location { return s.Expr.Location() }

// //////////////////////////////////////////////////
// declarations
// //////////////////////////////////////////////////
type (
	BadDecl struct {
		id       astutils.NodeId
		location token.Location
	}

	TypeSpec struct {
		id       astutils.NodeId
		location token.Location
		Name     *Identifier
	}

	ProcDecl struct {
		id       astutils.NodeId
		location token.Location
		Name     *Identifier
		Type     *ProcType
		Body     *BlockExpr
	}
)

func (BadDecl) declNode()  {}
func (TypeSpec) declNode() {}
func (ProcDecl) declNode() {}

func (d BadDecl) ID() astutils.NodeId  { return d.id }
func (d TypeSpec) ID() astutils.NodeId { return d.id }
func (d ProcDecl) ID() astutils.NodeId { return d.id }

func (d BadDecl) Location() token.Location  { return d.location }
func (d TypeSpec) Location() token.Location { return d.location }
func (d ProcDecl) Location() token.Location { return d.location }

// //////////////////////////////////////////////////
// types
// //////////////////////////////////////////////////
type (
	Field struct {
		id   astutils.NodeId
		Name *Identifier
		Type *TypeSpec
	}

	ProcType struct {
		id     astutils.NodeId
		Params []*Field
		Result *TypeSpec
	}
)

func (f Field) ID() astutils.NodeId    { return f.id }
func (t ProcType) ID() astutils.NodeId { return t.id }

func (f Field) Location() token.Location    { return f.Name.Location() }
func (t ProcType) Location() token.Location { return t.Result.Location() }

// //////////////////////////////////////////////////
// Source & Package
// //////////////////////////////////////////////////

type Source struct {
	id           astutils.NodeId
	Declarations []Declaration
}

func (s Source) ID() astutils.NodeId { return s.id }
func (s Source) Location() token.Location {
	if len(s.Declarations) > 0 {
		return s.Declarations[0].Location()
	}
	return token.Location{}
}
