/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package codec

func mapMerge(dst map[string]interface{}, src map[string]interface{}) error {
	for k, v := range src {
		vv, ok := dst[k]
		if !ok {
			dst[k] = v
			continue
		}
		vMap, vOk := v.(map[string]interface{})
		vvMap, vvOk := vv.(map[string]interface{})
		if vOk && vvOk {
			return mapMerge(vvMap, vMap)
		}
		dst[k] = v
	}

	return nil
}
