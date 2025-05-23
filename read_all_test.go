/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ioutils

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"go.osspkg.com/errors"
)

type mockReadCloser struct {
	Data     *bytes.Buffer
	ErrRead  error
	ErrClose error
}

func (v *mockReadCloser) Read(p []byte) (int, error) {
	if v.ErrRead != nil {
		return 0, v.ErrRead
	}
	return v.Data.Read(p)
}

func (v *mockReadCloser) Close() error {
	return v.ErrClose
}

func TestUnit_ReadAll(t *testing.T) {
	type args struct {
		r io.ReadCloser
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Case1",
			args: args{
				r: &mockReadCloser{
					Data:     bytes.NewBuffer([]byte(`hello`)),
					ErrRead:  nil,
					ErrClose: nil,
				},
			},
			want:    []byte(`hello`),
			wantErr: false,
		},
		{
			name: "Case2",
			args: args{
				r: &mockReadCloser{
					Data:     nil,
					ErrRead:  errors.New("read error"),
					ErrClose: nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Case3",
			args: args{
				r: &mockReadCloser{
					Data:     nil,
					ErrRead:  errors.New("read error"),
					ErrClose: errors.New("close error"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Case4",
			args: args{
				r: &mockReadCloser{
					Data:     bytes.NewBuffer([]byte(`hello`)),
					ErrRead:  nil,
					ErrClose: errors.New("close error"),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadAll(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}
