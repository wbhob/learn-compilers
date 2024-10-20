package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	main "github.com/wbhob/learn-compilers"
)

func TestParseNumberSequenceShorthand(t *testing.T) {
	t.Run("Basic Functionality", func(t *testing.T) {
		t.Run("handles repetition shorthand", func(t *testing.T) {
			assert.Equal(t, []float64{42, 42, 42}, main.ParseNumberSequenceShorthand("42x3"))
			assert.Equal(t, []float64{0, 2, 2, 2, 10, 42, 42, 42}, main.ParseNumberSequenceShorthand("0, 2x3, 10, 42x3"))
			assert.Equal(t, []float64{0.1, 0.1, 2.3, 2.3, -4.5, -4.5}, main.ParseNumberSequenceShorthand(".1x2, 2.3x2, -4.5x2"))
		})

		t.Run("handles a group", func(t *testing.T) {
			assert.Equal(t, []float64{1, 2, 3, 1, 2, 3}, main.ParseNumberSequenceShorthand("(1, 2, 3)x2"))
		})

		t.Run("handles a nested group", func(t *testing.T) {
			assert.Equal(t, []float64{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2}, main.ParseNumberSequenceShorthand("(((1, 2)x2)x2)x2"))
		})

		t.Run("complex example", func(t *testing.T) {
			expected := []float64{1, 2, 1, 2, 42, 42, 42, 0.5, 1.5, 1.5, 8, 16, 8, 16, 8, 16, 0.5, 1.5, 1.5, 8, 16, 8, 16, 8, 16, 5.5}
			assert.Equal(t, expected, main.ParseNumberSequenceShorthand("(1, 2)x2, 42x3, (.5, 1.5x2, (8, 16)x3)x2, 5.5"))
		})
	})

	t.Run("Specific Features", func(t *testing.T) {
		t.Run("handles single number", func(t *testing.T) {
			assert.Equal(t, []float64{42}, main.ParseNumberSequenceShorthand("42"))
		})

		t.Run("handles multiple numbers", func(t *testing.T) {
			assert.Equal(t, []float64{1, 2, 3, 4, 5, 6}, main.ParseNumberSequenceShorthand("1, 2, 3, 4, 5, 6"))
		})

		t.Run("handles multiple numbers (strange whitespace)", func(t *testing.T) {
			assert.Equal(t, []float64{1, 2, 3, 4, 5, 6}, main.ParseNumberSequenceShorthand("  1,2,  3  , 4 ,5, 6 "))
		})

		t.Run("handles decimal values", func(t *testing.T) {
			assert.Equal(t, []float64{0.1, 0.23, 0.45, 6.7}, main.ParseNumberSequenceShorthand(".1, .23, 0.45, 6.7"))
		})

		t.Run("handles negative values", func(t *testing.T) {
			assert.Equal(t, []float64{-42, -0.1, -0.25, -3.33}, main.ParseNumberSequenceShorthand("-42, -.1, -0.25, -3.33"))
		})

		t.Run("handles repetition shorthand, with 0 repetitions", func(t *testing.T) {
			assert.Equal(t, []float64{}, main.ParseNumberSequenceShorthand("42x0"))
			assert.Equal(t, []float64{1, 2, 3}, main.ParseNumberSequenceShorthand("1, 2, 42x0, 3"))
		})

		t.Run("handles a group with a single value", func(t *testing.T) {
			assert.Equal(t, []float64{1, 1, 1}, main.ParseNumberSequenceShorthand("(1)x3"))
		})

		t.Run("handles mixed values with a group", func(t *testing.T) {
			assert.Equal(t, []float64{0, 1, 1, 2, 3, 4, 4, 2, 3, 4, 4}, main.ParseNumberSequenceShorthand("0, 1x2, (2, 3, 4x2)x2"))
		})
	})

	t.Run("Validation", func(t *testing.T) {
		t.Run("failure cases", func(t *testing.T) {
			invalidInputs := []string{
				"1 2",
				"(((1 2)x2)x2)x2",
				"1,,2",
				"(1, 2, 3)x, (4, 5)x2",
				"(1, 2, 3)2, (4, 5)x2",
				"(((1, 2x2)x2)x2",
				"((1, 2)x2)x2)x2",
			}

			for _, input := range invalidInputs {
				t.Run(input, func(t *testing.T) {
					assert.Panics(t, func() { main.ValidateNumberSequenceShorthand(input) })
				})
			}
		})
	})
}
