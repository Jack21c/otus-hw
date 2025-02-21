package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type CacheItem struct {
	Key   Key
	Value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	cacheValue := CacheItem{
		Key:   key,
		Value: value,
	}
	if item, exist := c.items[key]; exist {
		item.Value = cacheValue
		c.queue.MoveToFront(item)
		return exist
	}
	if len(c.items) >= c.capacity {
		vBack, ok := c.queue.Back().Value.(CacheItem)
		if !ok {
			for k, v := range c.items {
				if v == c.queue.Back() {
					delete(c.items, k)
					break
				}
			}
		} else {
			delete(c.items, vBack.Key)
		}
		c.queue.Remove(c.queue.Back())
	}
	item := c.queue.PushFront(cacheValue)
	c.items[key] = item

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, exist := c.items[key]; exist {
		c.queue.MoveToFront(item)
		if value, ok := item.Value.(CacheItem); ok {
			return value.Value, true
		}
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	clear(c.items)
}
