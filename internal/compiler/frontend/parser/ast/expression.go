package ast

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/frontend/lexer"
	"github.com/certainty/go-braces/internal/compiler/location"
)

type UnaryOperator uint8

const (
	UnaryOpNot UnaryOperator = iota
	UnaryOpNeg
	UnaryOpPos
)

type BinaryOperator uint8

const (
	BinOpAdd BinaryOperator = iota
	BinOpSub
	BinOpMul
	BinOpDiv
	BinOpPow
	BinOpMod
	BinOpAnd
	BinOpOr
)

type UnaryExpression struct {
	location location.Location
	Operator UnaryOperator
	Operand  Expression
	Type     TypeDecl
}

var _ Node = (*UnaryExpression)(nil)
var _ Expression = (*UnaryExpression)(nil)

func UnaryOp(location location.Location, operator UnaryOperator, operand Expression) UnaryExpression {
	return UnaryExpression{
		location: location,
		Operator: operator,
		Operand:  operand,
	}
}

func (u UnaryExpression) Location() location.Location {
	return u.location
}

func TokenToUnaryOp(tok lexer.Token) UnaryOperator {
	switch tok.Type {
	case lexer.TOKEN_BANG:
		return UnaryOpNot
	case lexer.TOKEN_MINUS:
		return UnaryOpNeg
	case lexer.TOKEN_PLUS:
		return UnaryOpPos
	default:
		panic(fmt.Sprintf("invalid unary operator %v", tok))
	}
}

type BinaryExpression struct {
	location location.Location
	Left     Expression
	Right    Expression
	Operator BinaryOperator
	Type     TypeDecl
}

var _ Node = (*BinaryExpression)(nil)
var _ Expression = (*BinaryExpression)(nil)

func BinOp(location location.Location, operator BinaryOperator, left Expression, right Expression) BinaryExpression {
	return BinaryExpression{
		location: location,
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (b BinaryExpression) Location() location.Location {
	return b.location
}

func TokenToBinaryOp(token lexer.Token) BinaryOperator {
	switch token.Type {
	case lexer.TOKEN_PLUS:
		return BinOpAdd
	case lexer.TOKEN_MINUS:
		return BinOpSub
	case lexer.TOKEN_STAR:
		return BinOpMul
	case lexer.TOKEN_SLASH:
		return BinOpDiv
	case lexer.TOKEN_MOD:
		return BinOpMod
	case lexer.TOKEN_POWER:
		return BinOpPow
	case lexer.TOKEN_AMPERSAND_AMPERSAND:
		return BinOpAnd
	case lexer.TOKEN_PIPE_PIPE:
		return BinOpOr
	default:
		panic(fmt.Sprintf("invalid binary operator %v", token))
	}
}

type LiteralExpression struct {
	Token    lexer.Token
	Value    interface{}
	location location.Location
}

var _ Node = (*LiteralExpression)(nil)
var _ Expression = (*LiteralExpression)(nil)

func NewLiteralExpression(token lexer.Token, location location.Location) LiteralExpression {
	return LiteralExpression{
		Token:    token,
		Value:    token.Value,
		location: location,
	}
}

func (l LiteralExpression) String() string {
	return fmt.Sprintf("Lit{ %s }[%d:%d]", l.Value, l.Location().Line, l.Location().StartOffset)
}

func (l LiteralExpression) Location() location.Location {
	return l.location
}
