package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

// type cacheItem struct {
// 	key   Key
// 	value interface{}
// }

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    newList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	item, ok := lc.items[key]
	if !ok {
		item = lc.queue.PushFront(value)
		item.Key = key
		lc.items[key] = item
		if lc.queue.Len() > lc.capacity {
			delete(lc.items, lc.queue.Back().Key)
			lc.queue.Remove(lc.queue.Back())
		}
	}
	item.Value = value
	return ok
}

func (lc *lruCache) Get(key Key) (i interface{}, ok bool) {
	item, ok := lc.items[key]
	if ok {
		i = item.Value
	}
	return
}

func (lc *lruCache) Clear() {
	for lc.queue.Len() > 0 {
		delete(lc.items, lc.queue.Back().Key)
		lc.queue.Remove(lc.queue.Back())
	}
}
