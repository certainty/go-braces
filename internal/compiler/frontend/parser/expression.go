package parser

import (
	"fmt"
	"github.com/certainty/go-braces/internal/compiler/location"
	"github.com/certainty/go-braces/internal/isa"
)

type SchemeExpression interface {
	Location() location.Location
}

type LiteralExpression struct {
	Datum isa.Datum
}

func (l LiteralExpression) String() string {
	return fmt.Sprintf("Lit{ %s }[%d:%d]", l.Datum, l.Location().Line, l.Location().StartOffset)
}

func (l LiteralExpression) Location() location.Location {
	return l.Datum.Location()
}
