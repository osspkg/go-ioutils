package cache

import (
	"context"
	"sync"
	"time"

	"go.osspkg.com/routine"
)

type (
	withTTL[K comparable, T interface{}] struct {
		ttl  time.Duration
		list map[K]*ttlItem[T]
		mux  sync.RWMutex
	}

	ttlItem[T interface{}] struct {
		link T
		ts   int64
	}
)

func NewWithTTL[K comparable, T interface{}](ctx context.Context, ttl time.Duration) TCache[K, T] {
	cache := &withTTL[K, T]{
		ttl:  ttl,
		list: make(map[K]*ttlItem[T], 1000),
	}
	go cache.cleaner(ctx)
	return cache
}

func (v *withTTL[K, T]) cleaner(ctx context.Context) {
	routine.Interval(ctx, v.ttl, func(ctx context.Context) {
		curr := time.Now().Unix()

		for k, t := range v.list {
			if t.ts < curr {
				delete(v.list, k)
			}
		}
	})
}

func (v *withTTL[K, T]) Has(key K) bool {
	v.mux.RLock()
	defer v.mux.RUnlock()

	_, ok := v.list[key]

	return ok
}

func (v *withTTL[K, T]) Get(key K) (T, bool) {
	v.mux.RLock()
	defer v.mux.RUnlock()

	item, ok := v.list[key]
	if !ok {
		var zeroValue T
		return zeroValue, false
	}

	return item.link, true
}

func (v *withTTL[K, T]) Set(key K, value T) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.list[key] = &ttlItem[T]{
		link: value,
		ts:   time.Now().Add(v.ttl).Unix(),
	}
}

func (v *withTTL[K, T]) SetWithTTL(key K, value T, ttl time.Time) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.list[key] = &ttlItem[T]{
		link: value,
		ts:   ttl.Unix(),
	}
}

func (v *withTTL[K, T]) Del(key K) {
	v.mux.Lock()
	defer v.mux.Unlock()

	delete(v.list, key)
}

func (v *withTTL[K, T]) Keys() []K {
	v.mux.RLock()
	defer v.mux.RUnlock()

	result := make([]K, 0, len(v.list))
	for k := range v.list {
		result = append(result, k)
	}

	return result
}

func (v *withTTL[K, T]) Flush() {
	v.mux.Lock()
	defer v.mux.Unlock()

	for k := range v.list {
		delete(v.list, k)
	}
}
