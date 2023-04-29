package location

type Location struct {
	Input *Input
	// Line is the line number of the source code that this location
	Line int
	// Column is the column number of the source code that this location
	Column int
	// Offset is the byte offset of the source code that this location
	Offset int
}
