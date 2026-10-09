/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package shell_test

import (
	"bytes"
	"context"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"go.osspkg.com/casecheck"

	"go.osspkg.com/ioutils/shell"
)

func TestUnit_ShellCall(t *testing.T) {
	sh := shell.New()
	sh.SetDir(t.TempDir())
	sh.SetEnv("IOUTILS_TEST_VALUE", "shell-test-value")
	out, err := sh.Call(context.Background(), "echo go-ioutils-shell-test")
	casecheck.NoError(t, err)
	if !strings.Contains(string(out), "go-ioutils-shell-test") {
		t.Fatalf("unexpected shell output: %q", out)
	}
}

func TestShell_CallContext(t *testing.T) {
	sh := shell.New()
	sh.SetDir(t.TempDir())
	sh.SetEnv("IOUTILS_TEST_VALUE", "shell-test-value")

	var output bytes.Buffer
	sh.SetOut(&output)
	if runtime.GOOS == "windows" {
		casecheck.NoError(t, sh.CallContext(context.Background(), "echo %IOUTILS_TEST_VALUE%"))
	} else {
		casecheck.NoError(t, sh.CallContext(context.Background(), `echo "$IOUTILS_TEST_VALUE"`))
	}
	if !strings.Contains(output.String(), "shell-test-value") {
		t.Fatalf("CallContext output = %q, want environment value", output.String())
	}
}

func TestShell_CallPackageContext(t *testing.T) {
	sh := shell.New()
	sh.SetDir(t.TempDir())
	var output bytes.Buffer
	sh.SetOut(&output)

	casecheck.NoError(t, sh.CallPackageContext(context.Background(), "echo first", "echo second"))
	for _, want := range []string{"first", "second"} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("CallPackageContext output = %q, missing %q", output.String(), want)
		}
	}
}

func TestShell_CallPackageContextError(t *testing.T) {
	sh := shell.New()
	sh.SetDir(t.TempDir())
	command := "exit 7"
	if runtime.GOOS == "windows" {
		command = "exit /b 7"
	}
	if err := sh.CallPackageContext(context.Background(), "echo before", command, "echo after"); err == nil {
		t.Fatal("CallPackageContext() succeeded after a command failed")
	} else if !strings.Contains(err.Error(), "command #1") {
		t.Fatalf("CallPackageContext() error = %v, want command index", err)
	}
}

func TestShell_UseOSEnvFalse(t *testing.T) {
	sh := shell.New()
	sh.UseOSEnv(false)
	sh.SetEnv("IOUTILS_TEST_VALUE", "shell-test-value")
	var output bytes.Buffer
	sh.SetOut(&output)
	if runtime.GOOS == "windows" {
		casecheck.NoError(t, sh.CallContext(context.Background(), "echo %IOUTILS_TEST_VALUE%"))
	} else {
		casecheck.NoError(t, sh.CallContext(context.Background(), `echo "$IOUTILS_TEST_VALUE"`))
	}
	if !strings.Contains(output.String(), "shell-test-value") {
		t.Fatalf("CallContext output = %q, want explicitly configured environment value", output.String())
	}
}

func TestShell_CallContextMissingExecutable(t *testing.T) {
	sh := shell.New()
	sh.SetShell(filepath.Join(t.TempDir(), "missing-shell"), "-c")
	if err := sh.CallContext(context.Background(), "echo unused"); err == nil {
		t.Fatal("CallContext() succeeded with a missing shell executable")
	}
}
