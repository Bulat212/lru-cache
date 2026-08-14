package main

import (
	"container/list"
	"fmt"
	"sync"
)

type Element[k comparable, v any] struct {
	key   k
	value v
}

type Cache[k comparable, v any] struct {
	store    map[k]*list.Element
	capacity int
	list     *list.List
	mtx      sync.Mutex
}

func NewCache[k comparable, v any](capacity int) *Cache[k, v] {
	return &Cache[k, v]{
		store:    make(map[k]*list.Element),
		capacity: capacity,
		list:     list.New(),
		mtx:      sync.Mutex{},
	}
}
	
func (c *Cache[k, v]) Set(key k, value v) {

	c.mtx.Lock()
	defer c.mtx.Unlock()

	if el, exists := c.store[key]; exists {
		el.Value = Element[k, v]{key: key, value: value}
		c.list.MoveToFront(el)
		return
	}

	if c.list.Len() == c.capacity {
		removedValue := c.list.Remove(c.list.Back())

		removedElement, ok := removedValue.(Element[k, v])
		if !ok {
			fmt.Println("В интерфейсе лежит не Element!")
			return
		}
		delete(c.store, removedElement.key)
	}

	newElement := c.list.PushFront(Element[k, v]{key: key, value: value})
	c.store[key] = newElement
}

func (c *Cache[k, v]) Get(key k) (v, bool) {

	c.mtx.Lock()
	defer c.mtx.Unlock()

	if el, exists := c.store[key]; exists {
		c.list.MoveToFront(el)
		valueElement := el.Value.(Element[k, v])
		return valueElement.value, true
	}
	var zeroValue v

	return zeroValue, false
}

func (c *Cache[k, v]) Clear() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.list.Init()
	clear(c.store)
}
