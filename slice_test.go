package sly

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceSwap(t *testing.T) {
	s := []int{1, 2, 3}
	SliceSwap(s, 0, 2)
	assert.Equal(t, s, []int{3, 2, 1})
}

func TestSliceReshape(t *testing.T) {
	t.Run("reshape to larger size", func(t *testing.T) {
		ts := []int{1, 2, 3, 4, 5}
		newLen := 10
		expected := []int{1, 2, 3, 4, 5, 0, 0, 0, 0, 0}
		res := SliceReshape(ts, newLen)
		assert.Equal(t, expected, res)
	})

	t.Run("reshape to smaller size", func(t *testing.T) {
		ts := []int{1, 2, 3, 4, 5}
		newLen := 3
		expected := []int{1, 2, 3}
		res := SliceReshape(ts, newLen)
		assert.Equal(t, expected, res)
	})

	t.Run("reshape to same size", func(t *testing.T) {
		ts := []int{1, 2, 3, 4, 5}
		newLen := 5
		expected := []int{1, 2, 3, 4, 5}
		res := SliceReshape(ts, newLen)
		assert.Equal(t, expected, res)
	})

	t.Run("reshape to empty slice", func(t *testing.T) {
		ts := []int{1, 2, 3, 4, 5}
		newLen := 0
		expected := []int{}
		res := SliceReshape(ts, newLen)
		assert.Equal(t, expected, res)
	})
}

func TestSlicePop(t *testing.T) {
	t.Run("pop from non-empty slice", func(t *testing.T) {
		stack := []int{1, 2, 3, 4, 5}
		expectedPopValue := 5
		expectedStack := []int{1, 2, 3, 4}

		result := SlicePop(&stack)

		assert.Equal(t, expectedPopValue, result)
		assert.Equal(t, expectedStack, stack)
	})

	t.Run("pop from single element slice", func(t *testing.T) {
		stack := []int{1}
		expectedPopValue := 1
		expectedStack := []int{}

		result := SlicePop(&stack)

		assert.Equal(t, expectedPopValue, result)
		assert.Equal(t, expectedStack, stack)
	})
}
