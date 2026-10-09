/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ioutils

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

// --- Mocks ---

type mockReader struct {
	readFunc func(p []byte) (n int, err error)
}

func (m *mockReader) Read(p []byte) (n int, err error) {
	return m.readFunc(p)
}

type mockWriter struct {
	writeFunc func(p []byte) (n int, err error)
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	return m.writeFunc(p)
}

// --- Tests ---

func TestCopy(t *testing.T) {
	src := bytes.Repeat([]byte("a"), packSize+10)
	dst := &bytes.Buffer{}

	n, err := Copy(dst, bytes.NewReader(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(src) {
		t.Fatalf("expected %d bytes, got %d", len(src), n)
	}
	if !bytes.Equal(src, dst.Bytes()) {
		t.Fatal("data mismatch")
	}
}

func TestCopyN(t *testing.T) {
	tests := []struct {
		name       string
		reader     io.Reader
		writer     io.Writer
		size       int
		wantN      int
		wantErr    bool
		errContain string
	}{
		{
			name:   "happy path: exact multiple of size",
			reader: bytes.NewReader([]byte("123456")),
			writer: &bytes.Buffer{},
			size:   3,
			wantN:  6,
		},
		{
			name:   "happy path: less than size",
			reader: bytes.NewReader([]byte("12")),
			writer: &bytes.Buffer{},
			size:   10,
			wantN:  2,
		},
		{
			name:   "happy path: empty reader",
			reader: bytes.NewReader(nil),
			writer: &bytes.Buffer{},
			size:   10,
			wantN:  0,
		},
		{
			name: "reader error (non-EOF)",
			reader: &mockReader{readFunc: func(p []byte) (int, error) {
				return 0, errors.New("read failed")
			}},
			writer:     &bytes.Buffer{},
			size:       10,
			wantErr:    true,
			errContain: "read failed",
		},
		{
			name:   "writer error",
			reader: bytes.NewReader([]byte("data")),
			writer: &mockWriter{writeFunc: func(p []byte) (int, error) {
				return 0, errors.New("write failed")
			}},
			size:       10,
			wantErr:    true,
			errContain: "writer err",
		},
		{
			name: "defensive: negative read bytes",
			reader: &mockReader{readFunc: func(p []byte) (int, error) {
				return -1, nil // Нарушение контракта io.Reader, но код это проверяет
			}},
			writer:     &bytes.Buffer{},
			size:       10,
			wantErr:    true,
			errContain: "negative read bytes",
		},
		{
			name:   "panic: negative size",
			reader: bytes.NewReader([]byte("data")),
			writer: &bytes.Buffer{},
			size:   -1,
			// Ожидаем panic, обработаем отдельно или просто закомментируем,
			// если не хотим использовать recover. В table-driven лучше использовать recover.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Обработка panic для negative size
			if tt.size < 0 {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("got panit: %v", r)
					}
				}()
				CopyN(tt.writer, tt.reader, tt.size)
				return
			}

			n, err := CopyN(tt.writer, tt.reader, tt.size)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errContain)
				}
				if !strings.Contains(err.Error(), tt.errContain) {
					t.Fatalf("expected error containing %q, got %q", tt.errContain, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if n != tt.wantN {
				t.Fatalf("expected %d bytes copied, got %d", tt.wantN, n)
			}
		})
	}
}

func TestCopyN_ZeroSize(t *testing.T) {
	if _, err := CopyN(&bytes.Buffer{}, bytes.NewReader(nil), 0); err == nil {
		t.Fatal("CopyN() with a zero buffer size succeeded")
	}
}

func TestCopyB_ShortReads(t *testing.T) {
	input := []byte("short reads must not truncate the stream")
	output := &bytes.Buffer{}
	n, err := CopyB(output, &shortReader{data: input, chunkSize: 2}, make([]byte, 16))
	if err != nil {
		t.Fatal(err)
	}
	if n != len(input) || !bytes.Equal(output.Bytes(), input) {
		t.Fatalf("CopyB() copied %d bytes %q, want %d bytes %q", n, output.Bytes(), len(input), input)
	}
}

func TestCopyB_PartialWriteErrorCountsWrittenBytes(t *testing.T) {
	wantErr := errors.New("partial write")
	n, err := CopyB(&mockWriter{writeFunc: func([]byte) (int, error) {
		return 2, wantErr
	}}, strings.NewReader("data"), make([]byte, 8))
	if n != 2 {
		t.Fatalf("CopyB() count = %d, want 2", n)
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("CopyB() error = %v, want %v", err, wantErr)
	}
}

type shortReader struct {
	data      []byte
	chunkSize int
}

func (r *shortReader) Read(dst []byte) (int, error) {
	if len(r.data) == 0 {
		return 0, io.EOF
	}
	n := min(len(dst), min(len(r.data), r.chunkSize))
	copy(dst, r.data[:n])
	r.data = r.data[n:]
	return n, nil
}
