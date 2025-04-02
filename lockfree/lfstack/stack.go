package lfstack

type Stack[T any] interface {
	Capacity() uint32
	Size() uint32
	Push(val T) error
	Pop() (T, error)
	Peek() (T, error)
}
