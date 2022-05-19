package sly

import (
	"reflect"
	"unsafe"
)

// SliceSwap swaps two slice items by their indices.
//
//	values: Slice to swap the items in.
//	i: First item index.
//	j: Second item index.
func SliceSwap[T any](values []T, i, j int) {
	values[i], values[j] = values[j], values[i]
}

// SlicePop pops the last value from the slice.
//
//	ptr: Slice to pop from.
//
// Returns the popped value.
func SlicePop[T any](ptr *[]T) T {
	stack := *ptr

	val := stack[len(stack)-1]
	stack = stack[:len(stack)-1]

	*ptr = stack
	return val
}

// SliceReshape reshapes the slice to the new length.
//
//	input: Slice to reshape.
//	newLen: New length.
//
// Returns the reshaped slice.
func SliceReshape[T any](input []T, newLen int) []T {
	if newLen > cap(input) {
		input = append(make([]T, 0, newLen), input...)
	}

	tsh := *(*reflect.SliceHeader)(unsafe.Pointer(&input))
	input = *(*[]T)(unsafe.Pointer(&reflect.SliceHeader{
		Data: tsh.Data,
		Len:  newLen,
		Cap:  tsh.Cap,
	}))
	return input
}
