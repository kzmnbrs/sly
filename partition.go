package sly

// PartitionFat executes Fat (aka Dutch flag) partition of x around the pivot.
//
//	x: Slice to be partitioned.
//	pivot: Value to partition around.
//	compare: Comparator function.
//
// Returns a pair of indices:
//
//	less: all(y < pivot for y in x[:less]).
//	greater: all (y > pivot for y in x[greater+1:]).
//
//	all(y == pivot for y in x[less:greater+1])
func PartitionFat[T any](x []T, pivot T, compare Compare[T]) (less, greater int) {
	less = 0
	equal := 0
	greater = len(x) - 1

	// < (moves right) | == (moves right) | > (moves left).
	for equal <= greater {
		cmp := compare(x[equal], pivot)
		switch {
		case cmp < 0:
			SliceSwap(x, less, equal)
			less++
			equal++

		case cmp == 0:
			equal++

		default:
			SliceSwap(x, equal, greater)
			greater--
		}
	}
	return
}
