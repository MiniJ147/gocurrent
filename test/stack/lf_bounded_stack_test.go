package stack_test

import (
	"fmt"
	"testing"

	"github.com/minij147/gocurrent/lockfree/lfstack"
)

func TestFoo(t *testing.T) {
	s := lfstack.NewBounded[int](3)
	fmt.Println(s.Peek())
	for i := range 4 {
		fmt.Println(s.Push(i))
		fmt.Println(s.Peek())
	}
}
