/*
lock-free queue based on "The Art of Multiprocessor Programming"
Plans to further protect queue and insure correctness along with a
better more modern queue based off a more modern paper
*/
package queue

import (
	"sync/atomic"
)

type node[T any] struct {
	data T
	next atomic.Pointer[node[T]]
}

type queue[T any] struct {
	head    atomic.Pointer[node[T]]
	tail    atomic.Pointer[node[T]]
	sential *node[T]
}

// Push into the Queue
func (q *queue[T]) Push(data T) {
	newNode := node[T]{data: data}

	for {
		tail := q.tail.Load()    // snapshot of tail
		next := tail.next.Load() // next value of our snapshsot

		// someone has appended to our tail
		if next != nil || q.tail.Load() != tail {
			q.tail.CompareAndSwap(tail, next) // attempt to move this forward
			continue
		}

		// attempt to swap our tail will this value
		if tail.next.CompareAndSwap(next, &newNode) {
			q.tail.CompareAndSwap(tail, &newNode)
			break
		}
	}
}

// Pop from the Queue
// Does not block, returns true if successful
func (q *queue[T]) Pop() (T, bool) {
	var x T

	for {
		head := q.head.Load()
		tail := q.tail.Load()
		next := head.next.Load()

		if q.head.Load() != head {
			continue
		}

		if head == tail {
			if next == nil {
				return x, false
			}

			// attempt to move tail forward -- help out other threads
			q.tail.CompareAndSwap(tail, next)
			continue
		}

		val := next.data
		if q.head.CompareAndSwap(head, next) {
			return val, true
		}
	}
}

// Pop from Queue while blocking until we receive a value
func (q *queue[T]) PopBlocking() T {
	for {
		res, ok := q.Pop()
		if ok {
			return res
		}
	}
}

func New[T any]() queue[T] {
	var sential *node[T] = &node[T]{}
	var head, tail atomic.Pointer[node[T]]

	head.Store(sential)
	tail.Store(sential)

	return queue[T]{head, tail, sential}
}
