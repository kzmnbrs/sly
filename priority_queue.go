package sly

import (
	"context"
	"fmt"
	"sync"
)

type (
	// PriorityQueueOptions are used to construct a new priority queue.
	//
	//  InitialCap: Initial capacity.
	//  Limit: Max capacity. If 0, then unlimited.
	//  Locker: Queue lock. If nil, then SpinLock.
	//  Compare: Comparator function.
	PriorityQueueOptions[T any] struct {
		Limit   uint
		Locker  sync.Locker
		Compare Compare[T]
	}

	// The PriorityQueue is a thread-safe priority queue.
	PriorityQueue[T any] struct {
		heap    []T
		compare Compare[T]
		lim     int
		wait    chan struct{}
		locker  sync.Locker
	}
)

// NewPriorityQueue creates a new priority queue.
//
//	opts: See PriorityQueueOptions.
//
// Returns a pointer to the newly created priority queue, or
// an error if the options are invalid.
func NewPriorityQueue[T any](opts PriorityQueueOptions[T]) (*PriorityQueue[T], error) {
	if opts.Locker == nil {
		opts.Locker = new(SpinLock)
	}
	if opts.Compare == nil {
		return nil, fmt.Errorf("%w: nil comparator", ErrBadOptions)
	}
	if opts.Limit == 0 {
		return nil, fmt.Errorf("%w: unlimitied pq's are not supported", ErrBadOptions)
	}

	return &PriorityQueue[T]{
		heap:    make([]T, 0, opts.Limit),
		compare: opts.Compare,
		locker:  opts.Locker,
		wait:    make(chan struct{}, opts.Limit),
		lim:     int(opts.Limit),
	}, nil
}

// TryPush attempts to push an element onto the priority queue.
//
//	x: Element to push.
//
// Returns true if the element has been pushed or false is the queue is full.
func (pq *PriorityQueue[T]) TryPush(x T) bool {
	pq.locker.Lock()
	if len(pq.heap)+1 > pq.lim {
		pq.locker.Unlock()
		return false
	}

	HeapPush(&pq.heap, x, pq.compare)
	pq.wait <- struct{}{}
	pq.locker.Unlock()
	return true
}

// Pop the highest priority element from the priority queue.
//
//	ctx: Cancellation context. If nil, defaults to context.Background().
//
// Blocks indefinitely until there's either something to pop
// or the context is done.
//
// Returns the default value and false in case the context is done.
func (pq *PriorityQueue[T]) Pop(ctx context.Context) (T, bool) {
	if ctx == nil {
		ctx = context.Background()
	}

	pq.locker.Lock()
	if len(pq.heap) != 0 {
		x := HeapPop(&pq.heap, pq.compare)
		pq.locker.Unlock()
		return x, true
	}
	pq.locker.Unlock()

	var z T
	select {
	case <-pq.wait:
		pq.locker.Lock()
		// Another goroutine could've popped the heap already.
		if len(pq.heap) == 0 {
			pq.locker.Unlock()
			return z, false
		}
		x := HeapPop(&pq.heap, pq.compare)
		pq.locker.Unlock()
		return x, true
	case <-ctx.Done():
		return z, false
	}
}
