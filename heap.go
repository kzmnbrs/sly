package sly

// HeapPush pushes the value onto the heap.
//
//	heapPtr: Heap pointer.
//	x: Value to push.
//	compare: Comparator function.
//
// Returns the heap with the value pushed onto it.
func HeapPush[T any](heapPtr *[]T, x T, compare Compare[T]) []T {
	heap := *heapPtr

	heap = append(heap, x)
	heapSiftUp(heap, compare, len(heap)-1)

	*heapPtr = heap
	return heap
}

// HeapPop pops the max value from the heap.
//
//	heapPtr: Heap pointer.
//	compare: Comparator function.
//
// Returns the max value.
func HeapPop[T any](heapPtr *[]T, compare Compare[T]) T {
	heap := *heapPtr

	max := heap[0]
	SliceSwap(heap, 0, len(heap)-1)
	heap = heap[:len(heap)-1]
	heapSiftDown(heap, compare, 0)

	*heapPtr = heap
	return max
}

// HeapPopTail pops the smallest value from the heap and places it on top.
//
//	heapPtr: Heap pointer.
//	compare: Comparator function.
//
// Returns the smallest value.
func HeapPopTail[T any](heapPtr *[]T, compare Compare[T]) T {
	heap := *heapPtr

	min := heap[len(heap)-1]
	SliceSwap(heap, 0, len(heap)-1)
	heap = heap[:len(heap)-1]
	heapSiftDown(heap, compare, 0)

	*heapPtr = heap
	return min
}

// Heapify makes its input a heap.
//
//	input: Slice to heapify.
//	compare: Comparator function.
func Heapify[T any](input []T, compare Compare[T]) {
	// Any index beyond len(input)/2 will be a leaf node.
	for i := len(input)>>1 + 1; i >= 1; {
		i--
		heapSiftDown(input, compare, i)
	}
}

// SortHeap is an in-place heap sort.
//
//	input: Slice to sort.
//	compare: Comparator function.
func SortHeap[T any](input []T, compare Compare[T]) {
	Heapify(input, compare)
	// Copying slice to preserve `input` length.
	heapCpy := input
	for len(heapCpy) > 1 {
		_ = HeapPopTail(&heapCpy, compare)
	}
}

func heapSiftUp[T any](heap []T, compare Compare[T], current int) {
	for {
		if current == 0 {
			break
		}

		parent := (current - 1) >> 1
		// Stop if `current` is on its place.
		if compare.LessOrEqual(heap[current], heap[parent]) {
			break
		}

		SliceSwap(heap, current, parent)
		current = parent
	}
}

func heapSiftDown[T any](heap []T, compare Compare[T], current int) {
	n := len(heap)
	for {
		n--
		// Naively pick the left one.
		maxChild := current<<1 + 1
		// Stop if `current` is a leaf.
		if maxChild > n {
			return
		}

		// If `current` has a right child
		if maxChild < n &&
			// pick the one with the highest priority.
			compare.Greater(heap[maxChild+1], heap[maxChild]) {
			maxChild++
		}

		// Stop is `current` is on its place.
		if compare.GreaterOrEqual(heap[current], heap[maxChild]) {
			return
		}

		SliceSwap(heap, current, maxChild)
		current = maxChild
	}
}
