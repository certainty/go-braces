package location

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type ReadFunction func(reader io.Reader) (interface{}, error)

type Input interface {
	Type() string
	Buffer() ([]byte, error)
	WithReader(fn ReadFunction) (interface{}, error)
}

type FileInput struct {
	path string
}

func (f FileInput) Type() string {
	return fmt.Sprintf("file://%s", f.path)
}

func (f FileInput) Buffer() ([]byte, error) {
	return ioutil.ReadFile(f.path)
}

func (f FileInput) WithReader(r ReadFunction) (interface{}, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	return r(reader)
}

type StringInput struct {
	label string
	value string
}

func NewStringInput(label, value string) StringInput {
	return StringInput{
		label: label,
		value: value,
	}
}

func (s StringInput) Type() string {
	return fmt.Sprintf("string://%s", s.label)
}

func (s StringInput) Buffer() ([]byte, error) {
	return []byte(s.value), nil
}

func (s StringInput) WithReader(r ReadFunction) (interface{}, error) {
	reader := strings.NewReader(s.value)
	return r(reader)
}

type ReplInput struct {
	count int
	input StringInput
}

func NewReplInput(input string, count int) ReplInput {
	return ReplInput{
		count: count,
		input: NewStringInput("REPL", input),
	}
}

func (r ReplInput) Type() string {
	return fmt.Sprintf("repl://%d", r.count)
}

func (r ReplInput) Buffer() ([]byte, error) {
	return r.input.Buffer()
}

func (r ReplInput) WithReader(fn ReadFunction) (interface{}, error) {
	return r.input.WithReader(fn)
}
