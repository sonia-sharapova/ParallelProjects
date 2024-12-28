/*
I performed all operations on uint64 because since this is the largest size,
it will most likely hold any unsigned integer passed to it (reducing the chance
of overflow). The documentation for WaitGroup also used 64-bit integers.

Atomic operatoins used:
- atomic.AddUint64(): add to counter
- atomic.LoadUint64(): read from memory

References:
- Returning an interface: https://stackoverflow.com/questions/35006640/function-to-return-an-interface
- Implementing own counter (for uint32): https://stackoverflow.com/questions/68995144/how-to-get-the-number-of-goroutines-associated-with-a-waitgroup
- Custom WorkGroups: https://blog.stackademic.com/creating-your-own-waitgroup-in-go-b2c80178d6d0
- WaitGroup documentation: https://cs.opensource.google/go/go/+/refs/tags/go1.23.2:src/sync/waitgroup.go
- Atomic operations info: https://medium.com/@deckarep/the-go-1-19-atomic-wrappers-and-why-to-use-them-ae14c1177ad8
- Decrement by 1 in atomic.AddUint64: https://pkg.go.dev/sync/atomic#AddUint64
- atomic.AddUint64() documentation: https://pkg.go.dev/sync/atomic
- atomic.LoadUint64() documentation: https://pkg.go.dev/sync/atomic#LoadUint64
*/

package waitgroup

import (
	"sync/atomic"
)

// From the Go Documentation:
// WaitGroup: waits for a collection of goroutines to finish.
type WaitGroup interface {
	Add(amount uint)
	Done()
	Wait()
}

// Define custom WaitGroup struct
// counter: number of goroutines to wait for (using uint)
type CustomWaitGroup struct {
	counter uint64
}

// NewWaitGroup returns a instance of a waitgroup
// This instance must be a pointer and should not be copied after creation.
func NewWaitGroup() WaitGroup {
	return &CustomWaitGroup{}
}

// Add(): with uint passed
//   - set the number of goroutines to wait for
//   - increments the counter by the specified amount
func (wg *CustomWaitGroup) Add(delta uint) {
	// use atomic.AddUint32 for compatibility
	atomic.AddUint64(&wg.counter, uint64(delta))
}

// Done():
//   - Called when goroutine is finished
//   - Decrement the counter by 1
func (wg *CustomWaitGroup) Done() {
	// From documentation: to decrement x, you can use: AddUint32(&x, ^uint32(0)).
	atomic.AddUint64(&wg.counter, ^uint64(0))
}

// Wait():
//   - Block until all goroutines have finished
//   - AKA: Until the counter becomes zero
func (wg *CustomWaitGroup) Wait() {
	// read counter value from memory using atomic.LoadUint64
	for atomic.LoadUint64(&wg.counter) != 0 {
		// Waiting for counter
	}
}
