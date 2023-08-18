package sly

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestNewPriorityQueue(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		opts := PriorityQueueOptions[int]{
			InitialCap: 1000,
			Limit:      2000,
			Compare:    CompareOrdered[int],
		}
		pq, err := NewPriorityQueue(opts)
		assert.NotNil(t, pq)
		assert.NoError(t, err)
	})

	t.Run("nil comparator", func(t *testing.T) {
		opts := PriorityQueueOptions[int]{
			InitialCap: 1000,
			Limit:      2000,
		}
		pq, err := NewPriorityQueue(opts)
		assert.Nil(t, pq)
		assert.ErrorIs(t, err, ErrBadOptions)
	})
}

func TestPriorityQueue(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		pq, err := NewPriorityQueue(PriorityQueueOptions[int]{
			Limit:   3,
			Compare: CompareOrdered[int],
		})
		assert.NotNil(t, pq)
		assert.NoError(t, err)

		assert.True(t, pq.TryPush(1))
		assert.True(t, pq.TryPush(2))
		assert.True(t, pq.TryPush(3))
		assert.False(t, pq.TryPush(4))

		v, ok := pq.Pop()
		assert.Equal(t, 3, v)
		assert.True(t, ok)

		rest := pq.Close()
		assert.Equal(t, []int{2, 1}, rest)
		rest = pq.Close()
		assert.Equal(t, []int{2, 1}, rest)

		assert.False(t, pq.TryPush(5))
		v, ok = pq.Pop()
		assert.Equal(t, 0, v)
		assert.False(t, ok)
	})

	t.Run("pop block", func(t *testing.T) {
		pq, _ := NewPriorityQueue(PriorityQueueOptions[int]{
			Limit:   3,
			Compare: CompareOrdered[int],
		})

		recv := make(chan int)
		go func() {
			defer close(recv)
			for {
				v, ok := pq.Pop()
				if !ok {
					return
				}
				recv <- v
			}
		}()

		wg := sync.WaitGroup{}
		wg.Add(3)
		go func() {
			for range recv {
				wg.Done()
			}
		}()

		// Giving RW goroutines some time to start.
		time.Sleep(10 * time.Millisecond)
		pq.TryPush(1)
		pq.TryPush(2)
		pq.TryPush(3)
		wg.Wait()
		pq.Close()
	})

	t.Run("pop unblock on close", func(t *testing.T) {
		pq, _ := NewPriorityQueue(PriorityQueueOptions[int]{
			Limit:   3,
			Compare: CompareOrdered[int],
		})

		block := make(chan struct{})
		go func() {
			_, _ = pq.Pop()
			close(block)
		}()
		pq.Close()

		<-block
	})

	t.Run("close cutoff", func(t *testing.T) {
		pq, _ := NewPriorityQueue(PriorityQueueOptions[int]{
			Compare: CompareOrdered[int],
		})
		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				for {
					pq.TryPush(rand.Int())
				}
			}()
		}
		for i := 0; i < runtime.NumCPU()*4; i++ {
			go func() {
				_, _ = pq.Pop()
			}()
		}
		time.Sleep(10 * time.Millisecond)
		pq.Close()
	})
}

func BenchmarkPriorityQueue(b *testing.B) {
	pq, _ := NewPriorityQueue(PriorityQueueOptions[int]{
		Limit:   4096 * 32,
		Compare: CompareOrdered[int],
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pq.TryPush(rand.Int())
			_, _ = pq.Pop()
		}
	})
}
