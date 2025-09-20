/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cache

type Cache[K comparable, V any] interface {
	Has(key K) bool
	Get(key K) (V, bool)
	Extract(key K) (V, bool)
	Set(key K, value V)
	Replace(data map[K]V)
	Del(key K)
	Keys() []K
	Size() int
	Flush()
}

type Option[K comparable, V any] func(*_cache[K, V])

type Timestamp interface {
	Timestamp() int64
}
