package token_test

import (
	"fmt"
	"github.com/certainty/go-braces/internal/compiler/frontend/token"
	"testing"
)

func ExampleNewLocation() {
	origin := token.NewStringOrigin("example")
	location := token.NewLocation(origin, token.Line(10), token.Column(20), token.From(30), token.To(40))
	fmt.Println(location.String())
	// Output: string://example:10:20
}

func ExmapleNewStringOrigin() {
	origin := token.NewStringOrigin("example")
	fmt.Println(origin.Name())
	// Output: example
}

func TestStringOrigin(t *testing.T) {
	origin := token.NewStringOrigin("Test")

	if origin.Name() != "Test" {
		t.Errorf("Unexpected name for StringOrigin, got: %s, want: Test", origin.Name())
	}

	if origin.Description() != "string://Test" {
		t.Errorf("Unexpected description for StringOrigin, got: %s, want: string://Test", origin.Description())
	}
}

func TestLocation(t *testing.T) {
	origin := token.NewStringOrigin("Test")
	location := token.NewLocation(origin, token.Line(10), token.Column(20), token.From(30), token.To(40))

	if location.String() != "string://Test:10:20" {
		t.Errorf("Unexpected string for Location, got: %s, want: string://Test:10:20", location.String())
	}
}
