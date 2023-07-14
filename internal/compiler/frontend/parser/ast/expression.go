package ast

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/frontend/lexer"
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

func (b BinaryExpression) IsNumeric() bool {
	switch b.Operator {
	case BinOpAdd, BinOpSub, BinOpMul, BinOpDiv, BinOpPow, BinOpMod:
		return true
	default:
		return false
	}
}

func (b BinaryExpression) IsBoolean() bool {
	switch b.Operator {
	case BinOpAnd, BinOpOr:
		return true
	default:
		return false
	}
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

func (l LiteralExpression) String() string {
	return fmt.Sprintf("Lit{ %s }[%d:%d]", l.Value, l.Location().Line, l.Location().StartOffset)
}
