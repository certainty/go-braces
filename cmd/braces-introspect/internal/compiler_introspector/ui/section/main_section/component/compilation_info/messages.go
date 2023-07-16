package compilation_info

import "github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"

type MsgNewCompilation struct {
	Options []string
	Origin  token.Origin
}
