package lexer_test

import (
	"path/filepath"
	"testing"

	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/lexer"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
)

var fixtureDir = "../../../../testdata/compiler"

func TestMinimal(t *testing.T) {
	source, err := lexer.NewFileInput(filepath.Join(fixtureDir, "minimal.brace"))
	assert.NoError(t, err)

	scanner := lexer.New(source)
	tokens := readTokens(scanner)
	assert.Empty(t, scanner.Errors)

	snaps.MatchSnapshot(t, tokens)
}

func readTokens(scanner *lexer.Scanner) []token.Token {
	tokens := []token.Token{}

	for {
		tok := scanner.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}

	return tokens
}
