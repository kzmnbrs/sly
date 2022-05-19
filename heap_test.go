package sly

import (
	"sort"
	"testing"
)

func FuzzHeapPushPop(f *testing.F) {
	f.Add([]byte{87, 40, 12, 20, 33, 20, 31, 11, 3, 49})
	f.Add([]byte(nil))
	f.Fuzz(func(t *testing.T, bs []byte) {
		h := make([]int, 0, len(bs))
		for _, val := range bs {
			HeapPush(&h, int(val), CompareOrdered[int])
		}

		hs := make([]int, 0, len(bs))
		for len(h) > 0 {
			hs = append(hs, HeapPop(&h, CompareOrdered[int]))
		}
		h = SliceReshape(h, len(bs))

		sort.Slice(h, func(i, j int) bool {
			return h[i] > h[j]
		})

		for i := range bs {
			if h[i] != hs[i] {
				t.Fatalf("want: %v at %v, have: %v", h[i], i, hs[i])
			}
		}
	})
}

func FuzzSortHeap(f *testing.F) {
	f.Add([]byte{100, 48, 3, 82, 56, 62, 71, 39, 42, 22})
	f.Add([]byte(nil))
	f.Fuzz(func(t *testing.T, bs []byte) {
		h := make([]int, 0, len(bs))
		for _, val := range bs {
			h = append(h, int(val))
		}

		hs := append([]int(nil), h...)

		sort.Slice(h, func(i, j int) bool {
			return h[i] < h[j]
		})

		SortHeap(hs, CompareOrdered[int])

		for i := range bs {
			if h[i] != hs[i] {
				t.Fatalf("want: %v at %v, have: %v", h[i], i, hs[i])
			}
		}
	})
}
