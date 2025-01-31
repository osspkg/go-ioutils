/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package hash

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"

	"go.osspkg.com/errors"
)

func Verify(filename string, h hash.Hash, valid string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close() // nolint: errcheck
	if _, err = io.Copy(h, r); err != nil {
		return errors.Wrapf(err, "calculate file hash")
	}
	result := hex.EncodeToString(h.Sum(nil))
	h.Reset()
	if result != valid {
		return fmt.Errorf("invalid hash: expected[%s] actual[%s]", valid, result)
	}
	return nil
}

func Create(filename string, h hash.Hash) (string, error) {
	r, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer r.Close() // nolint: errcheck
	if _, err = io.Copy(h, r); err != nil {
		return "", errors.Wrapf(err, "calculate file hash")
	}
	result := hex.EncodeToString(h.Sum(nil))
	h.Reset()
	return result, nil
}
