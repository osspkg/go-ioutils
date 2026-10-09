# Instructions for agents

## Project

This repository is the Go module `go.osspkg.com/ioutils` and requires Go 1.26 or newer. It provides stream, filesystem, codec, cache, pool, and shell utilities.

## Repository map

- Root Go files contain stream helpers: `Copy`, `CopyN`, `CopyB`, `Pipe`, `ReadAll`, and `IsAsciiEOF`.
- `cache/`, `codec/`, `fs/`, `pool/`, and `shell/` contain the corresponding subpackages.
- `skills/go-ioutils/` contains library-specific guidance and runnable examples; these example packages are included by `go test ./...`.
- Tests live beside their packages in `*_test.go` files.
- `.github/workflows/ci.yml` runs tests and `go vet` on Linux, macOS, and Windows. The Ubuntu job also runs `make ci`.

## Commands

Run from the repository root:

```sh
go test ./...
go test -race ./...
go vet ./...
```

`make ci` runs `goppy setup-lib`, license checks, lint, tests, and build. The Makefile installs the latest `goppy` before running these tasks, so this command needs network access and is primarily the Ubuntu CI entry point.

## Package behavior to preserve

- `CopyN` uses `size` as its buffer size and copies until EOF; it does not limit the total bytes copied.
- `ReadAll` closes the supplied reader and returns read or close errors.
- `cache.Replace` copies the supplied map shallowly; values that contain references remain shared.
- `shell.New` selects a native shell, but command syntax is still shell-specific. Avoid hard-coded `/bin/sh`, `/tmp`, and Unix-only commands in cross-platform tests; use `t.TempDir()` and `filepath` for paths.
- `pool.Pool.Put` calls `Reset` before returning an item to `sync.Pool`; callers must not keep using an item after putting it back.

Keep changes scoped to the affected package, update tests for changed behavior, and avoid changing unrelated working-tree edits.
