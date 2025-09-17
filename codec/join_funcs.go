/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package codec

import (
	"bytes"
)

func BytesJoin(_ Codec, dst *[]byte, src ...[]byte) error {
	for _, next := range src {
		tmp := bytes.TrimSpace(*dst)
		tmp = append(tmp, '\n', '\n')
		tmp = append(tmp, next...)
		*dst = bytes.TrimSpace(tmp)
	}

	return nil
}

func MapJoin(c Codec, dst *[]byte, src ...[]byte) error {
	out := map[string]interface{}{}

	if len(*dst) > 0 {
		if err := c.Decode(*dst, &out); err != nil {
			return err
		}
	}

	for _, next := range src {
		tmp := map[string]interface{}{}
		if err := c.Decode(next, &tmp); err != nil {
			return err
		}

		mapMerge(out, tmp)
	}

	b, err := c.Encode(out)
	if err != nil {
		return err
	}

	*dst = b
	return nil
}

func mapMerge(dst map[string]interface{}, src ...map[string]interface{}) {
	for _, next := range src {
		for k, v := range next {
			vv, ok := dst[k]
			if !ok {
				dst[k] = v
				continue
			}

			m1, ok1 := vv.(map[string]interface{})
			m2, ok2 := v.(map[string]interface{})
			if ok2 && ok1 {
				mapMerge(m1, m2)
				continue
			}

			dst[k] = v
		}
	}
}
