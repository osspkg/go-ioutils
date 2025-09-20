/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cache

import (
	"sync"
)

type (
	_cache[K comparable, V any] struct {
		list map[K]V
		mux  sync.RWMutex
	}
)

func New[K comparable, V any](opts ...Option[K, V]) Cache[K, V] {
	obj := &_cache[K, V]{
		list: make(map[K]V, 100),
	}

	for _, opt := range opts {
		go opt(obj)
	}

	return obj
}

func (v *_cache[K, V]) Size() int {
	v.mux.RLock()
	defer v.mux.RUnlock()

	return len(v.list)
}

func (v *_cache[K, V]) Has(key K) bool {
	v.mux.RLock()
	defer v.mux.RUnlock()

	_, ok := v.list[key]

	return ok
}

func (v *_cache[K, V]) Get(key K) (V, bool) {
	v.mux.RLock()
	defer v.mux.RUnlock()

	item, ok := v.list[key]
	if !ok {
		var zeroValue V
		return zeroValue, false
	}

	return item, true
}

func (v *_cache[K, V]) Extract(key K) (V, bool) {
	v.mux.Lock()
	defer v.mux.Unlock()

	item, ok := v.list[key]
	if !ok {
		var zeroValue V
		return zeroValue, false
	}

	delete(v.list, key)

	return item, true
}

func (v *_cache[K, V]) Set(key K, value V) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.list[key] = value
}

func (v *_cache[K, V]) Replace(data map[K]V) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.list = data
}

func (v *_cache[K, V]) Del(key K) {
	v.mux.Lock()
	defer v.mux.Unlock()

	delete(v.list, key)
}

func (v *_cache[K, V]) Keys() []K {
	v.mux.RLock()
	defer v.mux.RUnlock()

	result := make([]K, 0, len(v.list))
	for k := range v.list {
		result = append(result, k)
	}

	return result
}

func (v *_cache[K, V]) Flush() {
	v.mux.Lock()
	defer v.mux.Unlock()

	for k := range v.list {
		delete(v.list, k)
	}
}
