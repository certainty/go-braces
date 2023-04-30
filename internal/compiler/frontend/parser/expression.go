package parser

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/reader"
	"github.com/certainty/go-braces/internal/compiler/location"
)

type SchemeExpression interface {
	Location() location.Location
}

type LiteralExpression struct {
	Datum reader.Datum
}

func (l LiteralExpression) Location() location.Location {
	return l.Datum.Location()
}
