/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ioutils

import "sync"

type TPool interface {
	Reset()
}

type Pool[T TPool] struct {
	callNew func() T
	pool    sync.Pool
}

func NewPool[T TPool](callNew func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{New: func() any { return callNew() }},
	}
}

func (v *Pool[T]) Get() T {
	buf, ok := v.pool.Get().(T)
	if !ok {
		buf = v.callNew()
	}
	return buf
}

func (v *Pool[T]) Put(t T) {
	t.Reset()
	v.pool.Put(t)
}

type SlicePool[T any] struct {
	B []T
}

func (v *SlicePool[T]) Reset() {
	v.B = v.B[:0]
}

func NewSlicePool[T any](l, c int) *Pool[*SlicePool[T]] {
	return NewPool(func() *SlicePool[T] {
		return &SlicePool[T]{B: make([]T, l, c)}
	})
}
