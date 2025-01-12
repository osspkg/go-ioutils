/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cache

import (
	"context"
	"sync"
	"time"

	"go.osspkg.com/routine"
)

type (
	cacheWithTTL[K comparable, V any] struct {
		ttl  time.Duration
		list map[K]*itemCacheTTL[V]
		mux  sync.RWMutex
	}

	itemCacheTTL[V interface{}] struct {
		link V
		ts   int64
	}
)

func NewWithTTL[K comparable, V any](ctx context.Context, ttl time.Duration) TCacheTTL[K, V] {
	cache := &cacheWithTTL[K, V]{
		ttl:  ttl,
		list: make(map[K]*itemCacheTTL[V], 1000),
	}
	go cache.cleaner(ctx)
	return cache
}

func (v *cacheWithTTL[K, V]) cleaner(ctx context.Context) {
	routine.Interval(ctx, v.ttl, func(ctx context.Context) {
		curr := time.Now().Unix()

		for k, t := range v.list {
			if t.ts < curr {
				delete(v.list, k)
			}
		}
	})
}

func (v *cacheWithTTL[K, V]) Has(key K) bool {
	v.mux.RLock()
	defer v.mux.RUnlock()

	_, ok := v.list[key]

	return ok
}

func (v *cacheWithTTL[K, V]) Get(key K) (V, bool) {
	v.mux.RLock()
	defer v.mux.RUnlock()

	item, ok := v.list[key]
	if !ok {
		var zeroValue V
		return zeroValue, false
	}

	return item.link, true
}

func (v *cacheWithTTL[K, V]) Set(key K, value V) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.list[key] = &itemCacheTTL[V]{
		link: value,
		ts:   time.Now().Add(v.ttl).Unix(),
	}
}

func (v *cacheWithTTL[K, V]) SetWithTTL(key K, value V, ttl time.Time) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.list[key] = &itemCacheTTL[V]{
		link: value,
		ts:   ttl.Unix(),
	}
}

func (v *cacheWithTTL[K, V]) Del(key K) {
	v.mux.Lock()
	defer v.mux.Unlock()

	delete(v.list, key)
}

func (v *cacheWithTTL[K, V]) Keys() []K {
	v.mux.RLock()
	defer v.mux.RUnlock()

	result := make([]K, 0, len(v.list))
	for k := range v.list {
		result = append(result, k)
	}

	return result
}

func (v *cacheWithTTL[K, V]) Flush() {
	v.mux.Lock()
	defer v.mux.Unlock()

	for k := range v.list {
		delete(v.list, k)
	}
}
