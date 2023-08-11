package introspection

import "net"

type WireClient[ControlType any, EventType any] struct {
	scope string
}

func NewWireClient[C any, E any](scope string) WireClient[C, E] {
	return WireClient[C, E]{scope: scope}
}

func (w WireClient[C, E]) Connect() (*WireControlConnection[C], *WireEventConnection[E], error) {
	controlConn, err := net.Dial("unix", ControlSocketPath(w.scope))
	if err != nil {
		return nil, nil, err
	}
	wireControlConnection := NewWireControlConnection[C](controlConn)

	eventsConn, err := net.Dial("unix", EventSocketPath(w.scope))
	if err != nil {
		wireControlConnection.Close()
		return nil, nil, err
	}
	wireEventConnection := NewWireEventConnection[E](eventsConn, WireEventSink)

	return wireControlConnection, wireEventConnection, nil
}
