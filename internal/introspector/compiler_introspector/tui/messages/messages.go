package messages

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"time"
)

type TickMsg time.Time

type RequestSentMsg struct{}
type RequestErrorMsg struct {
	Err error
}

type IntrospectionEventMsg struct {
	Event compiler_introspection.CompilerIntrospectionEvent
}

type IntrospectionConnectionError struct {
	Err error
}

type IntrospectionConnectionEstablished struct{}
