package common

import (
	"time"

	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
)

type MsgTick time.Time

type MsgIntrospectionEvent struct {
	Event compiler_introspection.CompilerIntrospectionEvent
}

func NewMsgIntrospectionEvent(event compiler_introspection.CompilerIntrospectionEvent) MsgIntrospectionEvent {
	return MsgIntrospectionEvent{Event: event}
}

type MsgError struct {
	Err error
}

func NewMsgError(err error) MsgError {
	return MsgError{Err: err}
}

type MsgResize struct {
	Width, Height int
}

func NewMsgResize(width, height int) MsgResize {
	return MsgResize{Width: width, Height: height}
}

type MsgModeChange struct {
	ActiveMode Mode
}

func NewMsgModeChange(activeMode Mode) MsgModeChange {
	return MsgModeChange{ActiveMode: activeMode}
}

type MsgClientConnected bool

type MsgRequestStatus struct {
	RequestStatus RequestStatus
}

func NewMsgRequestStatus(requestStatus RequestStatus) MsgRequestStatus {
	return MsgRequestStatus{RequestStatus: requestStatus}
}

type MsgActivateKeyMap KeyMap
