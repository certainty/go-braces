package introspection

type IntrospectionType uint8

const (
	IntrospectionTypeCompiler IntrospectionType = iota
	IntrospectionTypeVM
)

type Request interface{}

type Response interface{}

type WireRequest struct {
	Request
}

type WireResponse struct {
	Response
}

type WireEvent struct {
	IntrospectionEvent
}

type EventSingleStepBarrierReached struct{}

// Requests / Responses
type HeloRequest struct {
	IntrospectionType IntrospectionType
}

type HeloResponse struct {
	ClientID string
}

type StartSingleStepRequest struct{}

type StartSingleStepResponse struct{}

type NextSingleStepRequest struct{}

type NextSingleStepResponse struct{}

type ContinueSingleStepRequest struct{}

type ContinueSingleStepResponse struct{}
