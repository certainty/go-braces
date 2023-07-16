package vm

type InternedStringTable struct {
	stringsToIndices map[string]int
	indicesToStrings []string
}

func NewInternedStringTable() *InternedStringTable {
	return &InternedStringTable{
		stringsToIndices: make(map[string]int),
		indicesToStrings: []string{},
	}
}

func (table *InternedStringTable) Intern(str string) int {
	if index, exists := table.stringsToIndices[str]; exists {
		return index
	}

	index := len(table.indicesToStrings)
	table.stringsToIndices[str] = index
	table.indicesToStrings = append(table.indicesToStrings, str)

	return index
}

func (table *InternedStringTable) GetString(index int) string {
	return table.indicesToStrings[index]
}
