package sly

import "context"

type (
	chanBroadcastSub[T any] struct {
		ctx  context.Context
		sink chan<- T
	}

	// ChanBroadcast allows for broadcasting values from the source
	// channel to multiple sink channels. Ensure that sink channels
	// have adequate capacity to keep up with the broadcasts.
	ChanBroadcast[T any] struct {
		done   chan struct{}
		accept context.Context
		source <-chan T

		add    chan chanBroadcastSub[T]
		delete chan chan<- T
		sinks  map[chan<- T]context.Context
	}
)

// NewChanBroadcast creates a new ChanBroadcast with the given source
// channel and initial capacity for sinks.
//
//	accept: Cancellation context. If nil, defaults to context.Background().
//	source: Channel to read from.
//	nsinks: Initial capacity for the sinks map.
func NewChanBroadcast[T any](accept context.Context, source <-chan T, nsinks uint) ChanBroadcast[T] {
	if accept == nil {
		accept = context.Background()
	}
	b := ChanBroadcast[T]{
		done:   make(chan struct{}),
		accept: accept,
		source: source,
		add:    make(chan chanBroadcastSub[T]),
		delete: make(chan chan<- T),
		sinks:  make(map[chan<- T]context.Context, nsinks),
	}
	go b.run()
	return b
}

// AddContext registers a new sink channel to receive broadcasts.
//
//	broadcast: Cancellation context. If nil, defaults to context.Background().
//	sink: Sink channel.
//
// If the broadcast is already canceled, the subscription is ignored.
//
// If the sink channel is already registered, the subscription is replaced.
func (b *ChanBroadcast[T]) AddContext(broadcast context.Context, sink chan<- T) {
	if broadcast == nil {
		broadcast = context.Background()
	}

	select {
	case b.add <- chanBroadcastSub[T]{
		ctx:  broadcast,
		sink: sink,
	}:
	case <-broadcast.Done():
	// Using b.done here, instead of the b.accept.Done(),
	// because of the source closure case.
	case <-b.done:
	}
}

// Add a new sink channel to receive broadcasts.
//
// This is a convenience function for AddContext with background context.
//
// See AddContext for more details.
func (b *ChanBroadcast[T]) Add(sink chan<- T) {
	b.AddContext(context.Background(), sink)
}

// DeleteContext deletes the sink from broadcasting queue.
//
//	delete: Cancellation context. If nil, defaults to context.Background().
//	sink: Sink channel previously registered via Add.
//
// If the broadcast is canceled, the deletion is ignored.
//
// If the sink channel is not registered, the deletion is ignored.
func (b *ChanBroadcast[T]) DeleteContext(delete context.Context, sink chan<- T) {
	if delete == nil {
		delete = context.Background()
	}

	select {
	case b.delete <- sink:
	case <-delete.Done():
	// Using b.done here, instead of the b.accept.Done(),
	// because of the source closure case.
	case <-b.done:
	}
}

// Delete the sink from broadcasting queue.
//
// This is a convenience function for DeleteContext with background delete context.
//
// See DeleteContext for more details.
func (b *ChanBroadcast[T]) Delete(sink chan<- T) {
	b.DeleteContext(context.Background(), sink)
}

// WaitContext blocks until the broadcaster has finished.
//
//	wait: Cancellation context. If nil, defaults to context.Background().
//
// Returns an error if the wait was canceled.
func (b *ChanBroadcast[T]) WaitContext(wait context.Context) error {
	if wait == nil {
		wait = context.Background()
	}

	select {
	case <-wait.Done():
		return wait.Err()
	case <-b.done:
		return nil
	}
}

// Wait blocks until the broadcaster has finished.
//
// This is a convenience function for WaitContext with background wait context.
//
// See WaitContext for more details.
func (b *ChanBroadcast[T]) Wait() {
	_ = b.WaitContext(context.Background())
}

func (b *ChanBroadcast[T]) run() {
	defer func() {
		for sink := range b.sinks {
			b.closeAndDeleteLF(sink)
		}
		close(b.done)
	}()

	for {
		select {
		case value, more := <-b.source:
			if !more {
				return
			}

			for sink, ctx := range b.sinks {
				select {
				case sink <- value:
				case <-ctx.Done():
					b.closeAndDeleteLF(sink)
				default:
					b.closeAndDeleteLF(sink)
				}
			}

		case sub := <-b.add:
			b.sinks[sub.sink] = sub.ctx

		case sink := <-b.delete:
			b.closeAndDeleteLF(sink)

		case <-b.accept.Done():
			return
		}
	}
}

func (b *ChanBroadcast[T]) closeAndDeleteLF(sink chan<- T) {
	_, ok := b.sinks[sink]
	if !ok {
		return
	}

	close(sink)
	delete(b.sinks, sink)
}
