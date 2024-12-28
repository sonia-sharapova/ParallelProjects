package queue

import (
	"sync/atomic"
	"unsafe"
)

// Task to be processed
type Request struct {
	Command   string
	ID        int
	Body      string
	Timestamp float64
}

// node represents a single node in the queue
type node struct {
	value *Request			// Task at this node
	next  unsafe.Pointer 	// Pointer to next node
}

// LockFreeQueue represents a FIFO queue with lock-free operations
type LockFreeQueue struct {
	head unsafe.Pointer		// Pointer to head node
	tail unsafe.Pointer		// to tail node
}

// NewLockFreeQueue creates and initializes a LockFreeQueue
func NewLockFreeQueue() *LockFreeQueue {
	n := &node{}			// Node to initialize queue
	q := &LockFreeQueue{
		head: unsafe.Pointer(n),
		tail: unsafe.Pointer(n),
	}
	return q
}

// Enqueue adds a Request to the queue
func (q *LockFreeQueue) Enqueue(task *Request) {
	newNode := &node{value: task}	// Create a new node for the task
	var tail, next *node

	for {
		// Load the current tail and the next node in the queue
		tail = (*node)(atomic.LoadPointer(&q.tail))
		next = (*node)(atomic.LoadPointer(&tail.next))

		// Check if the tail is consistent
		if tail == (*node)(atomic.LoadPointer(&q.tail)) { 
			if next == nil {
				// If next is nil, the tail is at the end of the queue
				// Try to link the new node
				if atomic.CompareAndSwapPointer(
					(*unsafe.Pointer)(unsafe.Pointer(&tail.next)),
					unsafe.Pointer(next),
					unsafe.Pointer(newNode),
				) {
					// Advance the tail pointer
					atomic.CompareAndSwapPointer(
						&q.tail,
						unsafe.Pointer(tail),
						unsafe.Pointer(newNode),
					)
					return
				}
			} else {
				// Tail is behind, try to advance it
				atomic.CompareAndSwapPointer(
					&q.tail,
					unsafe.Pointer(tail),
					unsafe.Pointer(next),
				)
			}
		}
	}
}

// Dequeue removes a Request from the queue
func (q *LockFreeQueue) Dequeue() (*Request, bool) {
	var head, tail, next *node
	var value *Request

	for {
		// load current nodes
		head = (*node)(atomic.LoadPointer(&q.head))
		tail = (*node)(atomic.LoadPointer(&q.tail))
		next = (*node)(atomic.LoadPointer(&head.next))

		// Check if the head is consistent
		if head == (*node)(atomic.LoadPointer(&q.head)) { 
			if head == tail {
				// If the head equals the tail, check if the queue is empty
				if next == nil {
					// Queue is empty
					return nil, false
				}
				// Tail is behind, try to advance it
				atomic.CompareAndSwapPointer(
					&q.tail,
					unsafe.Pointer(tail),
					unsafe.Pointer(next),
				)
			} else {
				// Read value before CAS to avoid ABA problem
				value = next.value
				if atomic.CompareAndSwapPointer(
					&q.head,
					unsafe.Pointer(head),
					unsafe.Pointer(next),
				) {
					// Successfully dequeued the value
					return value, true
				}
			}
		}
	}
}

// IsEmpty checks if the queue is empty
func (q *LockFreeQueue) IsEmpty() bool {
	head := (*node)(atomic.LoadPointer(&q.head))
	tail := (*node)(atomic.LoadPointer(&q.tail))
	next := (*node)(atomic.LoadPointer(&head.next))

	// The queue is empty if head equals tail and there is no next node
	return head == tail && next == nil
}
