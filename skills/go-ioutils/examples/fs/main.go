/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	ioutilsfs "go.osspkg.com/ioutils/fs"
)

func main() {
	dir, err := os.MkdirTemp("", "ioutils-example-")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			panic(err)
		}
	}()

	source := filepath.Join(dir, "source.txt")
	destination := filepath.Join(dir, "copy.txt")
	if err := os.WriteFile(source, []byte("hash this file"), 0o600); err != nil {
		panic(err)
	}
	if err := ioutilsfs.CopyFile(destination, source, 0); err != nil {
		panic(err)
	}
	digest, err := ioutilsfs.HashCreate(destination, sha256.New())
	if err != nil {
		panic(err)
	}
	fmt.Println(digest)
}
