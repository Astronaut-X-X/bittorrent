package routing

import (
	"container/list"
	"sync"
)

type SyncList struct {
	*sync.RWMutex
	List *list.List
}

func NewSyncList() *SyncList {
	return &SyncList{
		RWMutex: &sync.RWMutex{},
		List:    list.New(),
	}
}

func (l *SyncList) Front() *list.Element {
	l.RLock()
	defer l.RUnlock()

	return l.List.Front()
}

func (l *SyncList) Back() *list.Element {
	l.RLock()
	defer l.RUnlock()

	return l.List.Back()
}

func (l *SyncList) PushFront(v interface{}) *list.Element {
	l.Lock()
	defer l.Unlock()

	return l.List.PushFront(v)
}

func (l *SyncList) PushBack(v interface{}) *list.Element {
	l.Lock()
	defer l.Unlock()

	return l.List.PushBack(v)
}

func (l *SyncList) InsertBefore(
	v interface{}, mark *list.Element) *list.Element {

	l.Lock()
	defer l.Unlock()

	return l.List.InsertBefore(v, mark)
}

func (l *SyncList) InsertAfter(
	v interface{}, mark *list.Element) *list.Element {

	l.Lock()
	defer l.Unlock()

	return l.List.InsertAfter(v, mark)
}

func (l *SyncList) Remove(e *list.Element) interface{} {
	l.Lock()
	defer l.Unlock()

	return l.List.Remove(e)
}

func (l *SyncList) Clear() {
	l.Lock()
	defer l.Unlock()

	l.List.Init()
}

func (l *SyncList) Len() int {
	l.RLock()
	defer l.RUnlock()

	return l.List.Len()
}

func (l *SyncList) Elements() []*list.Element {
	elems := make([]*list.Element, 0, l.List.Len())
	l.RLock()
	for e := l.List.Front(); e != nil; e = e.Next() {
		elems = append(elems, e)
	}
	l.RUnlock()
	return elems
}
