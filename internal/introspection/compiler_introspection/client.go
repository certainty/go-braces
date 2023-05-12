package compiler_introspection

import (
	"errors"
	"github.com/certainty/go-braces/internal/introspection"
)

type Client struct {
	wireClient *introspection.WireClient
	events     *introspection.WireEventConnection
	control    *introspection.WireControlConnection
}

func NewClient() (*Client, error) {
	client := introspection.NewWireClient(INTROSPECTION_TOOL_NAME)

	return &Client{
		wireClient: &client,
	}, nil
}

func (c *Client) Close() {
	if c.control != nil {
		c.control.Close()
	}

	if c.events != nil {
		c.events.Close()
	}
}

func (c *Client) IsConnected() bool {
	return c.events != nil && c.events.IsOpen() && c.control != nil && c.control.IsOpen()
}

func (c *Client) Connect() error {
	if c.IsConnected() {
		return nil
	}

	control, events, err := c.wireClient.Connect()
	if err != nil {
		return err
	}

	c.events = events
	c.control = control

	return nil
}

func (c *Client) EnableBreakPoints() error {
	return nil
}

func (c *Client) DisableBreakPoints() error {
	return nil
}

func (c *Client) PollEvent() (CompilerIntrospectionEvent, error) {
	if !c.IsConnected() {
		return nil, errors.New("No client connected")
	}

	event := <-c.events.Channel
	switch event := event.(type) {
	case CompilerIntrospectionEvent:
		return event, nil
	default:
		return nil, errors.New("Not a CompilerIntrospectionEvent")
	}
}
