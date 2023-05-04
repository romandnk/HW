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

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len  int
	head *ListItem
	tail *ListItem
}

func (list *list) Len() int {
	return list.len
}

func (list *list) Front() *ListItem {
	return list.head
}

func (list *list) Back() *ListItem {
	return list.tail
}

func (list *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if list.len == 0 {
		list.head = item
		list.tail = item
	} else {
		item.Next = list.head
		list.head.Prev = item
		list.head = item
	}
	list.len++
	return item
}

func (list *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if list.len == 0 {
		list.head = item
		list.tail = item
	} else {
		item.Prev = list.tail
		list.tail.Next = item
		list.tail = item
	}
	list.len++
	return item
}

func (list *list) Remove(i *ListItem) { // 10 (20) 30 40
	if i.Prev == nil {
		i.Next.Prev = nil
		list.head = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	if i.Next == nil {
		i.Prev.Next = nil
		list.tail = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	list.len--
}

func (list *list) MoveToFront(i *ListItem) { // 10 20 30 40 50
	list.Remove(i)
	list.len++
	list.Front().Prev = i
	i.Next = list.Front()
	i.Prev = nil
	list.head = i
}

func NewList() List {
	return new(list)
}
