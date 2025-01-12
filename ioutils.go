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

func ReadAll(r io.ReadCloser) ([]byte, error) {
	b, err := io.ReadAll(r)
	err = errors.Wrap(err, r.Close())
	if err != nil {
		return nil, err
	}
	return b, nil
}

const copyBufferSize = 512

func Copy(w io.Writer, r io.Reader) (int, error) {
	return CopyN(w, r, copyBufferSize)
}

func CopyN(w io.Writer, r io.Reader, size int) (int, error) {
	n := 0
	buff := make([]byte, size)
	for {
		m, err1 := r.Read(buff)
		if m < 0 {
			return 0, fmt.Errorf("reader err: negative read bytes")
		}
		if err1 != nil && !errors.Is(err1, io.EOF) {
			return 0, err1
		}
		n += m
		_, err2 := w.Write(buff[:m])
		if err2 != nil {
			return 0, fmt.Errorf("writer err: %w", err2)
		}
		if m < size || errors.Is(err1, io.EOF) {
			return n, nil
		}
	}
}
