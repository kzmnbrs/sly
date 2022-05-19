package sly

import (
	"fmt"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"
)

type (
	// PriorityQueueOptions are used to construct a new priority queue.
	//
	//  InitialCap: Initial capacity.
	//  MaxCap: Max capacity. If 0, then unlimited.
	//  Lock: Queue lock. If nil, then Spinlock.
	//  Compare: Comparator function.
	PriorityQueueOptions[T any] struct {
		InitialCap uint
		MaxCap     uint
		Lock       sync.Locker
		Compare    Compare[T]
	}

	// The PriorityQueue is a thread-safe priority queue.
	PriorityQueue[T any] struct {
		heap         []T
		compare      Compare[T]
		maxCap       int
		done         atomic.Int32
		pushInFlight atomic.Int64
		popInFlight  atomic.Int64
		popCond      sync.Cond
	}
)

// NewPriorityQueue creates a new priority queue.
//
//	opts: See PriorityQueueOptions.
//
// Returns a pointer to the newly created priority queue, or an error if the options are invalid.
func NewPriorityQueue[T any](opts PriorityQueueOptions[T]) (*PriorityQueue[T], error) {
	if opts.Lock == nil {
		opts.Lock = new(Spinlock)
	}
	if opts.Compare == nil {
		return nil, fmt.Errorf("%w: nil comparator", ErrBadOptions)
	}

	return &PriorityQueue[T]{
		heap:    make([]T, 0, opts.InitialCap),
		compare: opts.Compare,
		popCond: *sync.NewCond(opts.Lock),
		maxCap:  int(opts.MaxCap),
	}, nil
}

// TryPush attempts to push an element onto the priority queue.
//
//	x: Element to push.
//
// Returns true if the element was pushed, or false
// if the queue is either full or closed.
func (pq *PriorityQueue[T]) TryPush(x T) bool {
	if pq.done.Load() == 1 {
		return false
	}

	pq.pushInFlight.Add(1)
	defer pq.pushInFlight.Add(-1)

	pq.popCond.L.Lock()
	if pq.maxCap > 0 && len(pq.heap)+1 > pq.maxCap {
		pq.popCond.L.Unlock()
		return false
	}

	HeapPush(&pq.heap, x, pq.compare)
	pq.popCond.L.Unlock()
	pq.popCond.Signal()

	return true
}

// Pop the highest priority element from the priority queue.
//
// Blocks indefinitely until there's either something to pop
// or the queue is closed.
//
// Returns the default value and false in case the queue is closed
// and nothing remains to pop.
func (pq *PriorityQueue[T]) Pop() (T, bool) {
	pq.popInFlight.Add(1)
	defer pq.popInFlight.Add(-1)

	pq.popCond.L.Lock()
	defer pq.popCond.L.Unlock()
	for len(pq.heap) == 0 && pq.done.Load() == 0 {
		pq.popCond.Wait()
	}

	if pq.done.Load() == 1 {
		var def T
		return def, false
	}

	x := HeapPop(&pq.heap, pq.compare)
	return x, true
}

// Close the priority queue.
//
// Returns the remaining elements in order of priority.
func (pq *PriorityQueue[T]) Close() []T {
	if !pq.done.CompareAndSwap(0, 1) {
		return pq.heap
	}

	backoff := 1
	for pq.pushInFlight.Load() > 0 {
		for i := 0; i < backoff; i++ {
			runtime.Gosched()
		}
		if backoff < 16 {
			backoff <<= 1
		}
	}

	for pq.popInFlight.Load() > 0 {
		pq.popCond.Broadcast()
		for i := 0; i < backoff; i++ {
			runtime.Gosched()
		}
		if backoff < 16 {
			backoff <<= 1
		}
	}

	slices.SortFunc(pq.heap, func(a, b T) int {
		return pq.compare(a, b) * -1
	})
	return pq.heap
}
