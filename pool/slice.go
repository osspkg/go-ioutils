/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package pool

func NewSlicePool[T any](l, c int) *Pool[*SlicePool[T]] {
	return New(func() *SlicePool[T] {
		return &SlicePool[T]{B: make([]T, l, c)}
	})
}
