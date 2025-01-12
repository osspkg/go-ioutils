/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cache

import "time"

type TCache[K comparable, V interface{}] interface {
	Has(key K) bool
	Get(key K) (V, bool)
	Set(key K, value V)
	Del(key K)
	Keys() []K
	Flush()
}

type TCacheTTL[K comparable, V interface{}] interface {
	TCache[K, V]
	SetWithTTL(key K, value V, ttl time.Time)
}

type TCacheReplace[K comparable, V interface{}] interface {
	TCache[K, V]
	Replace(data map[K]V)
}
