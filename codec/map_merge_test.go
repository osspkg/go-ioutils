/*
 *  Copyright (c) 2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package codec

import (
	"testing"

	"go.osspkg.com/casecheck"
)

func TestUnit_mapMerge(t *testing.T) {
	mapA := map[string]interface{}{
		"qq": "ww",
		"aa": map[string]interface{}{
			"bb": "cc",
		},
	}
	mapB := map[string]interface{}{
		"zz": "xx",
		"aa": map[string]interface{}{
			"ss": "dd",
			"ee": map[string]interface{}{
				"rr": "tt",
			},
		},
	}

	casecheck.NoError(t, mapMerge(mapA, mapB))
	casecheck.Equal(t, map[string]interface{}{
		"zz": "xx",
		"qq": "ww",
		"aa": map[string]interface{}{
			"ss": "dd",
			"bb": "cc",
			"ee": map[string]interface{}{
				"rr": "tt",
			},
		},
	}, mapA)
}