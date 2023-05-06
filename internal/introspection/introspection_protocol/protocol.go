package introspection_protocol

import "encoding/gob"

type Event interface{}
type Request interface{}
type Response interface{}

type IntrospectionType uint8

const (
	IntrospectionTypeCompiler IntrospectionType = iota
	IntrospectionTypeVM
)

type HeloRequest struct {
	IntrospectionType IntrospectionType
}

type HeloResponse struct {
	ClientID string
}

func RegisterTypes() {
	gob.Register(HeloRequest{})
	gob.Register(HeloResponse{})
}
