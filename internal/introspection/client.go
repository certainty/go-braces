package introspection

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"path/filepath"
)

type Client struct {
	EventChan    chan IntrospectionEvent
	RequestChan  chan Request
	ResponseChan chan Response
	eventsSock   net.Conn
	controlSock  net.Conn
}

func NewClient(ipcDir string) (*Client, error) {
	eventsSock, err := net.Dial("unix", filepath.Join(ipcDir, "eventstream.ipc"))
	if err != nil {
		return nil, fmt.Errorf("failed to dial eventstream socket: %w", err)
	}

	controlSock, err := net.Dial("unix", filepath.Join(ipcDir, "control.ipc"))
	if err != nil {
		return nil, fmt.Errorf("failed to dial control socket: %w", err)
	}

	client := &Client{
		EventChan:    make(chan IntrospectionEvent),
		RequestChan:  make(chan Request),
		ResponseChan: make(chan Response),
		eventsSock:   eventsSock,
		controlSock:  controlSock,
	}

	go client.handleEventStream()
	go client.handleReqRep()

	return client, nil
}

func (c *Client) handleEventStream() error {
	dec := gob.NewDecoder(c.eventsSock)
	for {
		var event WireEvent
		err := dec.Decode(&event)
		if err != nil {
			if err == io.EOF {
				return err
			}
			log.Printf("error decoding event: %v", err)
			continue
		}
		c.EventChan <- event.IntrospectionEvent
	}
}

func (c *Client) handleReqRep() error {
	enc := gob.NewEncoder(c.controlSock)
	dec := gob.NewDecoder(c.controlSock)

	for req := range c.RequestChan {
		err := enc.Encode(WireRequest{Request: req})
		if err != nil {
			return fmt.Errorf("error encoding request: %w", err)
		}

		var resp WireResponse
		err = dec.Decode(&resp)
		if err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}
		c.ResponseChan <- resp.Response
	}
	return nil
}

func (c *Client) Close() {
	close(c.EventChan)
	close(c.RequestChan)
	close(c.ResponseChan)
	c.controlSock.Close()
	c.eventsSock.Close()
}
