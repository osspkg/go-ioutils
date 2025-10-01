/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package data

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"go.osspkg.com/casecheck"
)

type mockRW struct {
	Err error
	N   int
}

func (v *mockRW) Read(p []byte) (int, error) {
	if v.Err != nil {
		return 0, v.Err
	}
	return v.N, nil
}
func (v *mockRW) Write(p []byte) (int, error) {
	if v.Err != nil {
		return 0, v.Err
	}
	return v.N, nil
}

func TestUnit_Data1(t *testing.T) {
	d := NewBuffer(10)

	n, err := d.WriteString("aФ卉**")
	casecheck.NoError(t, err)
	casecheck.Equal(t, 9, n)

	casecheck.Error(t, d.UnreadRune())
	casecheck.Error(t, d.UnreadByte())

	r, s, err := d.ReadRune()
	casecheck.NoError(t, err)
	casecheck.Equal(t, rune('a'), r)
	casecheck.Equal(t, 1, s)

	r, s, err = d.ReadRune()
	casecheck.NoError(t, err)
	casecheck.Equal(t, rune('Ф'), r)
	casecheck.Equal(t, 2, s)

	b, err := d.ReadByte()
	casecheck.NoError(t, err)
	casecheck.Equal(t, byte(240), b)

	casecheck.NoError(t, d.UnreadRune())

	r, s, err = d.ReadRune()
	casecheck.NoError(t, err)
	casecheck.Equal(t, rune('Ф'), r)
	casecheck.Equal(t, 2, s)

	r, s, err = d.ReadRune()
	casecheck.NoError(t, err)
	casecheck.Equal(t, rune('卉'), r)
	casecheck.Equal(t, 4, s)

	bs, err := io.ReadAll(d)
	casecheck.NoError(t, err)
	casecheck.Equal(t, []byte("**"), bs)

	r, s, err = d.ReadRune()
	casecheck.Error(t, err)
	casecheck.Equal(t, rune(0), r)
	casecheck.Equal(t, 0, s)

	b, err = d.ReadByte()
	casecheck.Error(t, err)
	casecheck.Equal(t, byte(0), b)

	casecheck.NoError(t, d.UnreadByte())

	b, err = d.ReadByte()
	casecheck.NoError(t, err)
	casecheck.Equal(t, byte(42), b)

	i, err := d.Seek(0, 0)
	casecheck.NoError(t, err)
	casecheck.Equal(t, int64(0), i)

	i, err = d.Seek(-5, 0)
	casecheck.NoError(t, err)
	casecheck.Equal(t, int64(0), i)

	i, err = d.Seek(5, 1)
	casecheck.NoError(t, err)
	casecheck.Equal(t, int64(5), i)

	i, err = d.Seek(5, 2)
	casecheck.NoError(t, err)
	casecheck.Equal(t, int64(9), i)

	i, err = d.Seek(5, 3)
	casecheck.Error(t, err)
	casecheck.Equal(t, int64(0), i)

	d.Seek(0, 0)

	str, err := d.ReadString('*')
	casecheck.NoError(t, err)
	casecheck.Equal(t, "aФ卉*", str)

	str, err = d.ReadString('*')
	casecheck.NoError(t, err)
	casecheck.Equal(t, "*", str)

	d.Seek(0, 0)

	str, err = d.ReadString('+')
	casecheck.NoError(t, err)
	casecheck.Equal(t, "aФ卉**", str)

	str, err = d.ReadString('+')
	casecheck.Error(t, err)
	casecheck.Equal(t, "", str)

	d.Seek(0, 0)

	str, err = d.ReadNextString("卉*")
	casecheck.NoError(t, err)
	casecheck.Equal(t, "aФ卉*", str)

	d.Seek(0, 0)

	str, err = d.ReadNextString("*卉")
	casecheck.NoError(t, err)
	casecheck.Equal(t, "aФ卉**", str)

	str, err = d.ReadNextString("*卉")
	casecheck.Error(t, err)
	casecheck.Equal(t, "", str)

	d.Seek(0, 0)

	bs = d.Next(3)
	casecheck.NotNil(t, bs)
	casecheck.Equal(t, []byte("aФ"), bs)

	bs = d.Next(3)
	casecheck.NotNil(t, bs)
	casecheck.Equal(t, []byte{240, 175, 160}, bs)

	bs = d.Next(0)
	casecheck.Nil(t, bs)

	bs = d.Next(10)
	casecheck.NotNil(t, bs)
	casecheck.Equal(t, []byte{172, 42, 42}, bs)

	bs = d.Next(10)
	casecheck.Nil(t, bs)

	d.Seek(0, 0)

	bs = make([]byte, 2)

	n, err = d.ReadAt(bs, 2)
	casecheck.NoError(t, err)
	casecheck.Equal(t, 2, n)
	casecheck.Equal(t, []byte{164, 240}, bs)

	n, err = d.ReadAt(nil, 2)
	casecheck.Error(t, err)
	casecheck.Equal(t, 0, n)

	n, err = d.ReadAt(bs, 10)
	casecheck.Error(t, err)
	casecheck.Equal(t, 0, n)

	bb := bytes.NewBufferString("123")

	i, err = d.ReadFrom(bb)
	casecheck.NoError(t, err)
	casecheck.Equal(t, int64(3), i)
	casecheck.Equal(t, "aФ卉**123", d.String())

	i, err = d.ReadFrom(&mockRW{Err: fmt.Errorf("1")})
	casecheck.Error(t, err)
	casecheck.Equal(t, int64(0), i)

	i, err = d.ReadFrom(&mockRW{N: -1})
	casecheck.Error(t, err)
	casecheck.Equal(t, int64(0), i)

	d.Seek(0, 0)

	bb = bytes.NewBufferString("123")
	i, err = d.WriteTo(bb)
	casecheck.NoError(t, err)
	casecheck.Equal(t, int64(12), i)
	casecheck.Equal(t, "123aФ卉**123", bb.String())

	i, err = d.WriteTo(bb)
	casecheck.NoError(t, err)
	casecheck.Equal(t, int64(0), i)

	d.Seek(0, 0)

	i, err = d.WriteTo(&mockRW{Err: fmt.Errorf("1")})
	casecheck.Error(t, err)
	casecheck.Equal(t, int64(0), i)

	casecheck.NoError(t, d.WriteByte('1'))

	n, err = d.WriteRune('Ф')
	casecheck.NoError(t, err)
	casecheck.Equal(t, 2, n)

	casecheck.Equal(t, "aФ卉**1231Ф", d.String())

	d.Reset()
	d.WriteString("00000")
	n, err = d.WriteAt([]byte("11"), -1)
	casecheck.NoError(t, err)
	casecheck.Equal(t, 2, n)
	casecheck.Equal(t, "11000", d.String())

	d.Reset()
	d.WriteString("00000")
	n, err = d.WriteAt([]byte("11"), 2)
	casecheck.NoError(t, err)
	casecheck.Equal(t, 2, n)
	casecheck.Equal(t, "00110", d.String())

	d.Reset()
	d.WriteString("00000")
	n, err = d.WriteAt([]byte("11"), 4)
	casecheck.NoError(t, err)
	casecheck.Equal(t, 2, n)
	casecheck.Equal(t, "000011", d.String())

	d.Reset()
	d.WriteString("00000")
	n, err = d.WriteAt([]byte("11"), 10)
	casecheck.NoError(t, err)
	casecheck.Equal(t, 2, n)
	casecheck.Equal(t, []byte{48, 48, 48, 48, 48, 0, 0, 0, 0, 0, 49, 49}, d.Bytes())
}

func TestUnit_Data_Truncate(t *testing.T) {
	tests := []struct {
		in, out string
		trc     int
		ra      bool
		s, l    int
	}{
		{
			in:  "aФ卉**",
			out: "aФ卉**",
			trc: 0,
			ra:  false,
			s:   9,
			l:   9,
		},
		{
			in:  "aФ卉**",
			out: "aФ卉*",
			trc: 1,
			ra:  false,
			s:   8,
			l:   8,
		},
		{
			in:  "aФ卉**",
			out: "aФ卉",
			trc: 2,
			ra:  false,
			s:   7,
			l:   7,
		},
		{
			in:  "aФ卉**",
			out: "aФ",
			trc: 3,
			ra:  false,
			s:   3,
			l:   3,
		},
		{
			in:  "aФ卉**",
			out: "aФ",
			trc: 5,
			ra:  false,
			s:   3,
			l:   3,
		},
		{
			in:  "aФ卉**",
			out: "a",
			trc: 8,
			ra:  false,
			s:   1,
			l:   1,
		},
		{
			in:  "aФ卉**",
			out: "",
			trc: 9,
			ra:  false,
			s:   0,
			l:   0,
		},
		{
			in:  "aФ卉**",
			out: "",
			trc: 10,
			ra:  false,
			s:   0,
			l:   0,
		},
		{
			in:  "aФ卉**",
			out: "aФ",
			trc: 3,
			ra:  true,
			s:   3,
			l:   0,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			d := NewBuffer(10)
			_, err := d.WriteString(tt.in)
			casecheck.NoError(t, err)
			if tt.ra {
				n, err := d.Read(nil)
				casecheck.Error(t, err)
				casecheck.Equal(t, 0, n)

				b, err := io.ReadAll(d)
				casecheck.NoError(t, err)
				casecheck.Equal(t, len(tt.in), len(b))

				n, err = d.Read(b)
				casecheck.Error(t, err)
				casecheck.Equal(t, 0, n)
			}
			d.Truncate(tt.trc)
			casecheck.Equal(t, tt.out, d.String())
			casecheck.Equal(t, tt.s, d.Size())
			casecheck.Equal(t, tt.l, d.Len())
		})
	}
}

type mockLimitWriter struct {
	B   []byte
	Lim int
}

func (v *mockLimitWriter) Write(p []byte) (int, error) {
	n := min(len(p), v.Lim)
	v.B = append(v.B, p[:n]...)
	return n, nil
}

func TestUnit_WriteTo(t *testing.T) {
	d := NewBuffer(10)
	d.WriteString("1234567890")

	lw := &mockLimitWriter{Lim: 3}
	n, err := d.WriteTo(lw)
	casecheck.NoError(t, err)
	casecheck.Equal(t, int64(10), n)
	casecheck.Equal(t, "1234567890", string(lw.B))
}

func TestUnit_NextField_Accurate_False(t *testing.T) {
	d := NewBuffer(1)
	d.WriteString("123  456 `789`")

	b, s, err := d.NextField(" `", false)
	casecheck.NoError(t, err)
	casecheck.Equal(t, "123", string(b))
	casecheck.Equal(t, " ", string(s))

	b, s, err = d.NextField(" `", false)
	casecheck.NoError(t, err)
	casecheck.Equal(t, "", string(b))
	casecheck.Equal(t, " ", string(s))

	b, s, err = d.NextField(" `", false)
	casecheck.NoError(t, err)
	casecheck.Equal(t, "456", string(b))
	casecheck.Equal(t, " ", string(s))

	b, s, err = d.NextField(" `", false)
	casecheck.NoError(t, err)
	casecheck.Equal(t, "", string(b))
	casecheck.Equal(t, "`", string(s))

	b, s, err = d.NextField(" `", false)
	casecheck.NoError(t, err)
	casecheck.Equal(t, "789", string(b))
	casecheck.Equal(t, "`", string(s))

	b, s, err = d.NextField(" `", false)
	casecheck.Error(t, err)
	casecheck.Equal(t, "", string(b))
	casecheck.Equal(t, []byte{}, s)
}

func TestUnit_NextField_Accurate_True(t *testing.T) {
	d := NewBuffer(1)
	d.WriteString("123  456 `789`")

	b, s, err := d.NextField(" `", true)
	casecheck.NoError(t, err)
	casecheck.Equal(t, "123  456", string(b))
	casecheck.Equal(t, " `", string(s))

	b, s, err = d.NextField(" `", true)
	casecheck.NoError(t, err)
	casecheck.Equal(t, "789`", string(b))
	casecheck.Equal(t, " `", string(s))

	b, s, err = d.NextField(" `", false)
	casecheck.Error(t, err)
	casecheck.Equal(t, "", string(b))
	casecheck.Equal(t, []byte{}, s)
}
