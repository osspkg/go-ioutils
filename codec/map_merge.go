/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package codec

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
