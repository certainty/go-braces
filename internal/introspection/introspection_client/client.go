package introspection_client

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"path/filepath"

	"github.com/certainty/go-braces/internal/introspection/introspection_protocol"
)

type IntrospectionClient struct {
	EventChan    chan introspection_protocol.Event
	RequestChan  chan introspection_protocol.Request
	ResponseChan chan introspection_protocol.Response
	eventsSock   net.Conn
	controlSock  net.Conn
}

func NewClient(ipcDir string) (*IntrospectionClient, error) {
	eventsSock, err := net.Dial("unix", filepath.Join(ipcDir, "eventstream.ipc"))
	if err != nil {
		return nil, fmt.Errorf("failed to dial eventstream socket: %w", err)
	}

	controlSock, err := net.Dial("unix", filepath.Join(ipcDir, "control.ipc"))
	if err != nil {
		return nil, fmt.Errorf("failed to dial control socket: %w", err)
	}

	client := &IntrospectionClient{
		EventChan:    make(chan introspection_protocol.Event),
		RequestChan:  make(chan introspection_protocol.Request),
		ResponseChan: make(chan introspection_protocol.Response),
		eventsSock:   eventsSock,
		controlSock:  controlSock,
	}

	go client.handleEventStream()
	go client.handleReqRep()

	return client, nil
}

func (c *IntrospectionClient) handleEventStream() {
	dec := gob.NewDecoder(c.eventsSock)
	for {
		var event introspection_protocol.WireEvent
		err := dec.Decode(&event)
		if err != nil {
			c.EventChan <- errors.New("error decoding event")
			continue
		}
		c.EventChan <- event.Event
	}
}

func (c *IntrospectionClient) handleReqRep() {
	enc := gob.NewEncoder(c.controlSock)
	dec := gob.NewDecoder(c.controlSock)

	for req := range c.RequestChan {
		err := enc.Encode(introspection_protocol.WireRequest{Request: req})
		if err != nil {
			c.ResponseChan <- fmt.Errorf("error encoding request: %w", err)
		}

		var resp introspection_protocol.WireResponse
		err = dec.Decode(&resp)
		if err != nil {
			c.ResponseChan <- fmt.Errorf("error decoding response: %w", err)
		}
		c.ResponseChan <- resp.Response
	}
}

func (c *IntrospectionClient) Close() {
	close(c.EventChan)
	close(c.RequestChan)
	close(c.ResponseChan)
	c.controlSock.Close()
	c.eventsSock.Close()
}
