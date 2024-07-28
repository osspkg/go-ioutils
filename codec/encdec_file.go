/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package codec

import (
	"os"
	"path/filepath"
)

type FileEncoder string

func (v FileEncoder) Decode(configs ...interface{}) error {
	data, err := os.ReadFile(string(v))
	if err != nil {
		return err
	}
	ext := filepath.Ext(string(v))
	blob := &BlobEncoder{
		Blob: data,
		Ext:  ext,
	}
	return blob.Decode(configs...)
}

func (v FileEncoder) Encode(configs ...interface{}) error {
	ext := filepath.Ext(string(v))
	blob := &BlobEncoder{
		Blob: nil,
		Ext:  ext,
	}
	if err := blob.Encode(configs...); err != nil {
		return err
	}
	return os.WriteFile(string(v), blob.Blob, 0755)
}
