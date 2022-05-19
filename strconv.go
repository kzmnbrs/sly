package sly

import (
	"reflect"
	"unsafe"
)

// S2B converts the string to a byte slice with no allocations.
//
//	s: String to convert.
//
// Returns the byte slice representation of the string.
func S2B(s string) []byte {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}))
}

// B2S converts the byte slice to string with no allocations.
//
//	b: Byte slice to convert.
//
// Returns the string representation of the byte slice.
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
