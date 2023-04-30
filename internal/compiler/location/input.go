package location

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type ReadFunction func(reader io.Reader) (interface{}, error)

type Input interface {
	Type() string
	Buffer() ([]byte, error)
	Reader() io.Reader
}

type FileInput struct {
	path string
	file *os.File
}

func (f FileInput) Type() string {
	return fmt.Sprintf("file://%s", f.path)
}

func (f FileInput) Buffer() ([]byte, error) {
	return os.ReadFile(f.path)
}

func (f FileInput) Reader() io.Reader {
	return bufio.NewReader(f.file)
}

func NewFileInput(file os.File) (FileInput, error) {
	return FileInput{
		path: file.Name(),
		file: &file,
	}, nil
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

func (s StringInput) Reader() io.Reader {
	return strings.NewReader(s.value)
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

func (r ReplInput) Reader() io.Reader {
	return r.input.Reader()
}
