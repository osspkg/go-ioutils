/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package main

import (
	"context"
	"fmt"

	"go.osspkg.com/ioutils/shell"
)

func main() {
	sh := shell.New()
	out, err := sh.Call(context.Background(), "echo hello")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(out))
}
