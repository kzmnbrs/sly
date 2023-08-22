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
	//  Limit: Max capacity. If 0, then unlimited.
	//  Lock: Queue lock. If nil, then SpinLock.
	//  Compare: Comparator function.
	PriorityQueueOptions[T any] struct {
		InitialCap uint
		Limit      uint
		Lock       sync.Locker
		Compare    Compare[T]
	}

	// The PriorityQueue is a thread-safe priority queue.
	PriorityQueue[T any] struct {
		heap        []T
		compare     Compare[T]
		lim         int
		done        atomic.Int32
		pushPending atomic.Int64
		popPending  atomic.Int64
		popCond     sync.Cond
	}
)

// NewPriorityQueue creates a new priority queue.
//
//	opts: See PriorityQueueOptions.
//
// Returns a pointer to the newly created priority queue, or
// an error if the options are invalid.
func NewPriorityQueue[T any](opts PriorityQueueOptions[T]) (*PriorityQueue[T], error) {
	if opts.Lock == nil {
		opts.Lock = new(SpinLock)
	}
	if opts.Compare == nil {
		return nil, fmt.Errorf("%w: nil comparator", ErrBadOptions)
	}

	if opts.Limit > opts.InitialCap {
		opts.InitialCap = opts.Limit
	}
	return &PriorityQueue[T]{
		heap:    make([]T, 0, opts.InitialCap),
		compare: opts.Compare,
		popCond: *sync.NewCond(opts.Lock),
		lim:     int(opts.Limit),
	}, nil
}

// TryPush attempts to push an element onto the priority queue.
//
//	x: Element to push.
//
// Returns true if the element was pushed, or false
// if the queue is either full or closed.
func (pq *PriorityQueue[T]) TryPush(x T) bool {
	if pq.IsClosed() {
		return false
	}
	pq.pushPending.Add(1)
	defer pq.pushPending.Add(-1)

	pq.popCond.L.Lock()
	if pq.lim > 0 && len(pq.heap)+1 > pq.lim {
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
// or nothing remains to pop.
func (pq *PriorityQueue[T]) Pop() (T, bool) {
	if pq.IsClosed() {
		var z T
		return z, false
	}
	pq.popPending.Add(1)
	defer pq.popPending.Add(-1)

	pq.popCond.L.Lock()
	defer pq.popCond.L.Unlock()
	for len(pq.heap) == 0 && !pq.IsClosed() {
		pq.popCond.Wait()
	}

	// If cond has unlocked due to the queue closure.
	if pq.IsClosed() {
		var z T
		return z, false
	}

	x := HeapPop(&pq.heap, pq.compare)
	return x, true
}

// Close the priority queue. Denies all further mutations.
// Blocks until all pending operations are done.
//
// Returns the remaining elements in order of priority.
func (pq *PriorityQueue[T]) Close() []T {
	if !pq.done.CompareAndSwap(0, 1) {
		return pq.heap
	}

	pq.waitZero(&pq.pushPending, nil)
	pq.waitZero(&pq.popPending, func() {
		pq.popCond.Broadcast()
	})

	slices.SortFunc(pq.heap, func(a, b T) int {
		return pq.compare(a, b) * -1
	})
	return pq.heap
}

// IsClosed returns whether the queue is closed.
func (pq *PriorityQueue[T]) IsClosed() bool {
	return pq.done.Load() == 1
}

func (pq *PriorityQueue[T]) waitZero(v *atomic.Int64, effect func()) {
	backoff := 1
	for v.Load() > 0 {
		if effect != nil {
			effect()
		}
		for i := 0; i < backoff; i++ {
			runtime.Gosched()
		}
		if backoff < 16 {
			backoff <<= 1
		}
	}
}
