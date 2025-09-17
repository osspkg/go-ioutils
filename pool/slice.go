/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package pool

type Slice[T any] struct {
	B []T
}

func (v *Slice[T]) Reset() {
	v.B = v.B[:0]
}

func NewSlicePool[T any](l, c int) *Pool[*Slice[T]] {
	return New(func() *Slice[T] {
		return &Slice[T]{B: make([]T, l, c)}
	})
}
