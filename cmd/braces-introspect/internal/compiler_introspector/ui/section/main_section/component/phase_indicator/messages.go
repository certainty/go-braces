package phase_indicator

import "github.com/certainty/go-braces/pkg/introspection/compiler_introspection"

type MsgPhase compiler_introspection.CompilationPhase
type MsgReset struct{}
type MsgFinish struct{}
