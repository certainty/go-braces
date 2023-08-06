package introspection

import "net"

type WireClient struct {
	scope string
}

func NewWireClient(scope string) WireClient {
	return WireClient{scope: scope}
}

func (w WireClient) Connect() (*WireControlConnection, *WireEventConnection, error) {
	controlConn, err := net.Dial("unix", ControlSocketPath(w.scope))
	if err != nil {
		return nil, nil, err
	}
	wireControlConnection := NewWireControlConnection(controlConn)

	eventsConn, err := net.Dial("unix", EventSocketPath(w.scope))
	if err != nil {
		wireControlConnection.Close()
		return nil, nil, err
	}
	wireEventConnection := NewWireEventConnection(eventsConn, WireEventSink)

	return wireControlConnection, wireEventConnection, nil
}
