/*
 *  Copyright (c) 2022-2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package shell_test

import (
	"context"
	"testing"

	"go.osspkg.com/casecheck"
	"go.osspkg.com/ioutils/shell"
)

func TestUnit_ShellCall(t *testing.T) {
	sh := shell.New()
	sh.SetDir("/tmp")
	sh.SetEnv("LANG", "en_US.UTF-8")
	casecheck.NoError(t, sh.SetShell("/bin/bash", "x", "c"))
	out, err := sh.Call(context.TODO(), "ls -la /tmp")
	casecheck.NoError(t, err)
	t.Log(string(out))
}
