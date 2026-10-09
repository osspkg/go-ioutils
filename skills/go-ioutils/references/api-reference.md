# go-ioutils API reference

Module path: `go.osspkg.com/ioutils`. The root package is imported as `ioutils`; the subpackages are `cache`, `codec`, `fs`, `pool`, and `shell`.

Use only the sections relevant to the task. Confirm details against the implementation when changing package behavior.

## Root package

- `Copy(w, r)` copies a stream using a 512-byte buffer.
- `CopyN(w, r, size)` copies until EOF using a buffer of `size` bytes. Despite its name, `size` is a buffer size, not a byte limit.
- `CopyB(w, r, buffer)` copies until EOF using the supplied buffer.
- `Pipe(w, r, size)` copies until EOF using a buffer of `size` bytes.
- The copy helpers return the number of bytes accepted by the writer. A short read does not mean EOF. Repeated `(0, nil)` reads return `io.ErrNoProgress` after 100 consecutive attempts.
- `ReadAll(r)` reads all input and closes `r`, returning read or close errors.
- `IsAsciiEOF(b)` is true only when `b` equals the package's five-byte marker exactly.

## `fs`

- `CurrentDir()` returns the current directory, or `"."` if the OS call fails.
- `DirName(path)` returns the last component of the path's parent directory.
- `FileExist(path)` returns false for any `os.Stat` error, including permission errors.
- `SearchFiles` matches exact file names recursively; `SearchFilesByExt` matches the exact extension, including its leading dot.
- `ListFiles` invokes its handler for files only and returns an error for a nil handler.
- `RewriteFile` reads a file, passes its bytes to a callback, and writes the callback result only if the callback succeeds. A nil callback is invalid.
- `CopyFile(dst, src, mode)` copies file contents. It returns an error if source and destination refer to the same file. A zero mode uses the source mode when creating the destination.
- `HashCreate(path, h)` returns the lowercase hexadecimal digest; `HashVerify(path, h, expected)` compares against it. Both reset `h` after successful reading and digest calculation.

Use `filepath.Join` to build paths. File permission modes have platform-specific effects, especially on Windows.

## `codec`

`BlobEncoder` encodes or decodes in-memory data. Set `Ext` to one of the exported extension constants: `.yaml`, `.yml`, `.json`, `.toml`, `.xml`, `.unic`, or `.conf`. Extensions are matched exactly. `FileEncoder(path)` selects the codec from `filepath.Ext(path)`.

`AddCodec(ext, codec)` registers a process-wide codec. Registration is synchronized. Encode and decode callbacks must both be valid for operations that use them. JSON encoding of multiple input values merges their object fields recursively; later scalar values replace earlier values.

## `cache`

`cache.New[K, V]()` returns a concurrency-safe map-backed cache. `Get` and `Extract` use the usual value/boolean presence result; `Extract` also deletes the entry. `Yield(limit)` snapshots keys when the sequence is created, then looks up current values as it iterates. Entries removed before iteration are skipped.

`Replace(map)` makes a shallow copy of the supplied map. Values are still the caller's `V` values, so pointers or other reference values remain shared. `One()` chooses randomly from at most 30 keys sampled from map iteration; do not depend on uniform selection across the whole cache.

Cleanup options run in background goroutines and stop when their context is canceled. Pass a context whose lifetime matches the cache, and cancel it when cleanup is no longer needed. Expiration timestamps are Unix seconds; values expire when `Timestamp() < now.Unix()`.

## `pool`

`pool.New[T]` requires pooled values to implement `Reset()`. `Put` calls `Reset` before returning the value to `sync.Pool`; `Get` may create a new value because `sync.Pool` can discard stored values at any time. Do not rely on a pooled object being returned to a later `Get`.

`NewSlicePool[T](length, capacity)` creates a pool of `*pool.Slice[T]`. `Reset` clears the current slice contents and sets its length to zero while retaining capacity.

## `shell`

`shell.New()` uses `/bin/sh -xec` on Linux and macOS. On Windows it uses `%COMSPEC%` (or `cmd.exe`) with `/C`. The command string must use that shell's syntax; this API is not a cross-platform command-language translator.

- `Call(ctx, command)` returns combined stdout and stderr.
- `SetOut(w)` configures the writer used by `CallContext` and `CallPackageContext`.
- `SetEnv` adds or overrides child environment values; newline characters are allowed in values, but NUL bytes are not.
- `UseOSEnv(false)` starts the child with only the values explicitly added with `SetEnv`.
- `SetShell(executable, flags...)` joins the supplied flag strings into one argument. Keep flags in a form accepted as one argument by the selected shell.
- Context cancellation terminates the child process through `exec.CommandContext`.

The package intentionally executes shell command strings. Keep command strings trusted or use another API that passes executable arguments directly when shell interpretation is unnecessary.

## Examples

- [Copy a stream](../examples/copy/main.go)
- [Use the cache](../examples/cache/main.go)
- [Encode and decode JSON](../examples/codec/main.go)
- [Run a shell command](../examples/shell/main.go)
- [Copy and hash a file](../examples/fs/main.go)
