# Audit — go-directives

**Audited:** 2026-04-09
**Auditor:** general-purpose subagent (BOOT_STANDARDIZATION audit)
**Path:** libs/go-directives
**Kind:** lib

## Summary

`go-directives` is a self-contained pure-Go library that parses inline `::`
chat directives into structured values, with a companion CLI under
`cmd/directives-parse`. The code is well-documented, well-tested, and has
zero non-stdlib dependencies. The apply pass added the missing MIT
`LICENSE` and example coverage; `CHANGELOG.md` is still absent. The public
API is clean and consistent; no API-level blockers were found.

## Checklist

| # | Check | Status | Notes |
|---|---|---|---|
| 1 | `go.mod` present | pass | `github.com/hollis-labs/directives`, go 1.25.0 |
| 2 | `README.md` present (before this audit) | fail | none in tree; new standardized README written by this audit |
| 3 | `LICENSE` present | pass | MIT LICENSE added in apply session |
| 4 | `doc.go` with `// Package X ...` godoc comment | pass | no `doc.go`, but `directive.go` has a full `// Package directives ...` godoc block at the top |
| 5 | Module path matches intended repo layout | pass | resolved by apply-session decision D2: keep the standalone repo's canonical module path `github.com/hollis-labs/directives`; no rewrite to match the local framework folder |
| 6 | README has standard sections (title, desc, install, usage, API, examples) | pass | new README written to template; previous state was "missing" |
| 7 | Tests exist (`*_test.go`) | pass | 4 test files: `parser_test.go`, `lexer_test.go`, `config_test.go`, `hash_test.go`; coverage not measured in this read-only pass |
| 8 | Examples (`example_test.go` or `examples/`) | pass | `example_test.go` added with `ExampleParse` and `ExampleParse_configCascade` |
| 9 | State/session files NOT misclassified as library docs | pass | no `.agentrc/`, `BOOT.md`, `CLAUDE.md`, `bootstrap.md`, `boot-prompt.md`, or `boot/*.md` files present |
| 10 | Public API sanity: errors typed/sentinel, context.Context first arg | pass | `Parse` returns `ParseResult` (with embedded `Warnings`) and does no I/O — no error return and no context param needed; the library is intentionally side-effect-free, so the context/error conventions do not apply |
| 11 | `CHANGELOG.md` present (nice to have) | fail | absent |
| 12 | No circular/suspicious deps on other framework libs | pass | zero framework-internal deps; zero external deps |

## Findings — Required Fixes

1. **What:** No `LICENSE` file.
   **Why:** A Go library without a license cannot be legally consumed or
   redistributed; this blocks any external publish and is ambiguous even
   for internal consumers.
   **Suggested fix:** Add a top-level `LICENSE` file matching the
   framework's standard license (MIT).  
   **Status:** Resolved in apply session with [`LICENSE`](./LICENSE).

2. **What:** Module path is `github.com/hollis-labs/directives`, but the
   code lives at `framework/libs/go-directives` under a different org.
   **Why:** If this module is ever `go get`'d or vendored, the import path
   must resolve to a reachable repo. If `hollis-labs` is a legacy or
   personal org, this will break cross-lib references once framework libs
   start importing each other.
   **Suggested fix:** Decide the canonical publish org (likely the same
   namespace the rest of the framework uses) and update `go.mod` plus the
   import in `cmd/directives-parse/main.go` in one commit.  
   **Status:** Resolved by apply-session decision D2. The standalone repo
   keeps its canonical module path `github.com/hollis-labs/directives`; no
   rename was required for this framework-local standardization pass.

3. **What:** No `example_test.go`.
   **Why:** `go doc` and pkg.go.dev surface `Example*` functions as
   runnable documentation. Without one, discoverability of the public API
   through standard Go tooling is reduced.
   **Suggested fix:** Add an `example_test.go` with at least
   `ExampleParse` (single action) and `ExampleParse_configCascade`
   (matches the existing `TestParse_ConfigCascade` scenario).  
   **Status:** Resolved in apply session with [`example_test.go`](./example_test.go).

## Findings — Nice-to-Have

1. **What:** No `CHANGELOG.md`.
   **Why:** Downstream consumers can't track API churn.
   **Suggested fix:** Seed a `CHANGELOG.md` at the current state as
   "unreleased / 0.1.0" and start tracking changes going forward.

2. **What:** Package-level godoc is in `directive.go` rather than a
   dedicated `doc.go`.
   **Why:** Convention in most Go projects is to keep the package comment
   in `doc.go` so it is easy to find and does not get lost when
   `directive.go` is edited. This is stylistic only.
   **Suggested fix:** Move the leading `// Package directives ...` block
   from `directive.go` into a new `doc.go`.

3. **What:** Internal helper `lines_count` uses snake_case, which is
   non-idiomatic Go.
   **Why:** Style consistency; `golint`/`revive` will flag it.
   **Suggested fix:** Rename to `linesCount` or inline at the single call
   site in `parser.go`.

4. **What:** The CLI binary `cmd/directives-parse/main.go` is undocumented
   outside its own package comment — it is not mentioned in any existing
   top-level doc (because there is no top-level doc).
   **Why:** Users reading the new README will see it, but a dedicated
   `cmd/directives-parse/README.md` or an extended `Usage` section in the
   package doc would make the CLI discoverable via `go doc`.
   **Suggested fix:** Optional — consider a short CLI README or extend
   the CLI's package comment with a full flag table.

## Prior Documentation

- **Pre-existing README.md:** present and updated in place for the new LICENSE status.
- **Pre-existing README.original.md:** none (no rename needed).
- **Pre-existing LICENSE / CHANGELOG:** none.
- **Post-apply LICENSE:** [`LICENSE`](./LICENSE) added with MIT text.
- **Pre-existing `docs/` subfolder:** none.
- **Session/state files noted but excluded from library docs:** none
  found. The lib is clean of `.agentrc/`, `BOOT.md`, `CLAUDE.md`,
  `bootstrap.md`, `boot-prompt.md`, or `boot/*.md`.
- **Other notes:** The lib is its own nested git repo (has a `.git/`
  directory). `.gitignore` contains only `directives-parse` (the CLI
  binary output).

## Public API Snapshot

### `directive.go`

- `type Category int`
  - constants: `CategoryStructure`, `CategoryAction`, `CategoryConfig`, `CategoryMeta`
  - method: `(Category) String() string`
- `type Directive struct` — fields: `Command`, `Prompt`, `Config`, `ContextRange [2]int`, `Hash`, `Line`, `Source`, `Category`
- `type Warning struct` — fields: `Line`, `Message`
- `type ParseResult struct` — fields: `Directives []Directive`, `Warnings []Warning`
- `type ContextFrame struct` — fields: `Kind`, `Hint`, `Start`
- `func ClassifyCommand(cmd string) Category`

### `parser.go`

- `type ParserConfig struct` — fields: `Aliases map[string]string`, `Source string`
- `func Parse(input string, cfg ParserConfig) ParseResult`

### `lexer.go`

- `func DefaultAliases() map[string]string`

### `config.go`

- No exported symbols (internal `configStack`, `parseConfigPairs`).

### `stack.go`

- No exported symbols (internal `stack` type).

### `hash.go`

- No exported symbols (internal `computeHash`, `hashLines`).

### `cmd/directives-parse/main.go`

- `package main` — CLI entry point. Reads stdin or a file arg and writes
  `ParseResult` as JSON to stdout; warnings to stderr. Flags: `--source`,
  `--aliases`.

## Open Questions

1. Is `github.com/hollis-labs/directives` the intended publish path, or
   should this be re-homed under the framework's canonical module
   namespace? (Affects Finding #2 above.)
2. What license should the framework standardize on for in-tree libs?
   (Affects Finding #1 above.)
3. Is this library already consumed by another framework lib/plugin? The
   audit found no internal consumers from the code alone, but a grep
   across the broader framework tree was out of scope for this per-lib
   pass.
