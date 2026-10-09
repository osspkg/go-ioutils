/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package fs

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"testing"
)

func TestHashCreateAndVerify(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "data.bin")
	if err := os.WriteFile(filename, []byte("payload"), 0o600); err != nil {
		t.Fatal(err)
	}

	h := sha256.New()
	got, err := HashCreate(filename, h)
	if err != nil {
		t.Fatal(err)
	}
	want := "239f59ed55e737c77147cf55ad0c1b030b6d7ee748a7426952f9b852d5a935e5"
	if got != want {
		t.Fatalf("HashCreate() = %q, want %q", got, want)
	}
	if h.Size() != 32 {
		t.Fatalf("hash state size after HashCreate() = %d, want 32", h.Size())
	}

	if err := HashVerify(filename, h, want); err != nil {
		t.Fatalf("HashVerify() returned error for matching hash: %v", err)
	}
	if err := HashVerify(filename, h, "invalid"); err == nil {
		t.Fatal("HashVerify() succeeded for a mismatched hash")
	}
}

func TestHashFilesNotFound(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "missing")
	if _, err := HashCreate(filename, sha256.New()); err == nil {
		t.Fatal("HashCreate() succeeded for a missing file")
	}
	if err := HashVerify(filename, sha256.New(), "hash"); err == nil {
		t.Fatal("HashVerify() succeeded for a missing file")
	}
}

func TestHashVerifyResetsHasherAfterMismatch(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "data")
	if err := os.WriteFile(filename, []byte("payload"), 0o600); err != nil {
		t.Fatal(err)
	}
	h := sha256.New()
	if err := HashVerify(filename, h, "invalid"); err == nil {
		t.Fatal("HashVerify() succeeded for a mismatched hash")
	}
	emptyHash := sha256.Sum256(nil)
	if got := h.Sum(nil); string(got) != string(emptyHash[:]) {
		t.Fatalf("hasher state after mismatch = %x, want empty hash %x", got, emptyHash)
	}
}
