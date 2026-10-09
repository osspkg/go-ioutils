/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ioutils

import (
	"io"

	"go.osspkg.com/errors"
)

const packSize = 512

const maxConsecutiveEmptyReads = 100

func Copy(w io.Writer, r io.Reader) (int, error) {
	return CopyN(w, r, packSize)
}

func CopyN(w io.Writer, r io.Reader, size int) (int, error) {
	if size <= 0 {
		return 0, errors.New("size must be greater than zero")
	}

	buf := make([]byte, size)

	return CopyB(w, r, buf)
}

func CopyB(w io.Writer, r io.Reader, buff []byte) (int, error) {
	if len(buff) <= 0 {
		return 0, errors.New("size must be greater than zero")
	}
	if w == nil {
		return 0, errors.New("writer must not be nil")
	}
	if r == nil {
		return 0, errors.New("reader must not be nil")
	}

	return copyBuffer(w, r, buff)
}

func copyBuffer(w io.Writer, r io.Reader, buff []byte) (int, error) {
	total := 0
	emptyReads := 0

	for {
		rn, readErr := r.Read(buff)
		if rn < 0 {
			return total, errors.New("reader err: negative read bytes")
		}

		if rn > 0 {
			wn, writeErr := w.Write(buff[:rn])
			total += wn
			if writeErr != nil {
				return total, errors.Wrapf(writeErr, "writer err")
			}
			if wn != rn {
				return total, io.ErrShortWrite
			}
			emptyReads = 0
		} else if readErr == nil {
			emptyReads++
			if emptyReads >= maxConsecutiveEmptyReads {
				return total, io.ErrNoProgress
			}
		}

		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				return total, nil
			}
			return total, readErr
		}
	}
}
