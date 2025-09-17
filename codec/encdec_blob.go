/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package codec

import (
	"fmt"
	"sync"
)

type BlobEncoder struct {
	Blob []byte
	Ext  string
	mux  sync.RWMutex
}

func (v *BlobEncoder) Decode(configs ...interface{}) error {
	v.mux.RLock()
	defer v.mux.RUnlock()

	if len(v.Blob) == 0 {
		return nil
	}

	c, err := _default.Get(v.Ext)
	if err != nil {
		return fmt.Errorf("get codec: %w", err)
	}

	for _, conf := range configs {
		if err = c.Decode(v.Blob, conf); err != nil {
			return fmt.Errorf("decode bytes: %w", err)
		}
	}

	return nil
}

func (v *BlobEncoder) Encode(configs ...interface{}) error {
	v.mux.Lock()
	defer v.mux.Unlock()

	c, err := _default.Get(v.Ext)
	if err != nil {
		return fmt.Errorf("get codec: %w", err)
	}

	out := make([]byte, 0, 1024)
	for _, conf := range configs {
		bb, err0 := c.Encode(conf)
		if err0 != nil {
			return fmt.Errorf("encode bytes: %w", err0)
		}

		if err0 = c.Join(c, &out, bb); err0 != nil {
			return fmt.Errorf("join bytes: %w", err0)
		}
	}

	v.Blob = out

	return nil
}
