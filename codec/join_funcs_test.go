/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package codec

import (
	"testing"

	"go.osspkg.com/casecheck"
)

func TestUnit_mapMerge(t *testing.T) {
	mapA := map[string]any{
		"qq": "ww",
		"aa": map[string]any{
			"bb": "cc",
		},
		"yy": 123,
		"ww": 123,
	}
	mapB := map[string]any{
		"zz": "xx",
		"aa": map[string]any{
			"ss": "dd",
			"ee": map[string]any{
				"rr": "tt",
			},
		},
		"ww": map[string]any{
			"gg": "hh",
		},
	}

	mapA = mapMerge(mapA, mapB)

	casecheck.Equal(t, map[string]any{
		"yy": 123,
		"zz": "xx",
		"qq": "ww",
		"aa": map[string]any{
			"ss": "dd",
			"bb": "cc",
			"ee": map[string]any{
				"rr": "tt",
			},
		},
		"ww": map[string]any{
			"gg": "hh",
		},
	}, mapA)
}
