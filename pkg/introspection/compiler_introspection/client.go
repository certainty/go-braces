package compiler_introspection

import (
	"errors"
	"log"

	"github.com/certainty/go-braces/pkg/introspection"
)

type Client struct {
	wireClient *introspection.WireClient[CompilerIntrospectionControl, CompilerIntrospectionEvent]
	events     *introspection.WireEventConnection[CompilerIntrospectionEvent]
	control    *introspection.WireControlConnection[CompilerIntrospectionControl]
}

func NewClient() (*Client, error) {
	client := introspection.NewWireClient[CompilerIntrospectionControl, CompilerIntrospectionEvent](INTROSPECTION_TOOL_NAME)

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
		return CompilerIntrospectionEvent{}, errors.New("No client connected")
	}

	event := <-c.events.Channel
	return event, nil
}

func (c *Client) SendControl(control CompilerIntrospectionControl) error {
	if !c.IsConnected() {
		return errors.New("No client connected")
	}
	c.control.Out <- control
	return nil
}

func (c *Client) ReceiveControl() (CompilerIntrospectionControl, error) {
	if !c.IsConnected() {
		return nil, errors.New("No client connected")
	}
	control := <-c.control.In
	return control, nil
}

func (c *Client) BreakpointContinue() error {
	if !c.IsConnected() {
		return errors.New("No client connected")
	}

	log.Println("Sending breakpoint continue command")
	c.SendControl(BreakpointContinue{})
	log.Println("Waiting for response")
	response, err := c.ReceiveControl()
	log.Printf("Received response %v", response)

	if err != nil {
		return err
	}

	switch response := response.(type) {
	case CommandOk:
		return nil
	case CommandError:
		return errors.New(response.Message)
	default:
		log.Printf("Unexpected response: %v", response)
	}

	return nil
}
