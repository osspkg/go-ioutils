/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ioutils

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"go.osspkg.com/errors"
)

// --- Tests ---

func TestPipe(t *testing.T) {
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
			name:   "happy path: empty reader (immediate EOF)",
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
			name:   "defensive: short write",
			reader: bytes.NewReader([]byte("12345")),
			writer: &mockWriter{writeFunc: func(p []byte) (int, error) {
				return 2, nil // Writer вернул меньше байт, но без ошибки
			}},
			size:       10,
			wantN:      2,
			wantErr:    true,
			errContain: "short write", // io.ErrShortWrite.Error()
		},
		{
			name:       "validation: size <= 0",
			reader:     bytes.NewReader([]byte("data")),
			writer:     &bytes.Buffer{},
			size:       0,
			wantErr:    true,
			errContain: "size must be greater than zero",
		},
		{
			name:       "validation: nil writer",
			reader:     bytes.NewReader([]byte("data")),
			writer:     nil,
			size:       10,
			wantErr:    true,
			errContain: "writer must not be nil",
		},
		{
			name:       "validation: nil reader",
			reader:     nil,
			writer:     &bytes.Buffer{},
			size:       10,
			wantErr:    true,
			errContain: "reader must not be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := Pipe(tt.writer, tt.reader, tt.size)

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

// Тест на защиту от зависания при некорректном поведении io.Reader
func TestPipe_InfiniteLoopOnZeroRead(t *testing.T) {
	reader := &mockReader{readFunc: func(p []byte) (int, error) {
		return 0, nil
	}}

	if _, err := Pipe(&bytes.Buffer{}, reader, 10); !errors.Is(err, io.ErrNoProgress) {
		t.Fatalf("Pipe() error = %v, want %v", err, io.ErrNoProgress)
	}
}

func TestPipe_ShortReads(t *testing.T) {
	input := []byte("short reads must not truncate the stream")
	output := &bytes.Buffer{}
	n, err := Pipe(output, &shortReader{data: input, chunkSize: 2}, 16)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(input) || !bytes.Equal(output.Bytes(), input) {
		t.Fatalf("Pipe() copied %d bytes %q, want %d bytes %q", n, output.Bytes(), len(input), input)
	}
}

func TestPipe_PartialWriteErrorCountsWrittenBytes(t *testing.T) {
	wantErr := errors.New("partial write")
	n, err := Pipe(&mockWriter{writeFunc: func([]byte) (int, error) {
		return 2, wantErr
	}}, strings.NewReader("data"), 8)
	if n != 2 {
		t.Fatalf("Pipe() count = %d, want 2", n)
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("Pipe() error = %v, want %v", err, wantErr)
	}
}
