package sly

import (
	"reflect"
	"unsafe"
)

// Reshape reshapes the slice to the new length.
func Reshape[T any](ts []T, newLen int) []T {
	if newLen > cap(ts) {
		ts = append(make([]T, 0, newLen), ts...)
	}

	tsh := *(*reflect.SliceHeader)(unsafe.Pointer(&ts))
	ts = *(*[]T)(unsafe.Pointer(&reflect.SliceHeader{
		Data: tsh.Data,
		Len:  newLen,
		Cap:  tsh.Cap,
	}))
	return ts
}
