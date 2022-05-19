package sly

// IndexOfLinear returns the key index in the given slice or -1.
// Operates in linear lookup fashion.
func IndexOfLinear[T comparable](ts []T, key T) int {
	for i, v := range ts {
		if v == key {
			return i
		}
	}
	return -1
}

// IndexOfLinearWithComparator returns the key index in the given slice or -1.
// Operates in linear lookup fashion.
func IndexOfLinearWithComparator[T any](ts []T, key T, cmp Cmp[T]) int {
	for i, v := range ts {
		if cmp.Eq(v, key) {
			return i
		}
	}
	return -1
}

// IndexOfBinarySearch returns the key index in the given slice or -1.
// Lookup algorithm is binary search so the slice must be sorted.
func IndexOfBinarySearch[T any](ts []T, key T, cmp Cmp[T]) int {
	hi := len(ts) - 1
	if hi < 0 {
		return -1
	}
	_ = ts[hi]

	lo := 0
	for lo <= hi {
		mid := hi >> 1
		switch {
		case cmp.Lt(ts[mid], key):
			lo = mid + 1

		case cmp.Gt(ts[mid], key):
			hi = mid - 1

		default:
			return mid
		}
	}

	return -1
}
