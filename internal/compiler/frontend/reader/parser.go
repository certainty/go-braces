package reader

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/location"
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/prataprc/goparsec"
)

type Parser struct {
	introspectionAPI introspection.API
}

func NewParser(introspectionAPI introspection.API) *Parser {
	return &Parser{introspectionAPI: introspectionAPI}
}

func RecoverErrors(parser parsec.Parser, errors *[]error) parsec.Parser {
	return func(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
		result, newScanner := parser(s)

		if result == nil {
			// failed to parse, so we skip the scanning to till the next strategic choice and continue there
			// first get the error
			// now skip till the next strategic choice
			// finally restart the scan
			return nil, newScanner
		} else {
			return result, newScanner
		}
	}
}

// Parses a single datum
func (p *Parser) Parse(source *[]byte) (*DatumAST, error) {
	var nodes []Datum
	collectNodes := func(nodes []parsec.ParsecNode) parsec.ParsecNode {
		for _, node := range nodes {
			nodes = append(nodes, node.(Datum))
		}
		return nodes
	}
	scanner := parsec.NewScanner(*source).TrackLineno()
	parser := parsec.Many(collectNodes, parseDatum)

	_, err := parser(scanner)
	if err != nil {
		return nil, fmt.Errorf("error parsing datum")
	} else {
		return &DatumAST{Data: nodes}, nil
	}
}

func parseDatum(scanner parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	return parseBoolean(scanner)
}

func parseBoolean(scanner parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	boolParser := parsec.OrdChoice(nil, parsec.Token("#t", "TOK_TRUE"), parsec.Token("#f", "TOK_FALSE"))
	node, newScanner := boolParser(scanner)

	if node != nil {
		loc := buildLocation(scanner, newScanner)
		terminal := node.(parsec.Terminal)
		return NewDatumBool(terminal.GetName() == "TOK_TRUE", loc), newScanner
	} else {
		return nil, newScanner
	}
}

func buildLocation(before parsec.Scanner, after parsec.Scanner) location.Location {
	beforeCursor := before.GetCursor()
	afterCursor := after.GetCursor()

	return location.Location{
		Input:       nil,
		Line:        before.Lineno(),
		StartOffset: beforeCursor,
		EndOffset:   afterCursor,
	}
}
