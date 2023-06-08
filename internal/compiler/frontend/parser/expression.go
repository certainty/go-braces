package parser

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/location"
	"github.com/certainty/go-braces/internal/isa"
)

type Expression interface {
	Location() location.Location
}

type LiteralExpression struct {
	Value    isa.Value
	location location.Location
}

func (l LiteralExpression) String() string {
	return fmt.Sprintf("Lit{ %s }[%d:%d]", l.Value, l.Location().Line, l.Location().StartOffset)
}

func (l LiteralExpression) Location() location.Location {
	return l.location
}
