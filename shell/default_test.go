/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package shell

import "testing"

func TestDefaultShellFor(t *testing.T) {
	tests := []struct {
		name      string
		goos      string
		comspec   string
		wantShell string
		wantArgs  string
	}{
		{
			name:      "windows uses COMSPEC",
			goos:      "windows",
			comspec:   `C:\Windows\System32\cmd.exe`,
			wantShell: `C:\Windows\System32\cmd.exe`,
			wantArgs:  "/C",
		},
		{
			name:      "windows falls back to cmd.exe",
			goos:      "windows",
			wantShell: "cmd.exe",
			wantArgs:  "/C",
		},
		{
			name:      "unix uses sh",
			goos:      "linux",
			wantShell: "/bin/sh",
			wantArgs:  "-xec",
		},
		{
			name:      "macos uses sh",
			goos:      "darwin",
			wantShell: "/bin/sh",
			wantArgs:  "-xec",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotShell, gotArgs := defaultShellFor(tt.goos, tt.comspec)
			if gotShell != tt.wantShell || gotArgs != tt.wantArgs {
				t.Fatalf("defaultShellFor() = (%q, %q), want (%q, %q)", gotShell, gotArgs, tt.wantShell, tt.wantArgs)
			}
		})
	}
}
