package sly

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompareOrdered(t *testing.T) {
	assert.Equal(t, CompareOrdered(2, 3), -1)
	assert.Equal(t, CompareOrdered(3, 2), 1)
	assert.Equal(t, CompareOrdered(2, 2), 0)
}

func TestCompareReverse(t *testing.T) {
	rev := CompareReverse(CompareOrdered[int])
	assert.Equal(t, rev(2, 3), 1)
	assert.Equal(t, rev(3, 2), -1)
	assert.Equal(t, rev(2, 2), 0)
}

func TestCompareMethods(t *testing.T) {
	compare := Compare[int](CompareOrdered[int])
	assert.True(t, compare.Less(2, 3))
	assert.False(t, compare.Less(3, 2))
	assert.True(t, compare.Greater(3, 2))
	assert.False(t, compare.Greater(2, 3))
	assert.True(t, compare.Equal(2, 2))
	assert.False(t, compare.Equal(2, 3))
	assert.True(t, compare.NotEqual(2, 3))
	assert.False(t, compare.NotEqual(2, 2))
	assert.True(t, compare.LessOrEqual(2, 3))
	assert.True(t, compare.LessOrEqual(2, 2))
	assert.False(t, compare.LessOrEqual(3, 2))
	assert.True(t, compare.GreaterOrEqual(3, 2))
	assert.True(t, compare.GreaterOrEqual(2, 2))
	assert.False(t, compare.GreaterOrEqual(2, 3))
}
