/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package codec

import (
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

	call, err := _default.Get(v.Ext)
	if err != nil {
		return err
	}
	if v.Blob == nil {
		return nil
	}
	for _, conf := range configs {
		if err = call.Decode(v.Blob, conf); err != nil {
			return err
		}
	}
	return nil
}

func (v *BlobEncoder) Encode(configs ...interface{}) error {
	v.mux.Lock()
	defer v.mux.Unlock()

	codec, err0 := _default.Get(v.Ext)
	if err0 != nil {
		return err0
	}

	out := make(map[string]interface{}, 10)
	for _, conf := range configs {
		bb, err := codec.Encode(conf)
		if err != nil {
			return err
		}
		tmp := make(map[string]interface{}, 10)
		if err = codec.Decode(bb, &tmp); err != nil {
			return err
		}
		if err = codec.Merge(out, tmp); err != nil {
			return err
		}
	}
	v.Blob, err0 = codec.Encode(out)
	return err0
}
