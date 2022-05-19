package sly

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChanBroadcast(t *testing.T) {
	t.Run("add/delete no context", func(t *testing.T) {
		source := make(chan int)
		b := NewChanBroadcast(nil, source, 0)

		sink := make(chan int, 1)
		b.Add(sink)

		source <- 1
		assert.Equal(t, 1, <-sink)

		b.Delete(sink)
		// Must not panic.
		b.Delete(sink)

		// Giving the broadcast time to handle delete.
		time.Sleep(10 * time.Millisecond)

		source <- 2
		value, more := <-sink
		assert.Equal(t, 0, value)
		assert.False(t, more)
	})

	t.Run("add/delete cancel", func(t *testing.T) {
		b := ChanBroadcast[int]{
			add: make(chan chanBroadcastSub[int]),
		}

		ctx, cancel := context.WithCancel(context.TODO())
		go func() {
			<-time.After(time.Millisecond)
			cancel()
		}()
		// Broadcast hasn't started, add/delete must not block.
		b.AddContext(ctx, make(chan int))
		b.DeleteContext(ctx, make(chan int))
	})

	t.Run("add/delete broadcast done", func(t *testing.T) {
		source := make(chan int)
		b := NewChanBroadcast(nil, source, 0)
		close(source)

		b.Wait()
		assert.Nil(t, b.WaitContext(nil))

		// Broadcast has finished, add/delete must not block.
		b.AddContext(nil, make(chan int))
		b.DeleteContext(nil, make(chan<- int))
	})

	t.Run("sink cancel while broadcasting", func(t *testing.T) {
		source := make(chan int)
		b := NewChanBroadcast(nil, source, 0)

		sink := make(chan int, 10)
		ctx, cancel := context.WithCancel(context.TODO())

		b.AddContext(ctx, sink)
		cancel()

		done := make(chan struct{})
		go func() {
			for {
				select {
				case source <- 1:
				case <-done:
					return
				}
			}
		}()
		go func() {
			for {
				_, more := <-sink
				if !more {
					close(done)
					return
				}
			}
		}()
		<-done
	})

	t.Run("sink overflow", func(t *testing.T) {
		source := make(chan int)
		b := NewChanBroadcast(nil, source, 0)

		sink := make(chan int, 1000)
		b.Add(sink)

		for i := 0; i < 2000; i++ {
			source <- 1
		}

		for i := 0; i < 1000; i++ {
			value := <-sink
			assert.Equal(t, 1, value)
		}
		_, more := <-sink
		assert.False(t, more)
	})

	t.Run("broadcast cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		source := make(chan int)

		b := NewChanBroadcast(ctx, source, 0)
		sink := make(chan int)
		b.Add(sink)

		source <- 1
		assert.Equal(t, 1, <-sink)
		cancel()

		// Giving the broadcast time to handle cancellation.
		time.Sleep(10 * time.Millisecond)

		_, more := <-sink
		assert.False(t, more)
	})

	t.Run("wait cancel", func(t *testing.T) {
		source := make(chan int)
		b := NewChanBroadcast(nil, source, 0)

		ctx, cancel := context.WithCancel(context.TODO())
		cancel()

		assert.Equal(t, context.Canceled, b.WaitContext(ctx))
	})
}
