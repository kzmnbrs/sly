package sly

import (
	"context"
	"sync"
)

// ChanRelay routes data from the source channel to the sink channel.
//
//	ctx: Cancellation context. If nil, defaults to context.Background().
//	source: Channel to read from.
//	sink: Channel to write to.
func ChanRelay[T any](ctx context.Context, source <-chan T, sink chan<- T) {
	if ctx == nil {
		ctx = context.Background()
	}

	for {
		select {
		case value, more := <-source:
			if !more {
				return
			}

			select {
			case sink <- value:
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// ChanMerge takes multiple channels and merges their outputs into one.
//
//	ctx: Cancellation context. If nil, defaults to context.Background().
//	bufSize: Buffer size for the merged channel.
//	sources: Channels to merge.
//
// Returns the merged channel.
func ChanMerge[T any](ctx context.Context, bufSize uint, sources ...<-chan T) chan T {
	if ctx == nil {
		ctx = context.Background()
	}

	wg := sync.WaitGroup{}
	wg.Add(len(sources))

	sink := make(chan T, bufSize)
	for _, source := range sources {
		go func(source <-chan T) {
			defer wg.Done()
			ChanRelay(ctx, source, sink)
		}(source)
	}

	go func() {
		wg.Wait()
		close(sink)
	}()
	return sink
}

// ChanStream produces a stream of values on a channel.
//
//	ctx: Cancellation context. If nil, defaults to context.Background().
//	bufSize: Buffer size for the channel returned.
//	values: Values to stream.
//
// Returns the stream channel.
func ChanStream[T any](ctx context.Context, bufSize uint, values ...T) <-chan T {
	if ctx == nil {
		ctx = context.Background()
	}

	sink := make(chan T, bufSize)
	go func() {
		defer close(sink)
		for _, value := range values {
			select {
			case sink <- value:
			case <-ctx.Done():
				return
			}
		}
	}()
	return sink
}
