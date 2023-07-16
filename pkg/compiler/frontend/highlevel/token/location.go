package token

import "fmt"

// Origin is an interface that represents the origin of the source code.
// It could be a file, a string, or a repl.
// Each origin has a Name and a Description.
type (
	Origin interface {
		// Name returns the name of the origin.
		Name() string
		// Description returns a string representation of the origin.
		Description() string
	}

	// FileOrigin represents an origin that is a file.
	// It stores the path of the file.
	FileOrigin struct {
		// The path of the file.
		Path string
	}

	// StringOrigin represents an origin that is a string.
	// It stores the label of the string.
	StringOrigin struct {
		// A symbolic name for where the string came from
		Label string
	}

	// ReplOrigin represents an origin that is a repl.
	// It stores the input count of the repl.
	ReplOrigin struct {
		InputCount uint64
	}

	// Location represents a specific point or range in the source code.
	// It includes the origin of the source code and positional information
	// like line number, column number, start offset, and end offset.
	Location struct {
		// The origin of the source code
		Origin Origin

		// The line number of the source code that this location
		Line Line
		// The column number of the source code that this location
		Colunm Column

		// The byte offset of the start of the source code that this location
		StartOffset From

		// The byte offset of the end of the source code that this location
		EndOffset To
	}

	Line   uint64
	Column uint64
	From   uint64
	To     uint64
)

// NewStringOrigin creates a new StringOrigin with the given label.
func NewStringOrigin(label string) StringOrigin {
	return StringOrigin{Label: label}
}

func NewLocation(origin Origin, line Line, colunm Column, startOffset From, endOffset To) Location {
	return Location{
		Origin:      origin,
		Line:        line,
		Colunm:      colunm,
		StartOffset: startOffset,
		EndOffset:   endOffset,
	}
}

// Name returns the path of the FileOrigin
func (f FileOrigin) Name() string   { return f.Path }
func (s StringOrigin) Name() string { return s.Label }
func (r ReplOrigin) Name() string   { return fmt.Sprintf("repl-%d", r.InputCount) }

func (f FileOrigin) Description() string   { return fmt.Sprintf("file://%s", f.Path) }
func (s StringOrigin) Description() string { return fmt.Sprintf("string://%s", s.Label) }
func (r ReplOrigin) Description() string   { return fmt.Sprintf("repl://%d", r.InputCount) }

func (l Location) String() string {
	return fmt.Sprintf("%s:%d:%d", l.Origin.Description(), l.Line, l.Colunm)
}
