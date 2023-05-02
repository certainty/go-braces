package parser

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/frontend/reader"
	"github.com/certainty/go-braces/internal/compiler/location"
)

type SchemeExpression interface {
	Location() location.Location
}

type LiteralExpression struct {
	Datum reader.Datum
}

func (l LiteralExpression) String() string {
	return fmt.Sprintf("Lit{ %s }[%d:%d]", l.Datum, l.Location().Line, l.Location().StartOffset)
}

func (l LiteralExpression) Location() location.Location {
	return l.Datum.Location()
}
