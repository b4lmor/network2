package utils

import (
	"sync"
)

type SafeMap[K comparable, V any] struct {
	mu    sync.Mutex
	store map[K]V
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		store: make(map[K]V),
	}
}

func (sm *SafeMap[K, V]) Set(key K, value V) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.store[key] = value
}

func (sm *SafeMap[K, V]) Get(key K) (V, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	value, exists := sm.store[key]
	return value, exists
}

func (sm *SafeMap[K, V]) Len() int {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return len(sm.store)
}

func (sm *SafeMap[K, V]) GetAllKeys() []K {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	keys := make([]K, 0, len(sm.store))
	for key := range sm.store {
		keys = append(keys, key)
	}
	return keys
}

func (sm *SafeMap[K, V]) Delete(key K) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.store, key)
}
