/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cache

import (
	"iter"
	"sync"

	"go.osspkg.com/random"
)

type (
	store[K comparable, V any] struct {
		list map[K]V
		mux  sync.RWMutex
	}
)

func New[K comparable, V any](opts ...Option[K, V]) Cache[K, V] {
	obj := &store[K, V]{
		list: make(map[K]V, 100),
	}

	for _, opt := range opts {
		go opt(obj)
	}

	return obj
}

func (v *store[K, V]) Size() int {
	v.mux.RLock()
	defer v.mux.RUnlock()

	return len(v.list)
}

func (v *store[K, V]) Has(key K) bool {
	v.mux.RLock()
	defer v.mux.RUnlock()

	_, ok := v.list[key]

	return ok
}

func (v *store[K, V]) Get(key K) (V, bool) {
	v.mux.RLock()
	defer v.mux.RUnlock()

	item, ok := v.list[key]
	if !ok {
		var zeroValue V
		return zeroValue, false
	}

	return item, true
}

func (v *store[K, V]) One() (key K, val V, ok bool) {
	keys := v._keys(30)
	if len(keys) == 0 {
		return
	}

	random.Shuffle(keys)

	key = keys[0]
	val, ok = v.Get(key)

	return
}

func (v *store[K, V]) Extract(key K) (V, bool) {
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

func (v *store[K, V]) Set(key K, value V) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.list[key] = value
}

func (v *store[K, V]) Replace(data map[K]V) {
	v.mux.Lock()
	defer v.mux.Unlock()

	replacement := make(map[K]V, len(data))
	for key, value := range data {
		replacement[key] = value
	}
	v.list = replacement
}

func (v *store[K, V]) Del(key K) {
	v.mux.Lock()
	defer v.mux.Unlock()

	delete(v.list, key)
}

func (v *store[K, V]) Keys() []K {
	return v._keys(v.Size())
}

func (v *store[K, V]) _keys(limit int) []K {
	v.mux.RLock()
	defer v.mux.RUnlock()

	i := 0
	result := make([]K, 0, limit)
	for k := range v.list {
		result = append(result, k)
		i++
		if i >= limit {
			break
		}
	}

	return result
}

func (v *store[K, V]) Flush() {
	v.mux.Lock()
	defer v.mux.Unlock()

	clear(v.list)
}

func (v *store[K, V]) Yield(limit int) iter.Seq2[K, V] {
	if limit < 1 {
		limit = v.Size()
	}

	keys := v._keys(limit)

	return func(yield func(K, V) bool) {
		for _, key := range keys {
			if val, ok := v.Get(key); ok {
				if !yield(key, val) {
					return
				}
			}
		}
	}
}
