package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	if lru.capacity == 0 {
		return false
	}
	if item, ok := lru.items[key]; ok {
		item.Value = value
		lru.queue.MoveToFront(item)
		return true
	}
	if lru.capacity == lru.queue.Len() {
		delete(lru.items, key)
		lru.queue.Remove(lru.queue.Back())
	} else {
		lru.items[key] = lru.queue.PushFront(value)
	}
	return false
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := lru.items[key]; ok {
		lru.queue.MoveToFront(item)
		return item.Value, true
	}
	return nil, false
}

func (lru *lruCache) Clear() {
	*lru = *(NewCache(lru.capacity).(*lruCache))
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
