/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ioutils

import "testing"

func TestIsAsciiEOF(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{name: "matching marker", data: []byte{255, 244, 255, 253, 6}, want: true},
		{name: "different marker", data: []byte{255, 244, 255, 253, 7}},
		{name: "short input", data: []byte{255, 244, 255}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAsciiEOF(tt.data); got != tt.want {
				t.Fatalf("IsAsciiEOF(%v) = %t, want %t", tt.data, got, tt.want)
			}
		})
	}
}
