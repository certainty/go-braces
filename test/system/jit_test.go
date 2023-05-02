package system_test

import (
	"errors"
	"testing"

	"github.com/certainty/go-braces/internal/compiler"
	"github.com/certainty/go-braces/internal/vm"
	"github.com/stretchr/testify/assert"
)

// runJitTest compiles the given source code, executes it on the VM, and returns the result
func runJitTest(sourceCode string) (interface{}, error) {
	compiler := compiler.NewCompiler(compiler.DefaultOptions())
	assemblyModule, err := compiler.CompileString(sourceCode)
	if err != nil {
		return nil, err
	}

	// Create a new VM instance and execute the compilation unit
	virtualMachine := vm.NewVM(vm.DefaultOptions())
	result, err := virtualMachine.ExecuteModule(assemblyModule)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func assertCompilesAndRuns(t *testing.T, sourceCode string, expectedValue interface{}) {
	t.Helper()

	result, err := runJitTest(sourceCode)

	assert.NoError(t, err)
	assert.Equal(t, expectedValue, result)
}

func assertCompilationError(t *testing.T, sourceCode string) {
	t.Helper()

	_, err := runJitTest(sourceCode)
	assert.Error(t, err)
	if err != nil {
		var concreteError *compiler.CompilerError
		assert.True(t, errors.As(err, &concreteError), "Expected Compiler Error")
	}
}

func assertRuntimeError(t *testing.T, sourceCode string) {
	t.Helper()

	_, err := runJitTest(sourceCode)
	assert.Error(t, err)
	if err != nil {
		var concreteError *vm.VmError
		assert.True(t, errors.As(err, &concreteError), "Expected Runtime Error")
	}
}

func TestJitCanCompileAndExecuteSimpleProgram(t *testing.T) {
	assertCompilesAndRuns(t, "#t", true)
}
