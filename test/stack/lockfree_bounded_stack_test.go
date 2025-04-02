package stack_test

import (
	"testing"

	stack "github.com/minij147/gocurrent/stack/lockfree"
)

func TestFoo(t *testing.T) {
	s := stack.NewBounded[int](10)
	s.Push(10)
}
