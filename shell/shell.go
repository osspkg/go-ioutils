/*
 *  Copyright (c) 2024-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"go.osspkg.com/errors"
)

type (
	_shell struct {
		env   []string
		dir   string
		shell []string
		osenv bool
		mux   sync.RWMutex
	}

	TShell interface {
		SetEnv(key, value string)
		UseOSEnv(use bool)
		SetDir(dir string)
		SetShell(shell string, keys ...string) error
		CallPackageContext(ctx context.Context, out io.Writer, commands ...string) error
		CallContext(ctx context.Context, out io.Writer, command string) error
		Call(ctx context.Context, command string) ([]byte, error)
	}
)

func New() TShell {
	v := &_shell{
		osenv: true,
		env:   make([]string, 0, 10),
		dir:   os.TempDir(),
		shell: []string{"/bin/sh", "-xec"},
	}
	return v
}

func (v *_shell) SetEnv(key, value string) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.env = append(v.env, key+"="+value)
}

func (v *_shell) UseOSEnv(use bool) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.osenv = use
}

func (v *_shell) SetDir(dir string) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.dir = dir
}

func (v *_shell) SetShell(shell string, keys ...string) error {
	v.mux.Lock()
	defer v.mux.Unlock()

	keysSum := "-"
	for _, key := range keys {
		if len(key) != 1 {
			return fmt.Errorf("invalid key, must have 1 char: %s", key)
		}
		keysSum += key
	}

	v.shell = []string{shell, keysSum}
	return nil
}

func (v *_shell) CallPackageContext(ctx context.Context, out io.Writer, commands ...string) error {
	for i, command := range commands {
		if err := v.CallContext(ctx, out, command); err != nil {
			return errors.Wrapf(err, "call command #%d [%s]", i, command)
		}
	}
	return nil
}

func (v *_shell) CallContext(ctx context.Context, out io.Writer, command string) error {
	v.mux.RLock()
	defer v.mux.RUnlock()

	cmd := exec.CommandContext(ctx, v.shell[0], append(v.shell[1:], command, " <&-")...) //nolint:gosec
	cmd.Dir = v.dir
	cmd.Stdout = out
	cmd.Stderr = out

	if v.osenv {
		cmd.Env = append(os.Environ(), v.env...)
	} else {
		cmd.Env = v.env
	}

	return cmd.Run()
}

func (v *_shell) Call(ctx context.Context, command string) ([]byte, error) {
	v.mux.RLock()
	defer v.mux.RUnlock()

	cmd := exec.CommandContext(ctx, v.shell[0], append(v.shell[1:], command, " <&-")...) //nolint:gosec
	cmd.Dir = v.dir

	if v.osenv {
		cmd.Env = append(os.Environ(), v.env...)
	} else {
		cmd.Env = v.env
	}

	return cmd.CombinedOutput()
}
