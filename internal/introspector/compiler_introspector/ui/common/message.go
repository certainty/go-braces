package common

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"time"
)

type MsgTick time.Time

type MsgIntrospectionEvent struct {
	Event compiler_introspection.CompilerIntrospectionEvent
}

type MsgError struct {
	Err error
}

type MsgResize struct {
	Width, Height int
}

type MsgModeChange struct {
	ActiveMode Mode
}

type MsgClientConnected bool

type MsgRequestStatus struct {
	RequestStatus RequestStatus
}
