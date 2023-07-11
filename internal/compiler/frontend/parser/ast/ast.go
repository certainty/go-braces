package ast

import (
	"github.com/certainty/go-braces/internal/compiler/location"
)

type AST struct {
	Nodes []Node
}

type Node interface {
	Location() location.Location
}

type TypedNode interface {
	TypeName() string
}

type Declaration interface{ Node }
type Expression interface{ Node }
type Statement interface{ Node }

type Identifier struct {
	ID       string
	location location.Location
}

var _ Node = (*Identifier)(nil)

func NewIdentifier(id string, location location.Location) Identifier {
	return Identifier{
		ID:       id,
		location: location,
	}
}

func (id Identifier) Location() location.Location {
	return id.location
}

var _ Node = (*Identifier)(nil)
var _ Declaration = (*Identifier)(nil)

type TypeDecl struct {
	Name     Identifier
	location location.Location
}

func NewTypeDecl(name Identifier, location location.Location) TypeDecl {
	return TypeDecl{
		Name:     name,
		location: location,
	}
}

func (d TypeDecl) Location() location.Location {
	return d.location
}

var _ Node = (*TypeDecl)(nil)
var _ Declaration = (*TypeDecl)(nil)

type PackageDecl struct {
	Name     Identifier
	location location.Location
}

func NewPackageDecl(name Identifier, location location.Location) PackageDecl {
	return PackageDecl{
		Name:     name,
		location: location,
	}
}

func (d PackageDecl) Location() location.Location {
	return d.location
}

var _ Node = (*PackageDecl)(nil)
var _ Declaration = (*PackageDecl)(nil)

type CallableDecl struct {
	Type        TypeDecl
	IsProcedure bool
	Name        Identifier
	Arguments   []ArgumentDecl
	Body        Block
	location    location.Location
}

func NewFunctionDecl(
	t TypeDecl,
	name Identifier,
	arguments []ArgumentDecl,
	body Block,
	location location.Location,
) CallableDecl {
	return CallableDecl{
		Type:        t,
		IsProcedure: false,
		Name:        name,
		Arguments:   arguments,
		Body:        body,
		location:    location,
	}
}

func NewProcedureDecl(
	t TypeDecl,
	name Identifier,
	arguments []ArgumentDecl,
	body Block,
	location location.Location,
) CallableDecl {
	return CallableDecl{
		Type:        t,
		IsProcedure: true,
		Name:        name,
		Arguments:   arguments,
		Body:        body,
		location:    location,
	}
}

func (d CallableDecl) Location() location.Location {
	return d.location
}

func (d CallableDecl) TypeName() string {
	return d.Type.Name.ID
}

var _ Node = (*CallableDecl)(nil)
var _ Declaration = (*CallableDecl)(nil)
var _ TypedNode = (*CallableDecl)(nil)

type ArgumentDecl struct {
	Name     Identifier
	Type     TypeDecl
	location location.Location
}

func NewArgumentDecl(name Identifier, t TypeDecl, location location.Location) ArgumentDecl {
	return ArgumentDecl{
		Name:     name,
		Type:     t,
		location: location,
	}
}

func (def ArgumentDecl) Location() location.Location {
	return def.location
}

func (def ArgumentDecl) TypeName() string {
	return def.Type.Name.ID
}

var _ Node = (*ArgumentDecl)(nil)
var _ Declaration = (*ArgumentDecl)(nil)
var _ TypedNode = (*ArgumentDecl)(nil)

type Block struct {
	Code     []Node
	location location.Location
}

func NewBlock(body []Node, location location.Location) Block {
	return Block{
		Code:     body,
		location: location,
	}
}

func (b Block) Location() location.Location {
	return b.location
}

var _ Node = (*Block)(nil)
var _ Expression = (*Block)(nil)

func New() *AST {
	return &AST{
		Nodes: []Node{},
	}
}

func (ast *AST) ASTring() string {
	return NewASTWriter().Write(ast)
}

func (ast *AST) AddExpression(expression Expression) {
	ast.Nodes = append(ast.Nodes, expression)
}
