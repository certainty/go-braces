package compiler

import "github.com/certainty/go-braces/internal/introspection"

type CompilerOptions struct {
	introspectionAPI introspection.API
}

func DefaultOptions() CompilerOptions {
	api := introspection.NullAPI()

	return CompilerOptions{
		introspectionAPI: api,
	}
}
