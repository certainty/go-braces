package bytecode_test

import (
	"github.com/certainty/go-braces/vm/internal/bytecode"
	"github.com/certainty/go-braces/vm/internal/language/value"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChunkBuilderCanBuildChunkWithSingleConstant(t *testing.T) {
	valueFactory := value.NewFactory()
	chunkBuilder := bytecode.NewChunkBuilder()
	chunkBuilder.AddConstant(valueFactory.Bool(true))

	chunk := chunkBuilder.Build()

	assert.Equal(t, 1, len(chunk.Constants), "Chunk should have one constant")
}

func TestChunkBuilderCanBuildChunkWithMultipleConstants(t *testing.T) {
	valueFactory := value.NewFactory()
	chunkBuilder := bytecode.NewChunkBuilder()
	chunkBuilder.AddConstant(valueFactory.Bool(true))
	chunkBuilder.AddConstant(valueFactory.Bool(false))

	chunk := chunkBuilder.Build()

	assert.Equal(t, 2, len(chunk.Constants), "Chunk should have three constants")
}

func TestChunkBuilderDeduplicatesConstants(t *testing.T) {
	valueFactory := value.NewFactory()
	chunkBuilder := bytecode.NewChunkBuilder()
	chunkBuilder.AddConstant(valueFactory.Bool(true))
	chunkBuilder.AddConstant(valueFactory.Bool(true))

	chunk := chunkBuilder.Build()

	assert.Equal(t, 1, len(chunk.Constants), "Chunk should have one constant")
}

func TestChunkBuilderCanBuildChunkWithSingleInstruction(t *testing.T) {
	chunkBuilder := bytecode.NewChunkBuilder()
	address := chunkBuilder.AddInstruction(bytecode.Unspecified(), bytecode.NewSourceInformation(10, 3))

	chunk := chunkBuilder.Build()

	sourceInfo, err := chunk.IntrospectionInfo.SourceInformationAt(address)

	assert.Nil(t, err, "Chunk should have an instruction at the given address")
	assert.Equal(t, uint32(10), sourceInfo.Span.Line, "Chunk should return correct span for instruction")
	assert.Equal(t, uint32(3), sourceInfo.Span.Column, "Chunk should return correct span for instruction")
}

func TestChunkBuilderCanBuildChunkWithMultipleInstructions(t *testing.T) {
	chunkBuilder := bytecode.NewChunkBuilder()
	address1 := chunkBuilder.AddInstruction(bytecode.Unspecified(), bytecode.NewSourceInformation(10, 3))
	address2 := chunkBuilder.AddInstruction(bytecode.Nil(), bytecode.NewSourceInformation(10, 3))

	chunk := chunkBuilder.Build()

	sourceInfo, err := chunk.IntrospectionInfo.SourceInformationAt(address1)

	assert.Nil(t, err, "Chunk should have an instruction at the given address")
	assert.Equal(t, uint32(10), sourceInfo.Span.Line, "Chunk should return correct span for instruction")
	assert.Equal(t, uint32(3), sourceInfo.Span.Column, "Chunk should return correct span for instruction")

	sourceInfo, err = chunk.IntrospectionInfo.SourceInformationAt(address2)
	assert.Nil(t, err, "Chunk should have an instruction at the given address")
	assert.Equal(t, uint32(10), sourceInfo.Span.Line, "Chunk should return correct span for instruction")
	assert.Equal(t, uint32(3), sourceInfo.Span.Column, "Chunk should return correct span for instruction")
}
