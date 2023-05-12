package isa

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/location"
)

type Datum interface {
	Location() location.Location
}

func NewDatumBool(value bool, loc location.Location) DatumBool {
	return DatumBool{Value: value, Loc: loc}
}

type DatumBool struct {
	Value bool
	Loc   location.Location
}

func (d DatumBool) Location() location.Location {
	return d.Loc
}

func (d DatumBool) String() string {
	return fmt.Sprintf("<%v>[%d:%d]", d.Value, d.Location().Line, d.Location().StartOffset)
}