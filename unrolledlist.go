// The package unrolledlist implements an unrolled linked list, as
// described in http://en.wikipedia.org/wiki/Unrolled_linked_list. An
// unrolled linked list is similar to a normal linked list, but each
// of its CARs contains more than one value.
package unrolledlist

import "errors"

var (
	ErrOutOfBound = errors.New("index out of range")
)

type UnrolledList struct {
	next     *UnrolledList
	elements []interface{}
}

// Len Returns the length of l.
func (l UnrolledList) Len() int {
	if l.next == nil {
		return len(l.elements)
	}
	return len(l.elements) + l.next.Len()
}

// Append makes value the last element of l.
func (l *UnrolledList) Append(value interface{}) {
	if l.next == nil && len(l.elements) < cap(l.elements) {
		l.elements = append(l.elements, value)
	} else {
		if l.next == nil {
			l.grow()
		}
		l.next.Append(value)
	}
}

// grow adds a new node as the next of list l. The next pointer of the
// new node will be pointing at whatever was the next pointer of l
// pointing at.
func (l *UnrolledList) grow() {
	nextNext := l.next
	l.next = New(cap(l.elements))
	l.next.next = nextNext
}

func (l UnrolledList) iterate(ch chan interface{}) {
	for _, el := range l.elements {
		ch <- el
	}
	if l.next == nil {
		close(ch)
	} else {
		l.next.iterate(ch)
	}
}

// Iter gives you an iterator for list l (which you can use with
// range). Returns a channel that will have all the elements from l
// passed to by a goroutine.
func (l UnrolledList) Iter() chan interface{} {
	ch := make(chan interface{}, 1)
	go l.iterate(ch)
	return ch
}

// growDividing adds a new node after l, and puts half of l's elements
// into the new node.
func (l *UnrolledList) growDividing() {
	l.grow()
	half := len(l.elements) / 2
	l.next.elements = append(l.next.elements, l.elements[half:]...)
	l.elements = l.elements[:half]
}

// Get returns the element in the i-th position in l and nil if i <
// l.Len(). ok is true if a value from l was returned.
func (l UnrolledList) Get(i int) (value interface{}, ok bool) {
	length := len(l.elements)
	switch {
	case i < length:
		return l.elements[i], true
	case l.next != nil:
		return l.next.Get(i - length)
	}
	return nil, false
}

// insert is a helper function that inserts value at the i-th position
// in the slice sl. If 
func insert(sl []interface{}, i int, value interface{}) []interface{} {
	sl = append(sl[:i], append([]interface{}{value}, sl[i:]...)...)
	return sl
}

func sliceDelete(sl []interface{}, i int) (interface{}, []interface{}) {
	el := sl[i]
	return el, append(sl[:i], sl[i+1:]...)
}

// Insert inserts value at position i in the list l. If i >
// l.Len(), return ErrOutOfBound.
func (l *UnrolledList) Insert(i int, value interface{}) error {
	switch {
	case i >= len(l.elements):
		if l.next == nil {
			return ErrOutOfBound
		}
		l.next.Insert(i-len(l.elements), value)
	case len(l.elements) < cap(l.elements) && i <= cap(l.elements):
		l.elements = insert(l.elements, i, value)
	case len(l.elements) == cap(l.elements):
		l.growDividing()
		l.Insert(i, value)
	}
	return nil
}

// rebalance guarantees that all nodes except the last one are at
// least half full. If the current node becomes less than half full,
// move an element from the next ndoe to the current one. If it is
// possible to fit all the elements in both the current and next node
// in one node, do it.
func (l *UnrolledList) rebalance() {
	length := len(l.elements)
	capacity := cap(l.elements)
	switch {
	case l.next == nil:
		return
	case length+len(l.next.elements) <= capacity:
		l.elements = append(l.elements, l.next.elements...)
		l.next = l.next.next
	case capacity/2 > length:
		var toBeMoved interface{}
		toBeMoved, l.next.elements = sliceDelete(l.next.elements, 0)
		l.elements = append(l.elements, toBeMoved)
	}
}

// Delete and return the ith element of l. If i is out of bound,
// returns nil.
func (l *UnrolledList) Remove(i int) (popped interface{}) {
	length := len(l.elements)
	if i < length {
		popped, l.elements = sliceDelete(l.elements, i)
		l.rebalance()
	} else if l.next != nil {
		return l.next.Remove(i - length)
	}
	return
}

// Delete and return the first element of l.
func (l *UnrolledList) Pop() (popped interface{}) {
	return l.Remove(0)
}

// Return an empty Unrolled Linked List that will contain capacity
// elements in one cell.
func New(capacity int) *UnrolledList {
	return &UnrolledList{nil, make([]interface{}, 0, capacity)}
}
