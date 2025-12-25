/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cache

import (
	"context"
	"time"

	"go.osspkg.com/routine/tick"
)

func OptTimeClean[K comparable, V Timestamp](ctx context.Context, interval time.Duration) Option[K, V] {
	return func(v *_cache[K, V]) {

		tik := tick.Ticker{
			Calls: []tick.Config{{
				Interval: interval,
				Func: func(ctx context.Context, t time.Time) error {
					curr := t.Unix()

					v.mux.Lock()
					defer v.mux.Unlock()

					for key, value := range v.list {
						if value.Timestamp() < curr {
							delete(v.list, key)
						}
					}

					return nil
				},
			}},
		}

		tik.Run(ctx)
	}
}

func OptCountRandomClean[K comparable, V any](ctx context.Context, maxCount int, interval time.Duration) Option[K, V] {
	return func(v *_cache[K, V]) {

		if maxCount < 0 {
			panic("OptCountRandomClean: maxCount < 0")
		}

		tik := tick.Ticker{
			Calls: []tick.Config{{
				Interval: interval,
				Func: func(ctx context.Context, t time.Time) error {
					removeCount := v.Size() - maxCount
					if removeCount <= 0 {
						return nil
					}

					for key := range v.Yield(removeCount) {
						v.Del(key)
					}

					return nil
				},
			}},
		}

		tik.Run(ctx)
	}
}
