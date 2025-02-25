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
	Head    atomic.Pointer[node[T]]
	Tail    atomic.Pointer[node[T]]
	Sential *node[T]
}

func (q *queue[T]) Push(data T) {
	newNode := node[T]{data: data}

	for {
		tail := q.Tail.Load()    // snapshot of tail
		next := tail.next.Load() // next value of our snapshsot

		// someone has appended to our tail
		if next != nil || q.Tail.Load() != tail {
			q.Tail.CompareAndSwap(tail, next) // attempt to move this forward
			continue
		}

		// attempt to swap our tail will this value
		if tail.next.CompareAndSwap(next, &newNode) {
			q.Tail.CompareAndSwap(tail, &newNode)
			break
		}
	}
}

func (q *queue[T]) Pop() (T, bool) {
	var x T

	for {
		head := q.Head.Load()
		tail := q.Tail.Load()
		next := head.next.Load()

		if q.Head.Load() != head {
			continue
		}

		if head == tail {
			if next == nil {
				return x, false
			}

			// attempt to move tail forward -- help out other threads
			q.Tail.CompareAndSwap(tail, next)
			continue
		}

		val := next.data
		if q.Head.CompareAndSwap(head, next) {
			return val, true
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
