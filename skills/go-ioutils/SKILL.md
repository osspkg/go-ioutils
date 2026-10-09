---
name: go-ioutils
description: Use when writing or reviewing Go code that imports go.osspkg.com/ioutils; consult this library's package contracts and verified examples. Not for unrelated Go I/O work.
---

# go-ioutils

Use this skill when a task uses `go.osspkg.com/ioutils` or one of its subpackages. Check the current package API before proposing changes; this library has behavior that differs from similarly named standard-library functions.

Read [the API reference](references/api-reference.md) for the package involved. Use the runnable programs under `examples/` as starting points, and adapt them to the caller's error-handling and lifecycle requirements.

Keep shell commands appropriate to the selected operating system's shell. `shell.New()` chooses a native default, but it does not translate command syntax between shells. Avoid inserting untrusted data into command strings.
