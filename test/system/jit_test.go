package system_test

import (
	"github.com/certainty/go-braces/internal/compiler"
	"github.com/certainty/go-braces/internal/vm"
	"github.com/stretchr/testify/assert"
	"testing"
)

// runJitTest compiles the given source code, executes it on the VM, and returns the result
func runJitTest(sourceCode string) (interface{}, error) {
	compiler := compiler.NewCompiler(compiler.DefaultOptions())
	compilationUnit, err := compiler.JitCompile(sourceCode)
	if err != nil {
		return nil, err
	}

	// Create a new VM instance and execute the compilation unit
	virtualMachine := vm.NewVM(vm.DefaultOptions())
	result, err := virtualMachine.Execute(compilationUnit)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func assertCompiles(t *testing.T, sourceCode string, expectedValue interface{}) {
	t.Helper()

	result, err := runJitTest(sourceCode)

	assert.NoError(t, err)
	assert.Equal(t, expectedValue, result)
}

func assertCompilerError(t *testing.T, sourceCode string) {
	t.Helper()

	_, err := runJitTest(sourceCode)
	assert.Contains(t, err.Error(), "Compile error")
}

func assertRuntimeError(t *testing.T, sourceCode string) {
	t.Helper()

	_, err := runJitTest(sourceCode)
	assert.Contains(t, err.Error(), "Runtime error")
}
