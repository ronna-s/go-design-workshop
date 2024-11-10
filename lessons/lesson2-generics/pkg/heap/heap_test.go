package heap

import (
	"testing"

	"golang.org/x/exp/constraints"
)

type OrderedItem[T constraints.Ordered] struct {
	Val T
}

type intHeap = Heap[OrderedItem[int]]

func (i OrderedItem[T]) Less(j OrderedItem[T]) bool {
	return i.Val < j.Val
}

func verify(t *testing.T, h intHeap, i int) {
	t.Helper()
	n := len(h)
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if h[j1].Less(h[i]) {
			t.Errorf("error")
			return
		}
		verify(t, h, j1)
	}
	if j2 < n {
		if h[j2].Less(h[i]) {
			t.Errorf("error")
			return
		}
		verify(t, h, j2)
	}
}

func TestHeap(t *testing.T) {
	var h intHeap

	verify(t, h, 0)
	for i := 20; i > 10; i-- {
		h.Push(OrderedItem[int]{i})
	}

	verify(t, h, 0)
	for i := 10; i > 0; i-- {
		h.Push(OrderedItem[int]{i})
		verify(t, h, 0)
	}

	for i := 1; len(h) > 0; i++ {
		x := h.Pop()
		if i < 20 {
			h.Push(OrderedItem[int]{20 + i})
		}
		verify(t, h, 0)
		y := OrderedItem[int]{i}
		if x != y {
			t.Errorf("%v.th pop got %v; want %v", i, x.Val, i)
		}
	}
}