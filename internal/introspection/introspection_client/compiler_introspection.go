package introspection_client

import (
	"errors"
	"github.com/certainty/go-braces/internal/introspection/introspection_protocol"
)

type CompilerIntrospectionClient struct {
	*IntrospectionClient
}

func NewCompilerIntrospectionClient(ipcDir string) (*CompilerIntrospectionClient, error) {
	introspection_protocol.RegisterTypes()

	client, err := NewClient(ipcDir)
	if err != nil {
		return nil, err
	}

	return &CompilerIntrospectionClient{
		IntrospectionClient: client,
	}, nil
}

type Capabilities struct{}

func (c *CompilerIntrospectionClient) Helo() (*introspection_protocol.HeloResponse, error) {
	request := introspection_protocol.HeloRequest{
		IntrospectionType: introspection_protocol.IntrospectionTypeCompiler,
	}

	c.RequestChan <- request
	response := <-c.ResponseChan

	heloResponse, ok := response.(introspection_protocol.HeloResponse)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return &heloResponse, nil
}
