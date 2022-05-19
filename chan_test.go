package sly

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChanRelay(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		source, sink := make(chan int), make(chan int)
		go ChanRelay(nil, source, sink)

		source <- 1
		assert.Equal(t, 1, <-sink)
	})

	t.Run("cancel", func(t *testing.T) {
		source, sink := make(chan int), make(chan int)
		ctx, cancel := context.WithCancel(context.TODO())
		go ChanRelay(ctx, source, sink)

		cancel()

		writeSourceBlocks := false
		select {
		case source <- 1:
		default:
			writeSourceBlocks = true
		}
		assert.True(t, writeSourceBlocks)
	})

	t.Run("cancel write sink", func(t *testing.T) {
		source, sink := make(chan int), make(chan int)
		ctx, cancel := context.WithCancel(context.TODO())
		go ChanRelay(ctx, source, sink)

		source <- 1
		cancel()

		readSinkBlocks := false
		select {
		case <-sink:
		default:
			readSinkBlocks = true
		}
		assert.True(t, readSinkBlocks)
	})

	t.Run("close source", func(t *testing.T) {
		source, sink := make(chan int), make(chan int)
		go ChanRelay(nil, source, sink)
		close(source)

		readSinkBlocks := false
		select {
		case <-sink:
		default:
			readSinkBlocks = true
		}
		assert.True(t, readSinkBlocks)
	})
}

func TestChanMerge(t *testing.T) {
	source1, source2 := make(chan int), make(chan int)
	ctx, cancel := context.WithCancel(context.TODO())
	source := ChanMerge(ctx, 0, source1, source2)

	source1 <- 1
	source2 <- 2
	assert.Equal(t, 1, <-source)
	assert.Equal(t, 2, <-source)

	cancel()
	_, more := <-source
	assert.False(t, more)

	source = ChanMerge(nil, 0, make(<-chan int), make(<-chan int))
}

func TestChanStream(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		source := ChanStream(nil, 0, 1, 2, 3)

		assert.Equal(t, 1, <-source)
		assert.Equal(t, 2, <-source)
		assert.Equal(t, 3, <-source)
		_, more := <-source
		assert.False(t, more)
	})

	t.Run("quit", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		source := ChanStream(ctx, 0, 1, 2, 3)

		assert.Equal(t, 1, <-source)
		assert.Equal(t, 2, <-source)
		cancel()
		time.Sleep(time.Millisecond)

		_, more := <-source
		assert.False(t, more)
	})
}
