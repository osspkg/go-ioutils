/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ioutils

import (
	"fmt"
	"io"

	"go.osspkg.com/errors"
)

func Pipe(w io.Writer, r io.Reader, size int) (n int, err error) {
	if size <= 0 {
		return 0, fmt.Errorf("invalid buffer size")
	}

	buf := make([]byte, size)

	for {
		rn, re := r.Read(buf)
		if rn > 0 {
			wn, we := w.Write(buf[:rn])
			if rn != wn {
				wn = 0
				if we == nil {
					we = io.ErrShortWrite
				}
			}
			n += wn
			if we != nil {
				if !errors.Is(we, io.EOF) {
					err = we
				}
				break
			}
		}
		if re != nil {
			if !errors.Is(re, io.EOF) {
				err = re
			}
			break
		}
	}

	return
}
