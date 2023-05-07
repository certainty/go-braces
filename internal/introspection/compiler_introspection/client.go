package compiler_introspection

import (
	"errors"
	"github.com/certainty/go-braces/internal/introspection"
)

type Client struct {
	*introspection.Client
}

func NewClient(ipcDir string) (*Client, error) {
	client, err := introspection.NewClient(ipcDir)
	if err != nil {
		return nil, err
	}

	return &Client{
		client,
	}, nil
}

type Capabilities struct{}

func (c *Client) Helo() (*introspection.HeloResponse, error) {
	request := introspection.HeloRequest{
		IntrospectionType: introspection.IntrospectionTypeCompiler,
	}

	c.RequestChan <- request
	response := <-c.ResponseChan

	heloResponse, ok := response.(introspection.HeloResponse)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return &heloResponse, nil
}
