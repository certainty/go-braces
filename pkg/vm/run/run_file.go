package run

import (
	"fmt"

	"github.com/certainty/go-braces/pkg/compiler"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/lexer"
	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
	"github.com/certainty/go-braces/pkg/shared/isa"
	"github.com/certainty/go-braces/pkg/vm"
	log "github.com/sirupsen/logrus"
)

type Runner struct {
	vmInstance            *vm.VM
	compiler              *compiler.Compiler
	compilerIntrospection *compiler_introspection.Server
	options               Options
}

type Options struct {
	IntrospectCompiler bool
	IntrospectVM       bool
}

func NewRunner(options Options) (*Runner, error) {
	instance := vm.NewVM(vm.DefaultOptions())

	return &Runner{
		vmInstance: instance,
		options:    options,
	}, nil
}

func (r *Runner) Run(path string) {
	if r.options.IntrospectCompiler {
		println("Compiler Introspection is enabled")
		println("Waiting for introspection client ....")

		var err error
		r.compilerIntrospection, err = compiler_introspection.NewServer()
		if err != nil {
			return
		}

		instrumenter, err := r.compilerIntrospection.WaitForClient()
		if err != nil {
			return
		}

		r.compiler = compiler.NewCompiler(compiler.NewCompilerOptions(instrumenter))
		r.doRun(path)
	} else {
		r.compiler = compiler.NewCompiler(compiler.DefaultOptions())
		r.doRun(path)
	}
}

func (r *Runner) doRun(path string) {
	input, err := lexer.NewFileInput(path)

	if err != nil {
		println(err.Error())
		return
	}

	log.Infof("Compiling %s", path)
	result, err := r.compileAndRun(input)
	if err != nil {
		println(err.Error())
	} else {
		fmt.Printf("%s\n", r.vmInstance.WriteValue(result))
	}
}

func (r *Runner) compileAndRun(input *lexer.Input) (isa.Value, error) {
	assemblyModule, err := r.compiler.CompileModule(input)

	if err != nil {
		return nil, err
	}
	return r.vmInstance.ExecuteModule(assemblyModule)
}
