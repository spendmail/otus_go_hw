package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

// Returns count of list elements.
func (l *list) Len() int {
	return l.len
}

// Returns the first element of the list.
func (l *list) Front() *ListItem {
	return l.front
}

// Returns the last element of the list.
func (l *list) Back() *ListItem {
	return l.back
}

// Puts data to the ListItem and marks it as front.
func (l *list) PushFront(v interface{}) *ListItem {
	listItem := &ListItem{Value: v, Next: nil, Prev: nil}

	if l.front == nil {
		// If list is empty, marks element as both front and back
		l.front = listItem
		l.back = listItem
	} else {
		// Otherwise shifts front to the right
		listItem.Next = l.front
		l.front.Prev = listItem
		l.front = listItem
	}
	l.len++

	return l.front
}

// Puts data to the ListItem and marks it as back.
func (l *list) PushBack(v interface{}) *ListItem {
	listItem := &ListItem{Value: v, Next: nil, Prev: nil}

	if l.back == nil {
		// If list is empty, marks element as both front and back
		l.front = listItem
		l.back = listItem
	} else {
		// Otherwise shifts back to the left
		listItem.Prev = l.back
		l.back.Next = listItem
		l.back = listItem
	}
	l.len++

	return l.back
}

// Removes element from the list and adjusts appropriate pointers.
func (l *list) Remove(i *ListItem) {
	// If item is only one element
	if l.len == 1 {
		l.front = nil
		l.back = nil
		i = nil
		l.len--
		return
	}

	// If item is front element
	if i == l.front {
		l.front = l.front.Next
		l.front.Prev = nil
		i = nil
		l.len--
		return
	}

	// If item is back element
	if i == l.back {
		l.back.Prev = nil
		l.back = nil
		i = nil
		l.len--
		return
	}

	// If item is neither front nor back
	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	i = nil
	l.len--
}

// Shifts element to front and adjusts appropriate pointers.
func (l *list) MoveToFront(i *ListItem) {
	// If item is already at the front
	if i == l.front {
		return
	}

	// If item is at the back
	if i == l.back {
		i.Prev.Next = nil
		l.back = i.Prev
		i.Prev = nil
		i.Next = l.front
		l.front.Prev = i
		l.front = i
		return
	}

	// If item is neither at the front nor back
	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	i.Next = l.front
	l.front.Prev = i
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

// List constructor: returns new list instance.
func NewList() List {
	return new(list)
}
