package location

type Location struct {
	Input *Input
	// Line is the line number of the source code that this location
	Line int
	// Offset is the byte offset of the source code that this location
	StartOffset int
	EndOffset   int
}
