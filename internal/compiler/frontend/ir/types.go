package ir

type Tag uint8

const (
	IRTypeUint8 Tag = iota
	IRTypeUint16
	IRTypeUint32
	IRTypeUint64
	IRTypeInt8
	IRTypeInt16
	IRTypeInt32
	IRTypeInt64
	IRTypeUint1
	IRTypeVoid
)

type Type struct {
	Tag Tag
	// size in bits
	Size uint64
}

func IRTypeFromTag(tag Tag) Type {
	switch tag {
	case IRTypeUint8:
		return Type{Tag: IRTypeUint8, Size: 8}
	case IRTypeUint16:
		return Type{Tag: IRTypeUint16, Size: 16}
	case IRTypeUint32:
		return Type{Tag: IRTypeUint32, Size: 32}
	case IRTypeUint64:
		return Type{Tag: IRTypeUint64, Size: 64}
	case IRTypeInt8:
		return Type{Tag: IRTypeInt8, Size: 8}
	case IRTypeInt16:
		return Type{Tag: IRTypeInt16, Size: 16}
	case IRTypeInt32:
		return Type{Tag: IRTypeInt32, Size: 32}
	case IRTypeInt64:
		return Type{Tag: IRTypeInt64, Size: 64}
	case IRTypeUint1:
		return Type{Tag: IRTypeUint1, Size: 1}
	case IRTypeVoid:
		return Type{Tag: IRTypeVoid, Size: 0}
	}
	return Type{}
}
