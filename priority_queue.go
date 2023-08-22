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
		InitialCap uint
		Limit      uint
		Locker     sync.Locker
		Compare    Compare[T]
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

	if opts.Limit > opts.InitialCap {
		opts.InitialCap = opts.Limit
	}
	return &PriorityQueue[T]{
		heap:    make([]T, 0, opts.InitialCap),
		compare: opts.Compare,
		locker:  opts.Locker,
		wait:    make(chan struct{}),
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
	if pq.lim > 0 && len(pq.heap)+1 > pq.lim {
		pq.locker.Unlock()
		return false
	}

	HeapPush(&pq.heap, x, pq.compare)
	pq.locker.Unlock()
	select {
	case pq.wait <- struct{}{}:
	default:
	}
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
	defer pq.locker.Unlock()
	for len(pq.heap) == 0 {
		pq.locker.Unlock()
		select {
		case <-pq.wait:
			pq.locker.Lock()
		case <-ctx.Done():
			var z T
			return z, false
		}
	}

	x := HeapPop(&pq.heap, pq.compare)
	return x, true
}
