/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package codec

import (
	"fmt"
)

func SimpleCodec(
	enc func(in any) (out []byte, err error),
	dec func(in []byte, out any) (err error),
	join func(c Codec, src ...[]byte) ([]byte, error),
) Codec {
	c := Codec{}
	c.Encode = func(args ...any) ([]byte, error) {
		list := make([][]byte, 0, len(args))
		for _, arg := range args {
			b, err := enc(arg)
			if err != nil {
				return nil, fmt.Errorf("encode bytes: %w", err)
			}
			list = append(list, b)
		}

		switch len(list) {
		case 0:
			return nil, nil
		case 1:
			return list[0], nil
		default:
			b, err := join(c, list...)
			if err != nil {
				return nil, fmt.Errorf("join bytes: %w", err)
			}
			return b, nil
		}
	}
	c.Decode = func(b []byte, args ...any) error {
		for _, arg := range args {
			if err := dec(b, arg); err != nil {
				return fmt.Errorf("decode bytes: %w", err)
			}
		}
		return nil
	}
	return c
}
