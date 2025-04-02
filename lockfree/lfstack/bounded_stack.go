/*
BOUNDED:
Self designed stack with design principles borrowed from the queue in Art of Mulitprocessor programming
For reference and source go to the queue/lockfree README or source file

About:
Simple CAS lockfree design

Future:
add descriptor based design
add a actual proper modern paper implementation of a lockfree stack.
*/
package lfstack

import (
	"sync/atomic"
)

type stack[T any] struct {
	capacity uint32
	head     atomic.Uint32
	data     []atomic.Pointer[T]
}

func (s *stack[T]) Capacity() uint32 {
	return s.capacity
}

func (s *stack[T]) Size() uint32 {
	return 0
}

func (s *stack[T]) Push(val T) error {
	for {
		headExpected := s.head.Load()
		valueExpected := s.data[headExpected+1].Load()

		if headExpected == s.capacity {
			return ErrFullStack
		}

		// head is lagging
		if s.data[headExpected+1].Load() != nil {
			// help move it forward
			s.head.CompareAndSwap(headExpected, headExpected+1)
			continue
		}

		// now we can assume we are in the right spot
		// passing address to value's copy (which the garbage collector should take care of
		if s.data[headExpected+1].CompareAndSwap(valueExpected, &val) {
			s.head.CompareAndSwap(headExpected, headExpected+1)
			return nil
		}
	}
}

func (s *stack[T]) Pop() (T, error) {
	var tmp T
	return tmp, nil
}

func (s *stack[T]) Peek() (T, error) {
	var def T
	for {
		headExpected := s.head.Load()
		valueExpected := s.data[headExpected].Load()

		// head lagging
		if s.data[headExpected+1].Load() != nil {
			s.head.CompareAndSwap(headExpected, headExpected+1)
			continue
		}

		if valueExpected == nil {
			return def, ErrEmptyStack
		}

		return *valueExpected, nil
	}
}

func NewBounded[T any](capacity uint32) Stack[T] {
	return &stack[T]{
		capacity: capacity,
		head:     atomic.Uint32{},
		data:     make([]atomic.Pointer[T], capacity+2), // +2 for sential
	}
}
