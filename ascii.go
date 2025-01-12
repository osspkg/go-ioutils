/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ioutils

import "bytes"

var asciiEOF = []byte{255, 244, 255, 253, 6}

func IsAsciiEOF(b []byte) bool {
	return bytes.Equal(b, asciiEOF)
}
