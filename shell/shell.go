/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package shell

import (
	"context"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"go.osspkg.com/errors"
)

type (
	object struct {
		env   []string
		dir   string
		shell []string
		osenv bool
		out   io.Writer
		mux   sync.RWMutex
	}

	TShell interface {
		SetEnv(key, value string)
		UseOSEnv(use bool)
		SetDir(dir string)
		SetOut(out io.Writer)
		SetShell(shell string, keys ...string)
		CallPackageContext(ctx context.Context, commands ...string) error
		CallContext(ctx context.Context, command string) error
		Call(ctx context.Context, command string) ([]byte, error)
	}
)

func New() TShell {
	shell, args := defaultShell()
	v := &object{
		osenv: true,
		env:   make([]string, 0, 10),
		dir:   os.TempDir(),
		out:   io.Discard,
		shell: []string{shell, args},
	}
	return v
}

func defaultShell() (string, string) {
	return defaultShellFor(runtime.GOOS, os.Getenv("COMSPEC"))
}

func defaultShellFor(goos, comspec string) (string, string) {
	if goos == "windows" {
		if comspec == "" {
			comspec = "cmd.exe"
		}
		return comspec, "/C"
	}
	return "/bin/sh", "-xec"
}

func (v *object) SetEnv(key, value string) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.env = append(v.env, key+"="+value)
}

func (v *object) UseOSEnv(use bool) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.osenv = use
}

func (v *object) SetDir(dir string) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.dir = dir
}

func (v *object) SetOut(out io.Writer) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.out = out
}

func (v *object) SetShell(shell string, keys ...string) {
	v.mux.Lock()
	defer v.mux.Unlock()

	keysSum := make([]string, 0, len(keys))
	for _, key := range keys {
		if len(key) == 0 {
			continue
		}
		keysSum = append(keysSum, strings.TrimSpace(key))
	}

	v.shell = []string{shell, strings.Join(keysSum, " ")}
}

func (v *object) CallPackageContext(ctx context.Context, commands ...string) error {
	for i, command := range commands {
		if err := v.CallContext(ctx, command); err != nil {
			return errors.Wrapf(err, "call command #%d [%s]", i, command)
		}
	}
	return nil
}

func (v *object) CallContext(ctx context.Context, command string) error {
	v.mux.RLock()
	defer v.mux.RUnlock()

	cmd := exec.CommandContext(ctx, v.shell[0], v.shell[1], command) //nolint:gosec
	cmd.Dir = v.dir
	cmd.Stdout = v.out
	cmd.Stderr = v.out

	if v.osenv {
		cmd.Env = append(os.Environ(), v.env...)
	} else {
		cmd.Env = v.env
	}

	return cmd.Run()
}

func (v *object) Call(ctx context.Context, command string) ([]byte, error) {
	v.mux.RLock()
	defer v.mux.RUnlock()

	cmd := exec.CommandContext(ctx, v.shell[0], v.shell[1], command) //nolint:gosec
	cmd.Dir = v.dir

	if v.osenv {
		cmd.Env = append(os.Environ(), v.env...)
	} else {
		cmd.Env = v.env
	}

	return cmd.CombinedOutput()
}
