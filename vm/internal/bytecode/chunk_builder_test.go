package bytecode_test

import (
	"testing"

	"github.com/certainty/go-braces/vm/internal/bytecode"
	"github.com/certainty/go-braces/vm/internal/language/value"
	"github.com/stretchr/testify/assert"
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
