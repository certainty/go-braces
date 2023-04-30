package location

import "fmt"

type Origin interface {
	Name() string
	Description() string
}

type FileOrigin struct {
	Path string
}

type StringOrigin struct {
	Label string
}

type ReplOrigin struct {
	InputCount uint64
}

func (f FileOrigin) Name() string {
	return f.Path
}

func (f FileOrigin) Description() string {
	return fmt.Sprintf("file://%s", f.Path)
}

func (s StringOrigin) Name() string {
	return s.Label
}

func (s StringOrigin) Description() string {
	return fmt.Sprintf("string://%s", s.Label)
}

func (r ReplOrigin) Name() string {
	return fmt.Sprintf("repl-%d", r.InputCount)
}

func (r ReplOrigin) Description() string {
	return fmt.Sprintf("repl://%d", r.InputCount)
}

type Location struct {
	Origin *Origin
	// Line is the line number of the source code that this location
	Line uint64
	// Offset is the byte offset of the source code that this location
	StartOffset uint64
	EndOffset   uint64
}
