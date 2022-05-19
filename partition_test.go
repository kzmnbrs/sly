package sly

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func FuzzPartitionFat(f *testing.F) {
	f.Add([]byte{87, 40, 12, 20, 33, 20, 31, 11, 3, 49})
	f.Fuzz(func(t *testing.T, bs []byte) {
		if len(bs) == 0 {
			return
		}
		pivot := bs[rand.Intn(len(bs))]
		less, greater := PartitionFat(bs, pivot, CompareOrdered[byte])
		for _, v := range bs[:less] {
			assert.Less(t, v, pivot)
		}
		for _, v := range bs[greater+1:] {
			assert.Greater(t, v, pivot)
		}
		for _, v := range bs[less : greater+1] {
			assert.Equal(t, v, pivot)
		}
	})
}

func TestPartitionFat(t *testing.T) {
	less, greater := PartitionFat(nil, 0, CompareOrdered[int])
	assert.Equal(t, 0, less)
	assert.Equal(t, -1, greater)

	less, greater = PartitionFat([]int{0, 2, 1, 1, 1}, 1, CompareOrdered[int])
	assert.Equal(t, 1, less)
	assert.Equal(t, 3, greater)
}
