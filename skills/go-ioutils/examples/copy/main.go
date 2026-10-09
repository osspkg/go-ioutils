/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package main

import (
	"bytes"
	"fmt"
	"strings"

	"go.osspkg.com/ioutils"
)

func main() {
	var dst bytes.Buffer
	n, err := ioutils.Copy(&dst, strings.NewReader("copy this stream"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("copied %d bytes: %s\n", n, dst.String())
}
