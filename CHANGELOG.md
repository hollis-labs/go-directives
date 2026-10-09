# Changelog

## Repository retirement — 2026-10-09

### Changed

- Maintained development moved to [github.com/hollis-labs/libs/ui-go/directives](https://github.com/hollis-labs/libs/tree/ui-go%2Fv0.1.0/ui-go/directives) in
  `github.com/hollis-labs/libs/ui-go@v0.1.0` (`ui-go/v0.1.0`).
- This standalone repository is retired after the replacement release was
  verified fetchable with successful module CI. README migration instructions
  identify the new import prefix; existing standalone tags and history are preserved.

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## v0.1.0 — 2026-05-10

Initial public release.

### Added

- `Parse(input, ParserConfig) ParseResult` — pure-Go parser for inline `::`
  chat directives.
- `Directive`, `ParseResult`, `Warning`, `ContextFrame`, `Category` types
  and `ClassifyCommand` helper covering the full public API surface.
- `DefaultAliases()` with built-in aliases (e.g. `ctx` → `context_start`,
  `z` → `zoom`, `bd` → `blog-draft`, `n` → `note`).
- Context scope stack: `context_start` / `context_end` / `zoom` /
  `zoom_out` push and pop frames; action directives carry a
  `ContextRange` resolved from the innermost zoom or outermost context.
- Cascading `config` scopes — deeper scopes override shallower ones at
  directive-emit time.
- Deterministic SHA-256 hashes per directive (command + prompt + context
  window), independent of `Source`, suitable for idempotency keys.
- Graceful degradation on malformed input — unclosed contexts, stray
  `zoom out`, nested `context_start`, and empty `::config` directives all
  emit `Warning` entries instead of failing the parse.
- Unit tests for lexer, parser, config cascade, hash function, and
  graceful-degradation paths; runnable `Example*` functions for
  `pkg.go.dev`.
- `examples/parse` — runnable end-to-end demo program.
- `doc.go` package-level godoc summary.
- `CHANGELOG.md` (this file) and MIT `LICENSE`.

### Notes

- Pre-1.0: the public API may still evolve; breaking changes will be
  called out in this changelog and in a minor-version bump.
- Zero non-stdlib dependencies.
