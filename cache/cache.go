/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cache

import (
	"iter"
	"sync"

	"go.osspkg.com/random"
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

func (v *_cache[K, V]) One() (key K, val V, ok bool) {
	keys := v._keys(30)
	if len(keys) == 0 {
		return
	}

	random.Shuffle(keys)

	key = keys[0]
	val, ok = v.Get(key)

	return
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
	return v._keys(v.Size())
}

func (v *_cache[K, V]) _keys(limit int) []K {
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

func (v *_cache[K, V]) Flush() {
	v.mux.Lock()
	defer v.mux.Unlock()

	for k := range v.list {
		delete(v.list, k)
	}
}

func (v *_cache[K, V]) Yield(limit int) iter.Seq2[K, V] {
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
