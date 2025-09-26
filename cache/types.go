/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cache

import "iter"

type Cache[K comparable, V any] interface {
	Has(K) bool
	Get(K) (V, bool)
	//Yield if limit <=0 then got all elements for range
	Yield(limit int) iter.Seq2[K, V]
	//Extract getting element and delete form cache
	Extract(K) (V, bool)
	//One getting one random key-value element
	One() (K, V, bool)
	Set(K, V)
	//Replace replace all elements
	Replace(map[K]V)
	Del(K)
	Keys() []K
	Size() int
	Flush()
}

type Option[K comparable, V any] func(*_cache[K, V])

type Timestamp interface {
	Timestamp() int64
}
