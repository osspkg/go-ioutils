/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cache

import (
	"context"
	"time"

	"go.osspkg.com/routine"
)

func OptTimeClean[K comparable, V Timestamp](ctx context.Context, interval time.Duration) Option[K, V] {
	return func(v *_cache[K, V]) {

		tik := routine.Ticker{
			Interval: interval,
			OnStart:  false,
			Calls: []routine.TickFunc{
				func(ctx context.Context, t time.Time) {
					curr := t.Unix()
					keys := make([]K, 0, 10)

					v.mux.RLock()
					for key, value := range v.list {
						if value.Timestamp() < curr {
							keys = append(keys, key)
						}
					}
					v.mux.RUnlock()

					if len(keys) == 0 {
						return
					}

					v.mux.Lock()
					defer v.mux.Unlock()

					for _, key := range keys {
						delete(v.list, key)
					}
				},
			},
		}

		tik.Run(ctx)
	}
}

func OptCountRandomClean[K comparable, V any](ctx context.Context, maxCount int, interval time.Duration) Option[K, V] {
	return func(v *_cache[K, V]) {

		if maxCount < 0 {
			panic("OptCountRandomClean: maxCount < 0")
		}

		tik := routine.Ticker{
			Interval: interval,
			OnStart:  false,
			Calls: []routine.TickFunc{
				func(ctx context.Context, _ time.Time) {

					v.mux.RLock()
					removeCount := len(v.list) - maxCount
					v.mux.RUnlock()

					if removeCount <= 0 {
						return
					}

					v.mux.Lock()
					defer v.mux.Unlock()

					for key := range v.list {
						if removeCount <= 0 {
							return
						}

						delete(v.list, key)
						removeCount--
					}
				},
			},
		}

		tik.Run(ctx)
	}
}
