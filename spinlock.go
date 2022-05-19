package sly

import (
	"runtime"
	"sync/atomic"
)

// Spinlock is an atomic based active lock.
type Spinlock int32

// Lock acquires the lock.
func (v *Spinlock) Lock() {
	backoff := 1
	for !atomic.CompareAndSwapInt32((*int32)(v), 0, 1) {
		for i := 0; i < backoff; i++ {
			runtime.Gosched()
		}
		if backoff < 16 {
			backoff <<= 1
		}
	}
}

// Unlock releases the lock.
func (v *Spinlock) Unlock() {
	atomic.StoreInt32((*int32)(v), 0)
}
