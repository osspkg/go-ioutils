/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package data

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

const packSize = 65535

type Buffer struct {
	buf []byte
	pos int
}

func NewBuffer(size int) *Buffer {
	if size < 0 {
		size = 0
	}
	return &Buffer{
		buf: make([]byte, 0, size),
		pos: 0,
	}
}

func (v *Buffer) Reset() {
	l := len(v.buf)
	if l > 5e6 {
		v.buf = v.buf[:0:5e6]
	} else {
		v.buf = v.buf[:0]
	}
	v.pos = 0
}

func (v *Buffer) Bytes() []byte {
	return v.buf[:]
}

func (v *Buffer) String() string {
	return string(v.Bytes())
}

func (v *Buffer) Size() int {
	return len(v.buf)
}

func (v *Buffer) Len() int {
	return len(v.buf) - v.pos
}

func (v *Buffer) Truncate(c int) {
	if c <= 0 {
		return
	}
	n := v.Size() - c
	if n <= 0 {
		v.Reset()
		return
	}

	v.buf = v.buf[:n]

	for i := 1; i <= 3 && n-i >= 0; i++ {
		off := n - i
		if utf8.RuneStart(v.buf[off]) {
			if v.buf[off] > utf8.RuneSelf {
				n -= i
			}
			break
		}
	}

	v.buf = v.buf[:n]

	if v.pos > n {
		v.pos = n
	}
}

func (v *Buffer) Write(p []byte) (int, error) {
	v.buf = append(v.buf, p...)
	return len(p), nil
}

func (v *Buffer) WriteString(s string) (int, error) {
	return v.Write([]byte(s))
}

func (v *Buffer) WriteByte(b byte) error {
	v.buf = append(v.buf, b)
	return nil
}

func (v *Buffer) WriteRune(r rune) (n int, err error) {
	n = v.Size()
	v.buf = utf8.AppendRune(v.buf, r)
	n = v.Size() - n
	return
}

func (v *Buffer) WriteTo(w io.Writer) (n int64, err error) {
	return v.WriteToN(w, packSize)
}

func (v *Buffer) WriteToN(w io.Writer, size int) (n int64, err error) {
	if v.Len() <= 0 {
		return 0, io.EOF
	}

	b := make([]byte, size)

	for {
		rn, re := v.Read(b)
		if rn > 0 {
			wn, we := w.Write(b[:rn])
			if wn > rn || wn <= 0 {
				v.Seek(int64(-rn), SeekCurr) //nolint: errcheck
				err = fmt.Errorf("invalid write result")
				break
			}
			if wn != rn {
				v.Seek(int64(-rn+wn), SeekCurr) //nolint: errcheck
			}
			n += int64(wn)
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

func (v *Buffer) WriteAt(b []byte, off int64) (int, error) {
	if off < 0 {
		off = 0
	}

	if add := len(b) + int(off) - v.Size(); add > 0 {
		v.buf = append(v.buf, make([]byte, add)...)
	}

	copy(v.buf[off:], b[:])

	return len(b), nil
}

func (v *Buffer) ReadFrom(r io.Reader) (int64, error) {
	return v.ReadFromN(r, packSize)
}

func (v *Buffer) ReadFromN(r io.Reader, size int) (int64, error) {
	n := 0
	b := make([]byte, size)
	for {
		m, err := r.Read(b)
		if m < 0 {
			return 0, fmt.Errorf("negative read bytes")
		}
		if err != nil && !errors.Is(err, io.EOF) {
			return 0, err
		}
		n += m
		v.buf = append(v.buf, b[:m]...)
		if m < size || errors.Is(err, io.EOF) {
			break
		}
	}
	return int64(n), nil
}

func (v *Buffer) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, fmt.Errorf("got zero buffer arg")
	}
	if v.Len() == 0 {
		return 0, io.EOF
	}
	n := copy(p[:], v.buf[v.pos:])
	v.pos += n
	return n, nil
}

func (v *Buffer) ReadAt(p []byte, off int64) (int, error) {
	if len(p) == 0 {
		return 0, fmt.Errorf("got zero buffer arg")
	}
	if off < 0 || int(off) >= v.Size() {
		return 0, io.EOF
	}
	n := copy(p[:], v.buf[int(off):])
	return n, nil
}

func (v *Buffer) Discard(n int) int {
	if n <= 0 {
		return 0
	}
	np, _ := v.Seek(int64(n), SeekCurr) //nolint: errcheck
	return int(np)
}

func (v *Buffer) Resume(n int) int {
	if n <= 0 {
		return 0
	}
	np, _ := v.Seek(-int64(n), SeekCurr) //nolint: errcheck
	return int(np)
}

func (v *Buffer) Next(n int) []byte {
	if n <= 0 {
		return nil
	}
	m := v.Len()
	if m == 0 {
		return nil
	}
	if n > m {
		n = m
	}
	b := make([]byte, n)
	v.pos += copy(b[:], v.buf[v.pos:])
	return b
}

func (v *Buffer) ReadByte() (byte, error) {
	m := v.Len()
	if m == 0 {
		return 0, io.EOF
	}
	b := v.buf[v.pos]
	v.pos++
	return b, nil
}

func (v *Buffer) UnreadByte() error {
	if v.pos <= 0 {
		return fmt.Errorf("at beginning")
	}
	v.pos--
	return nil
}

func (v *Buffer) ReadRune() (rune, int, error) {
	m := v.Len()
	if m == 0 {
		return 0, 0, io.EOF
	}

	r, n := utf8.DecodeRune(v.buf[v.pos:])
	v.pos += n
	return r, n, nil
}

func (v *Buffer) UnreadRune() error {
	if v.pos <= 0 {
		return fmt.Errorf("at beginning")
	}

	n := v.pos

	for i := 1; i <= 4 && n-i >= 0; i++ {
		off := n - i
		if utf8.RuneStart(v.buf[off]) {
			n -= i
			break
		}
	}

	v.pos -= n

	return nil
}

func (v *Buffer) ReadBytes(delim byte) ([]byte, error) {
	m := v.Len()
	if m == 0 {
		return nil, io.EOF
	}
	i := bytes.IndexByte(v.buf[v.pos:], delim)
	end := v.pos + i + 1
	if i < 0 {
		end = v.Size()
	}
	b := v.buf[v.pos:end]
	v.pos = end
	return b, nil
}

func (v *Buffer) ReadNextBytes(delim []byte) ([]byte, error) {
	if v.Len() == 0 {
		return nil, io.EOF
	}
	i := bytes.Index(v.buf[v.pos:], delim)
	end := v.pos + i + len(delim)
	if i < 0 {
		end = v.Size()
	}
	b := v.buf[v.pos:end]
	v.pos = end
	return b, nil
}

func (v *Buffer) NextField(sep string, accurate bool) (field []byte, symbol []byte, err error) {
	if v.Len() == 0 {
		return nil, nil, io.EOF
	}

	var i int

	if accurate {
		i = bytes.Index(v.buf[v.pos:], []byte(sep))
	} else {
		i = bytes.IndexAny(v.buf[v.pos:], sep)
	}

	if i < 0 {
		i = v.Size()
	}

	end := v.pos + i
	if end > v.Size() {
		end = v.Size()
	}

	field = v.buf[v.pos:end]
	v.pos = end

	if accurate {
		symbol = []byte(sep)
		v.pos += len(symbol)
	} else {
		rv, rn := utf8.DecodeRune(v.buf[v.pos:])
		symbol = []byte(string(rv))
		v.pos += rn
	}

	if v.pos > v.Size() {
		v.pos = v.Size()
	}

	return
}

func (v *Buffer) ReadString(delim byte) (string, error) {
	b, err := v.ReadBytes(delim)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (v *Buffer) ReadNextString(delim string) (string, error) {
	b, err := v.ReadNextBytes([]byte(delim))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

const (
	SeekStart int = 0
	SeekCurr  int = 1
	SeekEnd   int = 2
)

func (v *Buffer) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case SeekStart:
		v.pos = int(offset)
	case SeekCurr:
		v.pos += int(offset)
	case SeekEnd:
		v.pos = v.Size() + int(offset)
	default:
		return 0, fmt.Errorf("invalid whence")
	}
	if v.pos < 0 {
		v.pos = 0
	} else if v.pos > v.Size() {
		v.pos = v.Size()
	}
	return int64(v.pos), nil
}
