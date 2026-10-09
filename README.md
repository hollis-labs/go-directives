# go-directives

## Maintenance moved to `github.com/hollis-labs/libs/ui-go`

This standalone repository is retired. Maintained source and documentation are
in [github.com/hollis-labs/libs/ui-go/directives](https://github.com/hollis-labs/libs/tree/ui-go%2Fv0.1.0/ui-go/directives), released in **`ui-go/v0.1.0`**.
Install the replacement module:

```sh
go get github.com/hollis-labs/libs/ui-go@v0.1.0
```

Replace the `github.com/hollis-labs/go-directives` import prefix with
`github.com/hollis-labs/libs/ui-go/directives`, retaining the package subpath. Review the replacement documentation
for any API changes before migrating. Existing standalone tags and history remain
available; old module pins do not automatically redirect to the new module.

The documentation below describes historical standalone usage.

[![Go Reference](https://pkg.go.dev/badge/github.com/hollis-labs/go-directives.svg)](https://pkg.go.dev/github.com/hollis-labs/go-directives)

A pure-Go parser for **chat directives** — inline `::` commands embedded in
conversation text. The library lexes directives, maintains a context-scope
stack, resolves cascading config, and computes deterministic SHA-256 hashes
for idempotency. It performs no side effects: it takes text in and emits
structured `Directive` values so the caller can route them to whatever
execution layer it wants.

## Status

**Pre-1.0 (v0.1.x).** The public API is exercised by a full unit test
suite (`parser_test.go`, `lexer_test.go`, `config_test.go`,
`hash_test.go`) and runnable `Example*` functions for `pkg.go.dev`. The
shape of the API is unlikely to change in incompatible ways before 1.0,
but breaking changes — if any — will be called out in
[`CHANGELOG.md`](./CHANGELOG.md) and bumped via a minor version.

## Install

```bash
go get github.com/hollis-labs/go-directives
```

## Quickstart

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	directives "github.com/hollis-labs/go-directives"
)

func main() {
	input := `::config voice=technical, project=carrier
::context_start Sprint planning
discussion about cron jobs
::blog-draft Post about the sprint
::context_end`

	result := directives.Parse(input, directives.ParserConfig{
		Source: "session-abc",
	})

	for _, w := range result.Warnings {
		fmt.Fprintf(os.Stderr, "warning: line %d: %s\n", w.Line, w.Message)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(result.Directives)
}
```

A more involved example that exercises the config cascade and zoom scopes
lives in [`examples/parse`](./examples/parse) and can be run directly:

```bash
go run ./examples/parse
```

## Documentation

Full API docs on godoc:

- <https://pkg.go.dev/github.com/hollis-labs/go-directives>

The package surface is small and is summarised below; see godoc for the
authoritative reference.

### Public API at a glance

Declared in `directive.go`:

- `Category` — `CategoryStructure`, `CategoryAction`, `CategoryConfig`,
  `CategoryMeta`; `Category.String()` returns the lowercase name.
- `Directive` — parsed output: `Command`, `Prompt`, `Config`,
  `ContextRange [2]int`, `Hash`, `Line`, `Source`, `Category`.
- `Warning` — non-fatal issue (`Line`, `Message`).
- `ParseResult` — `Directives []Directive` and `Warnings []Warning`.
- `ContextFrame` — a scope entry on the context stack (`Kind`, `Hint`,
  `Start`).
- `ClassifyCommand(cmd string) Category` — category lookup for a
  canonical command.

Declared in `parser.go`:

- `ParserConfig` — `Aliases map[string]string` (defaults to
  `DefaultAliases()` when nil) and `Source string`.
- `Parse(input string, cfg ParserConfig) ParseResult` — lexes and
  processes the input text and returns all directives plus any
  warnings.

Declared in `lexer.go`:

- `DefaultAliases() map[string]string` — built-in alias map (`ctx` →
  `context_start`, `z` → `zoom`, `bd` → `blog-draft`, `n` → `note`,
  etc.).

Structure commands recognised by the parser: `context_start`,
`context_end`, `zoom`, `zoom_out`. The config command is `config`. The
meta command is `retry`. Any other command is treated as an action and
passed through verbatim to the caller.

## Behaviour notes

- **Hashes** are computed as
  `sha256(command + "\n" + prompt + "\n" + sha256(contextLines))` and are
  deterministic across runs — independent of `Source`, so the same
  directive in two different sessions produces the same hash (verified by
  `TestParse_HashStability`).
- **Graceful degradation:** unclosed contexts, stray `zoom out` at
  context level, nested `context_start`, and empty `::config` directives
  all emit `Warning` entries rather than failing the parse (see the
  `TestParse_GracefulDegradation_*` tests).

## Dependencies

No non-stdlib Go dependencies. `go.mod` declares only the module path
`github.com/hollis-labs/go-directives` and the current Go directive.

## Testing

```bash
go test ./...
go test -race ./...
```

Pure unit tests, no external requirements (no database, network, or env
vars). Coverage spans the lexer, parser, config cascade, hash function,
and graceful-degradation paths.

## License

MIT License. See [`LICENSE`](./LICENSE).
