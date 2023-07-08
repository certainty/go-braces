package ir

import "github.com/certainty/go-braces/internal/isa"

type IRValue struct {
	Type  IRType
	Value isa.Value
}
