package sly

// HeapPush pushes the value onto the heap.
func HeapPush[T any](heapPtr *[]T, val T, cmp Cmp[T]) []T {
	heap := *heapPtr

	heap = append(heap, val)
	heapSiftUp(heap, cmp, len(heap)-1)

	*heapPtr = heap
	return heap
}

// HeapPop pops the max value from the heap.
func HeapPop[T any](heapPtr *[]T, cmp Cmp[T]) T {
	heap := *heapPtr

	max := heap[0]
	Swap(heap, 0, len(heap)-1)
	heap = heap[:len(heap)-1]
	heapSiftDown(heap, cmp, 0)

	*heapPtr = heap
	return max
}

// HeapBuild makes its input a heap.
func HeapBuild[T any](heap []T, cmp Cmp[T]) {
	for i := len(heap)>>1 + 1; i >= 1; {
		i--
		heapSiftDown(heap, cmp, i)
	}
}

// HeapSort is an in-place heap sort.
func HeapSort[T any](heap []T, cmp Cmp[T]) {
	HeapBuild(heap, cmp)
	// Copying slice to preserve `heap` length.
	heapCpy := heap
	for len(heapCpy) > 0 {
		_ = HeapPop(&heapCpy, cmp)
	}
}

func heapSiftUp[T any](heap []T, cmp Cmp[T], current int) {
	for {
		if current == 0 {
			break
		}

		parent := (current - 1) >> 1
		// Stop if `current` is on its place.
		if cmp.Leq(heap[current], heap[parent]) {
			break
		}

		Swap(heap, current, parent)
		current = parent
	}
}

func heapSiftDown[T any](heap []T, cmp Cmp[T], current int) {
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
			cmp.Gt(heap[maxChild+1], heap[maxChild]) {
			maxChild++
		}

		// Stop is `current` is on its place.
		if cmp.Geq(heap[current], heap[maxChild]) {
			return
		}

		Swap(heap, current, maxChild)
		current = maxChild
	}
}
