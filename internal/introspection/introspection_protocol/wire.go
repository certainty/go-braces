package introspection_protocol

// the on the wire representation of all data
// we might add meta data later

type WireRequest struct {
	Request
}

type WireResponse struct {
	Response
}

type WireEvent struct {
	Event
}
