package slot

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	s := NewSlot()

	for i := 0; i < 10; i++ {
		a := s.GetTable()

		// fmt.Println(a)
		fmt.Println(s.Result(a))
	}
}

func TestB(t *testing.T) {
	s := NewSlot()
	x := [][]int32{[]int32{1, 0, 2}, []int32{2, 0, 1}, []int32{3, 0, 4}, []int32{5, 0, 6}, []int32{5, 0, 3}}
	fmt.Println(s.Result(x))

	xx := [][]int32{[]int32{0, 0, 0}, []int32{0, 0, 0}, []int32{0, 0, 0}, []int32{0, 0, 0}, []int32{0, 0, 0}}
	fmt.Println(s.Result(xx))
}
