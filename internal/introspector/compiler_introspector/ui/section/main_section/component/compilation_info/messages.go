package compilation_info

import "github.com/certainty/go-braces/internal/compiler/location"

type MsgNewCompilation struct {
	Options []string
	Origin  location.Origin
}
