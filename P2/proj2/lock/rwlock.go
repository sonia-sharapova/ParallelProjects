// Package lock provides an implementation of a read-write lock
// that uses condition variables and mutexes.
package lock

import "sync"

// RWLock: custom read-write lock implementation.
// allows multiple readers or a single writer at a time.
type RWLock struct {
	mu             sync.Mutex  		// mutex ( for synchronizing access)
	cond           *sync.Cond  		// condition variable (for signaling threads)
	readers        int         		// current num of active readers
	writer         bool        		// whether writer is active
	waitingWriters int         		// num of writers waiting for the lock
}

const maxReaders = 32

// NewRWLock creates and returns a new RWLock instance
func NewRWLock() *RWLock {
	lock := &RWLock{}
	lock.cond = sync.NewCond(&lock.mu) // Initialize condition variable with the mutex
	return lock
}

// Lock acquires the write lock, blocking other readers and writers
func (rw *RWLock) Lock() {
	rw.mu.Lock() // Lock the mutex to synchronize access
	defer rw.mu.Unlock() // Ensure unlocking

	rw.waitingWriters++ // Increment count of waiting writers

	// Wait until there are no active readers or writers
	for rw.readers > 0 || rw.writer {
		rw.cond.Wait()
	}
	rw.waitingWriters-- 
	rw.writer = true    // Writer is active
}

// Unlock releases the write lock
func (rw *RWLock) Unlock() {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	// Ensure the lock is held by a writer before unlocking
	if !rw.writer {
		panic("Unlock called when no writer holds the lock")
	}
	rw.writer = false		// Wwriter is inactive
	rw.cond.Broadcast()		// Signal all waiting threads
}

// RLock acquires a read lock, allowing multiple readers but no writers
func (rw *RWLock) RLock() {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	// Wait until there is no active writer, and readers are below the limit
	for rw.writer || rw.readers >= maxReaders || rw.waitingWriters > 0 {
		rw.cond.Wait()
	}
	rw.readers++ // Increment the count of active readers
}

// RUnlock releases a read lock
func (rw *RWLock) RUnlock() {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	// Ensure there is at least one active reader before unlocking
	if rw.readers == 0 {
		panic("RUnlock called when no readers hold the lock")
	}
	rw.readers--
	// If there are no more readers, signal waiting threads
	if rw.readers == 0 {
		rw.cond.Broadcast()
	}
}
