/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package main

import (
	"fmt"

	"go.osspkg.com/ioutils/codec"
)

type Config struct {
	Endpoint string `json:"endpoint"`
	Retries  int    `json:"retries"`
}

func main() {
	blob := &codec.BlobEncoder{Ext: codec.ExtJSON}
	if err := blob.Encode(Config{Endpoint: "https://service.example", Retries: 3}); err != nil {
		panic(err)
	}

	var decoded Config
	if err := blob.Decode(&decoded); err != nil {
		panic(err)
	}
	fmt.Printf("%s (%d retries)\n", decoded.Endpoint, decoded.Retries)
}
