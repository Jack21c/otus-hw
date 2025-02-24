package hw04lrucache

import "fmt"

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
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newList := ListItem{
		Value: v,
		Prev:  nil,
		Next:  l.front,
	}
	if l.front != nil {
		l.front.Prev = &newList
	}
	l.front = &newList
	if l.back == nil {
		l.back = &newList
	}
	l.len++
	return &newList
}

func (l *list) PushBack(v interface{}) *ListItem {
	newList := ListItem{
		Value: v,
		Prev:  l.back,
		Next:  nil,
	}
	if l.back != nil {
		l.back.Next = &newList
	}
	l.back = &newList
	if l.front == nil {
		l.front = &newList
	}
	l.len++
	return &newList
}

func (l *list) Remove(i *ListItem) {
	l.len--
	switch i {
	case nil:
		l.len++
		fmt.Println("nil item")
	case l.front:
		l.front = i.Next
		if i.Next != nil {
			i.Next.Prev = nil
		}
		if l.back == i {
			l.back = nil
		}
	case l.back:
		l.back = i.Prev
		if i.Prev != nil {
			i.Prev.Next = nil
		}
	default:
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	i.Next = l.front
	i.Prev = nil
	if l.front != nil {
		l.front.Prev = i
	}
	l.front = i
}
