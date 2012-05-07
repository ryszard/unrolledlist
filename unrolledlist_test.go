package unrolledlist

import (
	"testing"
	"fmt"
)

// Helper methods that give you more introspection into the list.

// nodeLength returns the length of list l in nodes.
func (l UnrolledList) nodeLength() int {
	if l.next == nil {
		return 1
	}
	return 1 + l.next.nodeLength()
}

func (l UnrolledList) repr() string {
	if l.next == nil {
		return fmt.Sprintf("%v", l.elements)
	}
	return fmt.Sprintf("%v %v", l.elements, l.next.repr())
}

// newPopulatedList creates a new UnrolledList and appends integers
// from 0 to n-1 to it.
func newPopulatedList(capacity, n int) (list *UnrolledList) {
	list = New(capacity)
	for i := 0; i < n; i++ {
		list.Append(i)
	}
	return
}

func ListLike(t *testing.T, list *UnrolledList, values ...int) {
	wasError := false
	for i, wanted := range values {
		if value := list.Get(i); value != wanted {
			wasError = true
			t.Errorf("Wrong value for index %d: %v (should be %v).", i, value, wanted)
		}
	}
	if wasError {
		t.Errorf("(%v)", list.repr())
	}
}

func TestAppendOneNode(t *testing.T) {
	list := newPopulatedList(10, 5)
	ListLike(t, list, 0, 1, 2, 3, 4)
}

func TestMoreNodes(t *testing.T) {
	list := newPopulatedList(3, 5)
	ListLike(t, list, 0, 1, 2, 3, 4)
}

func TestDowncaseInsert(t *testing.T) {
	sl := make([]interface{}, 3, 3)
	sl[0], sl[1] = 0, 2
	sl = insert(sl, 1, 1)
	if sl[0] != 0 || sl[1] != 1 || sl[2] != 2 {
		t.Errorf("slice is not ordered: %v.", sl)
	}
}

func TestInsertOneNode(t *testing.T) {
	list := New(5)
	list.Append(0)
	list.Append(2)
	list.Insert(1, 1)
	ListLike(t, list, 0, 1, 2)

}

func TestInsertIntoLastNode(t *testing.T) {
	list := newPopulatedList(4, 7)
	list.Insert(6, 1000)
	ListLike(t, list, 0, 1, 2, 3, 4, 5, 1000, 6)
}

func TestInsertIntoMiddleNode(t *testing.T) {
	list := newPopulatedList(3, 7)
	list.Insert(3, 1000)
	ListLike(t, list, 0, 1, 2, 1000, 3, 4, 5, 6)
}

func TestInsertOutOfBounds(t *testing.T) {
	list := New(3)
	if err := list.Insert(100, 1); err == nil {
		t.Error("Out of bound insert didn't return an error.")
	}
}

func TestIteration(t *testing.T) {
	list := newPopulatedList(3, 10)
	wanted := 0
	for el := range list.Iter() {
		if wanted != el {
			t.Errorf("Wanted %v, got %v.", wanted, el)
		}
		wanted++
	}
	if wanted != 10 {
		t.Errorf("Iter was supposed to yield 10 times, yielded %d.", wanted)
	}
}

func TestSliceDelete(t *testing.T) {
	arr := []interface{}{0, 1, 2}
	el, arr := sliceDelete(arr, 1)
	if el != 1 {
		t.Errorf("Got %v, expected 1 (%v)", el, arr)
	}
	if arr[0] != 0 || arr[1] != 2 {
		t.Errorf("Got %v, expected [0, 2]", arr)
	}
}

func TestPopSimple(t *testing.T) {
	list := newPopulatedList(5, 3)

	nodeLength := list.nodeLength()

	if el := list.Pop(1); el != 1 {
		t.Errorf("Got %v, expected 1.", el)
	}
	ListLike(t, list, 0, 2)
	if newNodeLength := list.nodeLength(); newNodeLength != nodeLength {
		t.Errorf("Node length has changed even though it shouldn't have: %v, not %v.", newNodeLength, nodeLength)
	}
}

func TestPopNotInFirstNode(t *testing.T) {
	list := newPopulatedList(3, 10)

	nodeLength := list.nodeLength()

	if el := list.Pop(5); el != 5 {
		t.Errorf("Got %v, expected 1", el)
	}
	ListLike(t, list, 0, 1, 2, 3, 4, 6, 7, 8, 9)
	if newNodeLength := list.nodeLength(); newNodeLength != nodeLength {
		t.Errorf("Node length has changed even though it shouldn't have: %v, not %v.", newNodeLength, nodeLength)
	}

}

func TestPopNodeMoveElementsFromAdjacent(t *testing.T) {
	list := newPopulatedList(4, 8)
	// [0, 1, 2, 3], [4, 5, 6, 7]

	for i := 0; i < 3; i++ {
		list.Pop(0)
	}

	if l1, l2 := len(list.elements), len(list.next.elements); l1 != 2 || l2 != 3 {
		t.Errorf("The layout of the elements is wrong. It should be ([3 4] [5 6 7]), (%v) found.", list.repr())
	}
	ListLike(t, list, 3, 4, 5, 6, 7)

}

func TestPopNodeMoveElementsWithMerge(t *testing.T) {
	list := newPopulatedList(4, 12)
	if nl := list.nodeLength(); nl != 3 {
		t.Errorf("The node length isn't 3: %v", nl)
	}

	for i := 0; i < 4; i++ {
		list.Pop(0)
	}
	if nl, l1, l2 := list.nodeLength(), len(list.elements), len(list.next.elements); nl != 2 || l1 != 4 || l2 != 4 {
		t.Errorf("The layout of the elements is wrong. It should be ([4 5 6 7] [8 9 10 11]), (%v) found.", list.repr())
	}

	ListLike(t, list, 4, 5, 6, 7, 8, 9, 10, 11)

}

func TestOutOfBounds(t *testing.T) {
	list := newPopulatedList(3, 10)

	if el := list.Pop(100); el != nil {
		t.Errorf("Out of bound element should be nil, not %v.", el)
	}

	list = newPopulatedList(3, 1)
	if el := list.Pop(1); el != nil {
		t.Errorf("Out of bound element should be nil, not %v.", el)
	}

}

func TestLength(t *testing.T) {
	list := New(3)
	if l := list.Length(); l != 0 {
		t.Errorf("An empty list should have length 0, not %v.", l)
	}
	for i := 0; i < 10; i++ {
		list.Append(i)
	}
	if l := list.Length(); l != 10 {
		t.Errorf("Wrong length: %v instead of 10.", l)
	}

}

func TestAppendReallyAppends(t *testing.T) {
	list := newPopulatedList(5, 10)
	list.Pop(0)
	list.Pop(0)
	if l := list.nodeLength(); l != 2 {
		t.Errorf("Expected to see two nodes.Seen %d: %v", l, list.repr())
	}

	list.Append(1000)
	ListLike(t, list, 2, 3, 4, 5, 6, 7, 8, 9, 1000)
}
