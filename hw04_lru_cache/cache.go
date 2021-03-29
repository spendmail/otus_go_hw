package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

// LruCache setter: sets or updates value, depends on whether the value it exists or not.
func (l *lruCache) Set(key Key, value interface{}) bool {
	listItem, exists := l.items[key]
	cacheItemElement := cacheItem{string(key), value}

	if exists {
		// If cache element exists, move it to front
		listItem.Value = cacheItemElement
		l.queue.MoveToFront(listItem)
	} else {
		// If cache element doesn't exist, create
		listItem = l.queue.PushFront(cacheItemElement)

		// If list exceeds capacity, remove last element from list and map
		if l.queue.Len() > l.capacity {
			item := l.queue.Back()
			backCacheItem := item.Value.(cacheItem)
			delete(l.items, Key(backCacheItem.key))
			l.queue.Remove(item)
		}
	}

	// Update map value anyway
	l.items[key] = listItem

	return exists
}

// LruCache getter: returns value if exists, or nil, if doesnt.
func (l *lruCache) Get(key Key) (interface{}, bool) {
	item, exists := l.items[key]

	// If cache element doesn't exist, return nil, false
	if !exists {
		return nil, exists
	}

	// If cache element exists, moves it to front
	l.queue.MoveToFront(item)

	// To get actual value, interface{} needs to be casted to cacheItem
	cacheItemElement := item.Value.(cacheItem)
	return cacheItemElement.value, exists
}

// Reinitialize lruCache instance.
func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   string
	value interface{}
}

// Cache constructor: returns lruCache instance pointer.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
