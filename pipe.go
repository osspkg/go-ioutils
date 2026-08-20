/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ioutils

import (
	"io"

	"go.osspkg.com/errors"
)

func Pipe(w io.Writer, r io.Reader, size int) (int, error) {
	if size <= 0 {
		return 0, errors.New("size must be greater than zero")
	}
	if w == nil {
		return 0, errors.New("writer must not be nil")
	}
	if r == nil {
		return 0, errors.New("reader must not be nil")
	}

	buff := make([]byte, size)

	n := 0

	for {
		rn, re := r.Read(buff)
		if rn < 0 {
			return n, errors.New("reader err: negative read bytes")
		}

		if rn > 0 {
			wn, we := w.Write(buff[:rn])
			if we != nil {
				return n, errors.Wrapf(we, "writer err")
			}

			n += wn

			if re != nil {
				if errors.Is(re, io.EOF) {
					return n, nil
				}
				return n, re
			}

			if wn != rn {
				return n, io.ErrShortWrite
			}
		}

		if re != nil {
			if errors.Is(re, io.EOF) {
				return n, nil
			}
			return n, re
		}
	}
}
