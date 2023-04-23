package bytecode

import "fmt"

type SourceSpan struct {
	line   uint32
	column uint32
}

type SourceInformation struct {
	Span SourceSpan
	// we will store additional information in the future
}

type IntrospectionInfo struct {
	sourceInformation []SourceInformation
}

func (ii IntrospectionInfo) SourceInformationAt(address Address) (*SourceInformation, error) {
	if address >= uint32(len(ii.sourceInformation)) {
		return nil, fmt.Errorf("no source information for instruction at address %d", address)
	}
	return &ii.sourceInformation[address], nil
}
