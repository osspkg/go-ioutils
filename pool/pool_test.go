/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package pool

import (
	"testing"

	"go.osspkg.com/casecheck"
)

func TestUnit_Pool(t *testing.T) {
	buf := NewSlicePool[byte](2, 10)

	item := buf.Get()
	casecheck.True(t, len(item.B) == 2)
	casecheck.True(t, cap(item.B) == 10)
	buf.Put(item)
}

func Benchmark_Pool(b *testing.B) {
	buf := NewSlicePool[byte](0, 10)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			item := buf.Get()
			if len(item.B) != 0 {
				b.FailNow()
			}
			item.B = append(item.B, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}...)
			buf.Put(item)
		}
	})

}
