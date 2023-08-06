package token_test

import (
	"fmt"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	"testing"
)

func ExampleNew() {
	origin := token.NewStringOrigin("example")
	loc := token.NewLocation(origin, token.Line(10), token.Column(20), token.From(30), token.To(40))

	tk := token.New(loc, token.IDENTIFIER, []rune("exampleIdentifier"))
	fmt.Println(tk.String())
	// Output: (IDENTIFIER, exampleIdentifier)
}

func ExampleByKeyword() {
	origin := token.NewStringOrigin("example")
	loc := token.NewLocation(origin, token.Line(10), token.Column(20), token.From(30), token.To(40))

	tk := token.ByKeyword(loc, "for")
	fmt.Println(tk.String())
	// Output: (FOR, for)
}

// TestTokenTypeString tests the String method for TokenType.
func TestTokenTypeString(t *testing.T) {
	tt := token.IDENTIFIER
	if tt.String() != "IDENTIFIER" {
		t.Errorf("Expected IDENTIFIER, got %s", tt.String())
	}
}

// TestTokenIsKeyword tests the IsKeyword method for Token.
func TestTokenIsKeyword(t *testing.T) {
	tk := token.New(token.Location{}, token.PACKAGE, []rune("package"))
	if !tk.IsKeyword() {
		t.Errorf("Expected true, got %v", tk.IsKeyword())
	}
}

// TestTokenIsLiteral tests the IsLiteral method for Token.
func TestTokenIsLiteral(t *testing.T) {
	tl := token.New(token.Location{}, token.STRING, []rune("hello"), "hello")
	if !tl.IsLiteral() {
		t.Errorf("Expected true, got %v", tl.IsLiteral())
	}
}

// TestTokenIsOperator tests the IsOperator method for Token.
func TestTokenIsOperator(t *testing.T) {
	to := token.New(token.Location{}, token.ADD, []rune("+"))
	if !to.IsOperator() {
		t.Errorf("Expected true, got %v", to.IsOperator())
	}
}

// TestTokenIsEOF tests the IsEOF method for Token.
func TestTokenIsEOF(t *testing.T) {
	teof := token.New(token.Location{}, token.EOF, nil)
	if !teof.IsEOF() {
		t.Errorf("Expected true, got %v", teof.IsEOF())
	}
}

// TestTokenIsIllegal tests the IsIllegal method for Token.
func TestTokenIsIllegal(t *testing.T) {
	ti := token.Illegal(token.Location{}, []rune("unknown"))
	if !ti.IsIllegal() {
		t.Errorf("Expected true, got %v", ti.IsIllegal())
	}
}

// TestTokenString tests the String method for Token.
func TestTokenString(t *testing.T) {
	tl := token.New(token.Location{}, token.STRING, []rune("hello"), "hello")
	if tl.String() != "(STRING, hello)" {
		t.Errorf("Expected (STRING, hello), got %s", tl.String())
	}
}

// TestTokenIllegal tests the creation of an Illegal token.
func TestTokenIllegal(t *testing.T) {
	ti := token.Illegal(token.Location{}, []rune("unknown"))
	if !ti.IsIllegal() || string(ti.Text) != "unknown" {
		t.Errorf("Expected illegal token with text unknown, got %s with text %s", ti.Type.String(), string(ti.Text))
	}
}

// TestTokenByKeyword tests the creation of a Token by keyword.
func TestTokenByKeyword(t *testing.T) {
	tb := token.ByKeyword(token.Location{}, "package")
	if !tb.IsKeyword() || tb.Type.String() != "PACKAGE" || string(tb.Text) != "package" {
		t.Errorf("Expected keyword token with text package, got %s with text %s", tb.Type.String(), string(tb.Text))
	}
}

func TestTokenByUnknownKeyword(t *testing.T) {
	tb := token.ByKeyword(token.Location{}, "unknown")

	if !tb.IsIllegal() || string(tb.Text) != "unknown" {
		t.Errorf("Expected illegal token with text unknown, got %s with text %s", tb.Type.String(), string(tb.Text))
	}
}
