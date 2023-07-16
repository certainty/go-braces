package ir

import "github.com/certainty/go-braces/pkg/shared/isa"
import "github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/types"

type Value struct {
	Type  types.Type
	Value isa.Value
}
