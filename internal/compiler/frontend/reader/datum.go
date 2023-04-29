package reader

import "github.com/certainty/go-braces/internal/compiler/location"

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

func (d *DatumBool) Location() location.Location {
	return d.Loc
}
