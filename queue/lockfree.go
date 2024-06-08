package queue

import (
	"sync/atomic"
	"unsafe"
)

type Request struct {
	Command   string
	Id        int
	Body      string
	Timestamp int
	next      *Request
}

// LockfreeQueue represents a FIFO structure with operations to enqueue
// and dequeue tasks represented as Request
type LockFreeQueue struct {
	head *Request
	tail *Request
}

func NewRequest(v *Request) *Request {
	return &Request{Command: v.Command, Id: v.Id, Body: v.Body, Timestamp: v.Timestamp}
}

func NewLockFreeQueue() *LockFreeQueue {
	dummyNode := &Request{}
	return &LockFreeQueue{head: dummyNode, tail: dummyNode}
}

func (q *LockFreeQueue) Enqueue(request *Request) {
	for {
		tail := q.tail
		next := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next)))
		if tail == q.tail {
			if next == nil {
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next)), unsafe.Pointer(next), unsafe.Pointer(request)) {
					atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(tail), unsafe.Pointer(request))
					return
				}
			} else {
				atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(tail), next)
			}
		}
	}
}

func (q *LockFreeQueue) Dequeue() *Request {
	for {
		head := q.head
		tail := q.tail
		first := head.next
		if head == q.head {
			if head == tail {
				if first == nil {
					return nil
				}
				atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(tail), unsafe.Pointer(first))
			} else {
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)), unsafe.Pointer(head), unsafe.Pointer(first)) {
					return first
				}
			}
		}
	}
}

func (q *LockFreeQueue) IsEmpty() bool {
	head := q.head
	tail := q.tail
	return head == tail && head.next == nil
}
