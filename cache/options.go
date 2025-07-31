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

func AutoClean[K comparable, V Timestamp](ctx context.Context, interval time.Duration) Option[K, V] {
	return func(v *_cache[K, V]) {
		routine.Interval(ctx, interval, func(ctx context.Context) {
			curr := time.Now().Unix()
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
		})
	}
}
