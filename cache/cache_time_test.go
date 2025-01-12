/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cache

import (
	"context"
	"testing"
	"time"

	"go.osspkg.com/casecheck"
)

func TestUnit_WithTTL_Pointer(t *testing.T) {
	type A struct {
		Data uint64
	}

	ctx, cncl := context.WithTimeout(context.TODO(), time.Second)
	defer cncl()
	c := NewWithTTL[string, *A](ctx, time.Minute)
	v, ok := c.Get("a")
	casecheck.False(t, ok)
	casecheck.Nil(t, v)
}
