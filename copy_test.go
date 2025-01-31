/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ioutils

import (
	"bytes"
	"io"
	"testing"

	"go.osspkg.com/casecheck"
)

func TestUnit_Copy(t *testing.T) {
	b := make([]byte, 521, 1024)
	in := bytes.NewBuffer(b)
	out := bytes.NewBuffer(nil)
	n, err := Copy(out, in)
	casecheck.NoError(t, err)
	casecheck.Equal(t, 521, n)
	bb, err := io.ReadAll(out)
	casecheck.NoError(t, err)
	casecheck.Equal(t, 521, len(bb))
}

func TestUnit_CopyN(t *testing.T) {
	b := make([]byte, 521, 1024)
	in := bytes.NewBuffer(b)
	out := bytes.NewBuffer(nil)
	n, err := CopyN(out, in, 1)
	casecheck.NoError(t, err)
	casecheck.Equal(t, 521, n)
	bb, err := io.ReadAll(out)
	casecheck.NoError(t, err)
	casecheck.Equal(t, 521, len(bb))
}
