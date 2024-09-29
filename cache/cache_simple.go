package cache

import (
	"sync"
)

type (
	cacheReplace[K comparable, V any] struct {
		list map[K]V
		mux  sync.RWMutex
	}
)

func NewWithReplace[K comparable, V any]() TCacheReplace[K, V] {
	return &cacheReplace[K, V]{
		list: make(map[K]V, 1000),
	}
}

func (v *cacheReplace[K, V]) Has(key K) bool {
	v.mux.RLock()
	defer v.mux.RUnlock()

	_, ok := v.list[key]

	return ok
}

func (v *cacheReplace[K, V]) Get(key K) (V, bool) {
	v.mux.RLock()
	defer v.mux.RUnlock()

	item, ok := v.list[key]
	if !ok {
		var zeroValue V
		return zeroValue, false
	}

	return item, true
}

func (v *cacheReplace[K, V]) Set(key K, value V) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.list[key] = value
}

func (v *cacheReplace[K, V]) Replace(data map[K]V) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.list = data
}

func (v *cacheReplace[K, V]) Del(key K) {
	v.mux.Lock()
	defer v.mux.Unlock()

	delete(v.list, key)
}

func (v *cacheReplace[K, V]) Keys() []K {
	v.mux.RLock()
	defer v.mux.RUnlock()

	result := make([]K, 0, len(v.list))
	for k := range v.list {
		result = append(result, k)
	}

	return result
}

func (v *cacheReplace[K, V]) Flush() {
	v.mux.Lock()
	defer v.mux.Unlock()

	for k := range v.list {
		delete(v.list, k)
	}
}
