/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package codec

import (
	"bytes"
)

func BytesJoin(_ Codec, src ...[]byte) ([]byte, error) {
	out := make([]byte, 0, 1024)

	for _, next := range src {
		tmp := bytes.TrimSpace(out)
		tmp = append(tmp, '\n', '\n')
		tmp = append(tmp, next...)
		out = bytes.TrimSpace(tmp)
	}

	return out, nil
}

func MapJoin(c Codec, src ...[]byte) ([]byte, error) {
	list := make([]map[string]any, 0, len(src))

	for _, next := range src {
		tmp := map[string]any{}
		if err := c.Decode(next, &tmp); err != nil {
			return nil, err
		}
		if len(tmp) == 0 {
			continue
		}
		list = append(list, tmp)
	}

	out := mapMerge(list...)

	b, err := c.Encode(out)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func mapMerge(src ...map[string]any) map[string]any {
	dst := make(map[string]any, len(src))

	for _, next := range src {
		for k, v := range next {
			vv, ok := dst[k]
			if !ok {
				dst[k] = v
				continue
			}

			m1, ok1 := vv.(map[string]any)
			m2, ok2 := v.(map[string]any)
			if ok2 && ok1 {
				v = mapMerge(m1, m2)
			}

			dst[k] = v
		}
	}

	return dst
}
