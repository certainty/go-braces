package compiler

type CompilerOptions struct {
	// enable debug mode which intercepts, phases, lexing, parsing etc.
}

type Compiler struct {
	options CompilerOptions
}

type CompilationUnit struct {
}

func DefaultOptions() CompilerOptions {
	return CompilerOptions{}
}

func NewCompiler(options CompilerOptions) Compiler {
	return Compiler{options: options}
}

func (c Compiler) JitCompile(code string) (*CompilationUnit, error) {
	compilationUnit := CompilationUnit{}

	return &compilationUnit, nil
}
