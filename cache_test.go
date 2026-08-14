package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestDeleteOldElement(t *testing.T) {

	cache := NewCache[int, string](2)
	cache.Set(1, "apple")
	cache.Set(2, "banana")
	cache.Set(3, "potato")

	if _, exists := cache.Get(1); exists {
		t.Error("Ожидалось, что элемент под ключом 1 со значением удалится")
	}

	if _, exists := cache.Get(3); !exists {
		t.Error("Ожидалось, что элемент под ключом 3 добавился")
	}

}

func TestOverflowCache(t *testing.T) {

	cache := NewCache[int, string](3) 
	cache.Set(1, "apple")
	cache.Set(2, "banana")
	cache.Set(3, "potato")  

	cache.Get(1)     
	cache.Set(4, "mandarin")

	if _, exists := cache.Get(2); exists {
		t.Error("Ожидалось, что элемент под ключом 2 со значением удалится")
	}

	if _, exists := cache.Get(4); !exists {
		t.Error("Ожидалось, что элемент под ключом 4 добавился")
	}
	if _, exists := cache.Get(1); !exists {
		t.Error("Ожидалось, что элемент под ключом 1 существует")
	}
	if _, exists := cache.Get(3); !exists {
		t.Error("Ожидалось, что элемент под ключом 3 существует")
	}
}

func TestСoncurrency(t *testing.T) {

	cache := NewCache[int, string](3)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 20; i++ {
			cache.Set(i, fmt.Sprintf("element %d", i))
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 20; i++ {
			cache.Set(i, fmt.Sprintf("element %d", i))
		}
	}()

	wg.Wait()

	if cache.list.Len() > cache.capacity {
		t.Errorf("список превысил ёмкость: len=%d, capacity=%d", cache.list.Len(), cache.capacity)
	}

	if len(cache.store) != cache.list.Len() {
		t.Errorf("store и list рассинхронизированы: len(store)=%d, list.Len()=%d", len(cache.store), cache.list.Len())
	}
}

func TestClear(t *testing.T) {
	cache := NewCache[int, string](3)

	cache.Set(1, "apple")
	cache.Set(2, "banana")
	cache.Set(3, "potato")

	cache.Clear()

	if _, exists := cache.Get(1); exists {
		t.Error("Ожидалось, что элемент под ключом 1 будет удалён после Clear")
	}
	if _, exists := cache.Get(2); exists {
		t.Error("Ожидалось, что элемент под ключом 2 будет удалён после Clear")
	}
	if _, exists := cache.Get(3); exists {
		t.Error("Ожидалось, что элемент под ключом 3 будет удалён после Clear")
	}

	if cache.list.Len() != 0 {
		t.Errorf("Ожидалась пустая длина списка, получено: %d", cache.list.Len())
	}
	if len(cache.store) != 0 {
		t.Errorf("Ожидалась пустая map, получено: %d", len(cache.store))
	}

	cache.Set(4, "mandarin")
	if val, exists := cache.Get(4); !exists || val != "mandarin" {
		t.Errorf("Кэш должен работать после Clear, получено: val=%v exists=%v", val, exists)
	}
}