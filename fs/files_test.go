/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package fs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestCurrentDir(t *testing.T) {
	if CurrentDir() == "" {
		t.Fatal("CurrentDir returned an empty path")
	}
}

func TestDirName(t *testing.T) {
	path := filepath.Join("parent", "child", "file.txt")
	if got, want := DirName(path), "child"; got != want {
		t.Fatalf("DirName(%q) = %q, want %q", path, got, want)
	}
}

func TestFileExist(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "present.txt")
	if err := os.WriteFile(filename, []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	if !FileExist(filename) {
		t.Fatalf("FileExist(%q) = false, want true", filename)
	}
	if FileExist(filepath.Join(t.TempDir(), "missing.txt")) {
		t.Fatal("FileExist returned true for a missing file")
	}
}

func TestSearchFiles(t *testing.T) {
	root := t.TempDir()
	want := filepath.Join(root, "nested", "match.txt")
	writeTestFile(t, want, "match")
	writeTestFile(t, filepath.Join(root, "other.txt"), "other")
	if err := os.Mkdir(filepath.Join(root, "match.txt"), 0o700); err != nil {
		t.Fatal(err)
	}

	got, err := SearchFiles(root, "match.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != want {
		t.Fatalf("SearchFiles() = %v, want [%q]", got, want)
	}

	if _, err := SearchFiles(filepath.Join(root, "missing"), "match.txt"); err == nil {
		t.Fatal("SearchFiles() with a missing directory succeeded")
	}
}

func TestSearchFilesByExt(t *testing.T) {
	root := t.TempDir()
	want := filepath.Join(root, "nested", "match.yaml")
	writeTestFile(t, want, "match")
	writeTestFile(t, filepath.Join(root, "other.yml"), "other")

	got, err := SearchFilesByExt(root, ".yaml")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != want {
		t.Fatalf("SearchFilesByExt() = %v, want [%q]", got, want)
	}

	if _, err := SearchFilesByExt(filepath.Join(root, "missing"), ".yaml"); err == nil {
		t.Fatal("SearchFilesByExt() with a missing directory succeeded")
	}
}

func TestListFiles(t *testing.T) {
	root := t.TempDir()
	first := filepath.Join(root, "one.txt")
	second := filepath.Join(root, "nested", "two.txt")
	writeTestFile(t, first, "1")
	writeTestFile(t, second, "22")

	if err := ListFiles(root, nil); err == nil {
		t.Fatal("ListFiles() with a nil handler succeeded")
	}

	got := make(map[string]fs.FileInfo)
	err := ListFiles(root, func(path string, info fs.FileInfo) {
		got[path] = info
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[first] == nil || got[second] == nil {
		t.Fatalf("ListFiles() returned paths %v, want %q and %q", got, first, second)
	}
	if err := ListFiles(filepath.Join(root, "missing"), func(string, fs.FileInfo) {}); err == nil {
		t.Fatal("ListFiles() with a missing directory succeeded")
	}
}

func TestRewriteFile(t *testing.T) {
	t.Run("create and rewrite", func(t *testing.T) {
		filename := filepath.Join(t.TempDir(), "config.txt")
		err := RewriteFile(filename, func(src []byte) ([]byte, error) {
			if len(src) != 0 {
				t.Fatalf("callback input = %q, want empty", src)
			}
			return []byte("written"), nil
		})
		if err != nil {
			t.Fatal(err)
		}
		assertFileContent(t, filename, "written")
	})

	t.Run("callback error leaves existing file unchanged", func(t *testing.T) {
		filename := filepath.Join(t.TempDir(), "config.txt")
		writeTestFile(t, filename, "original")
		wantErr := errors.New("callback failed")
		err := RewriteFile(filename, func([]byte) ([]byte, error) {
			return nil, wantErr
		})
		if !errors.Is(err, wantErr) {
			t.Fatalf("RewriteFile() error = %v, want %v", err, wantErr)
		}
		assertFileContent(t, filename, "original")
	})
}

func TestRewriteFile_InvalidPath(t *testing.T) {
	err := RewriteFile(filepath.Join(t.TempDir(), "missing", "file.txt"), func(src []byte) ([]byte, error) {
		return src, nil
	})
	if err == nil {
		t.Fatal("RewriteFile() with an invalid path succeeded")
	}
}

func TestCopyFile(t *testing.T) {
	t.Run("copy with source mode", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "source.txt")
		dst := filepath.Join(dir, "destination.txt")
		writeTestFile(t, src, "copy me")
		if err := CopyFile(dst, src, 0); err != nil {
			t.Fatal(err)
		}
		assertFileContent(t, dst, "copy me")
	})

	t.Run("copy with explicit mode", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "source.txt")
		dst := filepath.Join(dir, "destination.txt")
		writeTestFile(t, src, "copy me")
		if err := CopyFile(dst, src, 0o600); err != nil {
			t.Fatal(err)
		}
		assertFileContent(t, dst, "copy me")
	})

	t.Run("source does not exist", func(t *testing.T) {
		err := CopyFile(filepath.Join(t.TempDir(), "destination.txt"), filepath.Join(t.TempDir(), "missing.txt"), 0)
		if err == nil {
			t.Fatal("CopyFile() with a missing source succeeded")
		}
	})

	t.Run("source and destination are the same file", func(t *testing.T) {
		filename := filepath.Join(t.TempDir(), "same.txt")
		writeTestFile(t, filename, "preserve this")
		if err := CopyFile(filename, filename, 0); err == nil {
			t.Fatal("CopyFile() with the same source and destination succeeded")
		}
		assertFileContent(t, filename, "preserve this")
	})

	t.Run("destination directory does not exist", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "source.txt")
		writeTestFile(t, src, "copy me")
		err := CopyFile(filepath.Join(dir, "missing", "destination.txt"), src, 0)
		if err == nil {
			t.Fatal("CopyFile() with an invalid destination succeeded")
		}
	})
}

func writeTestFile(t *testing.T, filename, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filename), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filename, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func assertFileContent(t *testing.T, filename, want string) {
	t.Helper()
	got, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("file content = %q, want %q", got, want)
	}
}
