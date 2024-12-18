/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ioutils

import (
	"bytes"
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

const buffSize = 128

var (
	ErrMaximumSize = errors.New("maximum buffer size reached")
	ErrInvalidSize = errors.New("invalid size")
)

func ReadFull(w io.Writer, r io.Reader, maxSize int) error {
	if maxSize < 0 {
		return ErrInvalidSize
	}

	// nolint: ineffassign
	var (
		total = 0
		n     = 0
		buff  = make([]byte, buffSize)
		err   error
	)

	for {
		n, err = r.Read(buff)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if n < 0 {
			return ErrInvalidSize
		}
		if _, err = w.Write(buff[:n]); err != nil {
			return err
		}
		total += n
		if maxSize > 0 && total > maxSize {
			return ErrMaximumSize
		}
		if n < buffSize {
			break
		}
	}
	return nil
}

func ReadBytes(v io.Reader, divide string) ([]byte, error) {
	var (
		n   int
		err error
		b   = make([]byte, 0, 512)
		db  = []byte(divide)
		dl  = len(db)
	)

	for {
		if len(b) == cap(b) {
			b = append(b, 0)[:len(b)]
		}
		n, err = v.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if len(b) < dl {
			return b, io.EOF
		}
		if bytes.Equal(db, b[len(b)-dl:]) {
			b = b[:len(b)-dl]
			break
		}
	}
	return b, nil
}

func WriteBytes(v io.Writer, b []byte, divide string) error {
	var (
		db = []byte(divide)
		dl = len(db)
	)
	if len(b) < dl || !bytes.Equal(db, b[len(b)-dl:]) {
		b = append(b, db...)
	}
	if _, err := v.Write(b); err != nil {
		return err
	}
	return nil
}

const copyBufferSize = 512

func Copy(w io.Writer, r io.Reader) (int, error) {
	return CopyPack(w, r, copyBufferSize)
}

func CopyPack(w io.Writer, r io.Reader, size int) (int, error) {
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
