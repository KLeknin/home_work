package hw04lrucache

import (
	"sync/atomic"
)

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
	Key   Key
}

type list struct {
	len   int32
	front *ListItem
	back  *ListItem
}

func newList() *list {
	return new(list)
}

func (l *list) Len() int {
	return (int(l.len))
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	second := l.front
	l.front = new(ListItem)
	l.front.Value = v
	l.front.Next = second
	if second != nil {
		second.Prev = l.front
	}
	if l.back == nil {
		l.back = l.front
	}
	atomic.AddInt32(&l.len, 1)
	return (l.front)
}

func (l *list) PushBack(v interface{}) *ListItem {
	last := l.back
	l.back = new(ListItem)
	l.back.Value = v
	l.back.Prev = last
	if last != nil {
		last.Next = l.back
	}
	if l.front == nil {
		l.front = l.back
	}
	atomic.AddInt32(&l.len, 1)
	return (l.back)
}

func (l *list) Remove(i *ListItem) {
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}
	atomic.AddInt32(&l.len, -1)
	if l.len == 0 {
		l.front = nil
		l.back = nil
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	i.Prev.Next = i.Next
	first := l.front
	i.Next = first
	i.Prev = nil
	first.Prev = i
	l.front = i
}
