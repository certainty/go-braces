package compiler

import "github.com/certainty/go-braces/internal/introspection"

type CompilerOptions struct {
	introspectionAPI *introspection.API
}

type Compiler struct {
	introspectionAPI *introspection.API
}

type CompilationUnit struct {
}

func DefaultOptions() CompilerOptions {
	return CompilerOptions{}
}

func NewCompiler(options CompilerOptions) Compiler {
	if options.introspectionAPI == nil {
		return Compiler{introspectionAPI: introspection.NullAPI()}
	} else {
		return Compiler{introspectionAPI: options.introspectionAPI}
	}
}

func (c Compiler) JitCompile(code string) (*CompilationUnit, error) {
	compilationUnit := CompilationUnit{}

	return &compilationUnit, nil
}
