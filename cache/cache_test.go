/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cache_test

import (
	"context"
	"testing"
	"time"

	"go.osspkg.com/casecheck"

	"go.osspkg.com/ioutils/cache"
)

func TestUnit_New(t *testing.T) {
	c := cache.New[string, string]()

	c.Set("foo", "bar")
	casecheck.True(t, c.Has("foo"))

	casecheck.Equal(t, []string{"foo"}, c.Keys())

	v, ok := c.Get("foo")
	casecheck.True(t, ok)
	casecheck.Equal(t, v, "bar")

	v, ok = c.Extract("foo")
	casecheck.True(t, ok)
	casecheck.Equal(t, v, "bar")

	v, ok = c.Extract("foo")
	casecheck.False(t, ok)
	casecheck.Equal(t, v, "")

	v, ok = c.Get("foo")
	casecheck.False(t, ok)
	casecheck.Equal(t, v, "")

	casecheck.False(t, c.Has("foo"))
	casecheck.Equal(t, []string{}, c.Keys())

	c.Set("foo", "bar")
	casecheck.True(t, c.Has("foo"))

	c.Del("foo")
	casecheck.False(t, c.Has("foo"))

	c.Replace(map[string]string{"foo": "bar"})
	casecheck.Equal(t, []string{"foo"}, c.Keys())

	c.Flush()
	casecheck.Equal(t, []string{}, c.Keys())
}

type testValue struct {
	Val string
	TS  int64
}

func (v testValue) Timestamp() int64 { return v.TS }

func TestUnit_AutoClean(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := cache.New[string, testValue](
		cache.AutoClean[string, testValue](ctx, time.Millisecond*100),
	)

	c.Set("foo", testValue{Val: "bar", TS: time.Now().Add(time.Millisecond * 200).Unix()})
	casecheck.True(t, c.Has("foo"))

	time.Sleep(time.Second)

	casecheck.False(t, c.Has("foo"))
}

func Benchmark_New(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := cache.New[string, testValue](
		cache.AutoClean[string, testValue](ctx, time.Millisecond*100),
	)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

			c.Set("foo", testValue{Val: "bar", TS: time.Now().Add(time.Millisecond * 200).Unix()})
			c.Get("foo")
			c.Has("foo")
			c.Extract("foo")
			c.Replace(map[string]testValue{"foo": {Val: "bar", TS: time.Now().Add(time.Millisecond * 200).Unix()}})
			c.Keys()
			c.Del("foo")
			c.Set("foo", testValue{Val: "bar", TS: time.Now().Add(time.Millisecond * 200).Unix()})
			c.Flush()
		}
	})
}
