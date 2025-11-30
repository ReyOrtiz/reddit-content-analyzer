package services

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// CosineSimilarity Tests
// ============================================================================

func TestCosineSimilarity(t *testing.T) {
	t.Run("IdenticalVectors", func(t *testing.T) {
		// Arrange
		vec := []float32{1.0, 2.0, 3.0, 4.0, 5.0}

		// Act
		result := CosineSimilarity(vec, vec)

		// Assert
		assert.InDelta(t, 1.0, result, 0.0001)
	})

	t.Run("OrthogonalVectors", func(t *testing.T) {
		// Arrange
		vec1 := []float32{1.0, 0.0, 0.0}
		vec2 := []float32{0.0, 1.0, 0.0}

		// Act
		result := CosineSimilarity(vec1, vec2)

		// Assert
		assert.InDelta(t, 0.0, result, 0.0001)
	})

	t.Run("OppositeVectors", func(t *testing.T) {
		// Arrange
		vec1 := []float32{1.0, 0.0, 0.0}
		vec2 := []float32{-1.0, 0.0, 0.0}

		// Act
		result := CosineSimilarity(vec1, vec2)

		// Assert
		assert.InDelta(t, -1.0, result, 0.0001)
	})

	t.Run("HighSimilarity", func(t *testing.T) {
		// Arrange
		vec1 := []float32{1.0, 2.0, 3.0}
		vec2 := []float32{1.1, 2.1, 3.1}

		// Act
		result := CosineSimilarity(vec1, vec2)

		// Assert
		assert.Greater(t, result, 0.99)
		assert.LessOrEqual(t, result, 1.0)
	})

	t.Run("LowSimilarity", func(t *testing.T) {
		// Arrange
		vec1 := []float32{1.0, 0.0, 0.0, 0.0}
		vec2 := []float32{0.0, 0.0, 0.0, 1.0}

		// Act
		result := CosineSimilarity(vec1, vec2)

		// Assert
		assert.InDelta(t, 0.0, result, 0.0001)
	})

	t.Run("DifferentLengths", func(t *testing.T) {
		// Arrange
		vec1 := []float32{1.0, 2.0, 3.0}
		vec2 := []float32{1.0, 2.0}

		// Act
		result := CosineSimilarity(vec1, vec2)

		// Assert
		assert.Equal(t, 0.0, result)
	})

	t.Run("ZeroVectors", func(t *testing.T) {
		// Arrange
		vec1 := []float32{0.0, 0.0, 0.0}
		vec2 := []float32{0.0, 0.0, 0.0}

		// Act
		result := CosineSimilarity(vec1, vec2)

		// Assert
		assert.Equal(t, 0.0, result)
	})

	t.Run("OneZeroVector", func(t *testing.T) {
		// Arrange
		vec1 := []float32{1.0, 2.0, 3.0}
		vec2 := []float32{0.0, 0.0, 0.0}

		// Act
		result := CosineSimilarity(vec1, vec2)

		// Assert
		assert.Equal(t, 0.0, result)
	})

	t.Run("NegativeValues", func(t *testing.T) {
		// Arrange
		vec1 := []float32{1.0, -2.0, 3.0}
		vec2 := []float32{-1.0, 2.0, -3.0}

		// Act
		result := CosineSimilarity(vec1, vec2)

		// Assert
		// Should be negative since vectors point in opposite directions
		assert.Less(t, result, 0.0)
		assert.GreaterOrEqual(t, result, -1.0)
	})

	t.Run("LargeVectors", func(t *testing.T) {
		// Arrange
		vec1 := make([]float32, 1000)
		vec2 := make([]float32, 1000)
		for i := range vec1 {
			vec1[i] = float32(i)
			vec2[i] = float32(i) * 0.5
		}

		// Act
		result := CosineSimilarity(vec1, vec2)

		// Assert
		// Should be 1.0 since vec2 is a scalar multiple of vec1
		assert.InDelta(t, 1.0, result, 0.0001)
	})

	t.Run("NormalizedVectors", func(t *testing.T) {
		// Arrange
		// Create normalized vectors (unit vectors)
		magnitude1 := math.Sqrt(1.0*1.0 + 2.0*2.0 + 3.0*3.0)
		magnitude2 := math.Sqrt(2.0*2.0 + 4.0*4.0 + 6.0*6.0)

		vec1 := []float32{float32(1.0 / magnitude1), float32(2.0 / magnitude1), float32(3.0 / magnitude1)}
		vec2 := []float32{float32(2.0 / magnitude2), float32(4.0 / magnitude2), float32(6.0 / magnitude2)}

		// Act
		result := CosineSimilarity(vec1, vec2)

		// Assert
		// Should be 1.0 since vec2 is a scalar multiple of vec1
		assert.InDelta(t, 1.0, result, 0.0001)
	})

	t.Run("SmallValues", func(t *testing.T) {
		// Arrange
		vec1 := []float32{0.001, 0.002, 0.003}
		vec2 := []float32{0.0011, 0.0021, 0.0031}

		// Act
		result := CosineSimilarity(vec1, vec2)

		// Assert
		assert.Greater(t, result, 0.99)
	})

	t.Run("SingleElement", func(t *testing.T) {
		// Arrange
		vec1 := []float32{5.0}
		vec2 := []float32{5.0}

		// Act
		result := CosineSimilarity(vec1, vec2)

		// Assert
		assert.InDelta(t, 1.0, result, 0.0001)
	})

	t.Run("SingleElementDifferent", func(t *testing.T) {
		// Arrange
		vec1 := []float32{5.0}
		vec2 := []float32{-5.0}

		// Act
		result := CosineSimilarity(vec1, vec2)

		// Assert
		assert.InDelta(t, -1.0, result, 0.0001)
	})
}
