package compilation_info

import "github.com/certainty/go-braces/internal/compiler/frontend/token"

type MsgNewCompilation struct {
	Options []string
	Origin  token.Origin
}
