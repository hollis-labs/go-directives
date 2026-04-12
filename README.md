# Directives

A pure-Go parser for chat directives — inline `::` commands embedded in
conversation text. The library lexes directives, maintains a context scope
stack, resolves cascading config, and computes deterministic SHA-256 hashes
for idempotency. It performs no side effects: it takes text in and emits
structured `Directive` values so the caller can route them to whatever
execution layer it wants.

## Status

Beta — public API is exercised by a full unit test suite (`parser_test.go`,
`lexer_test.go`, `config_test.go`, `hash_test.go`) and ships with a working
CLI under `cmd/directives-parse`, and the tree now includes an MIT
`LICENSE`. There is still no `CHANGELOG.md` or versioned release in the
tree yet.

## Install

```bash
go get github.com/hollis-labs/directives
```

A CLI is included:

```bash
go install github.com/hollis-labs/directives/cmd/directives-parse@latest
```

## Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hollis-labs/directives"
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

The bundled CLI does the same thing end-to-end:

```bash
echo "::blog-draft Write about testing" | directives-parse
directives-parse conversation.md
directives-parse --source "session-abc" --aliases aliases.json < input.txt
```

## API Overview

Declared in `directive.go`:

- `Category` (`CategoryStructure`, `CategoryAction`, `CategoryConfig`, `CategoryMeta`) — classifies each directive; `Category.String()` returns the lowercase name.
- `Directive` — parsed output struct: `Command`, `Prompt`, `Config`, `ContextRange [2]int`, `Hash`, `Line`, `Source`, `Category`.
- `Warning` — non-fatal issue (`Line`, `Message`).
- `ParseResult` — `Directives []Directive` and `Warnings []Warning`.
- `ContextFrame` — a scope entry on the context stack (`Kind`, `Hint`, `Start`).
- `ClassifyCommand(cmd string) Category` — category lookup for a canonical command.

Declared in `parser.go`:

- `ParserConfig` — `Aliases map[string]string` (defaults to `DefaultAliases()` when nil) and `Source string`.
- `Parse(input string, cfg ParserConfig) ParseResult` — lexes and processes the input text and returns all directives plus any warnings.

Declared in `lexer.go`:

- `DefaultAliases() map[string]string` — built-in alias map (`ctx` → `context_start`, `z` → `zoom`, `bd` → `blog-draft`, `n` → `note`, etc.).

Structure commands recognised by the parser: `context_start`, `context_end`,
`zoom`, `zoom_out`. The config command is `config`. The meta command is
`retry`. Any other command is treated as an action and passed through
verbatim to the caller.

## Architecture Notes

The parser is a small pipeline: `lex` scans the input line-by-line for
entries beginning with `::`, resolves aliases against the alias map, and
emits raw `token` values (see `lexer.go`). `Parse` then walks the tokens and
dispatches on `ClassifyCommand`:

- Structure tokens push and pop `ContextFrame` entries on an internal
  `stack` (`stack.go`) and parallel `configStack` scopes (`config.go`).
- Config tokens call `parseConfigPairs` and `set` key/value pairs on the
  innermost scope. Deeper scopes override shallower ones when `merged()` is
  called at directive-emit time.
- Action tokens become `Directive` values. Their `ContextRange` is resolved
  by `resolveContextRange`: innermost zoom → action line, else outermost
  context → action line, else the last-action-line → action line (full
  conversation fallback).
- Meta tokens (currently only `retry`) become `Directive` values with no
  context range, since they are out-of-band parser directives.

Hashes are computed in `hash.go` as
`sha256(command + "\n" + prompt + "\n" + sha256(contextLines))`. The hash is
deterministic across runs and changes when command, prompt, or the context
window changes — but is independent of `Source`, so the same directive in
two different sessions produces the same hash (verified by
`TestParse_HashStability`).

The parser degrades gracefully on malformed input: unclosed contexts, stray
`zoom out` at context level, nested `context_start`, and empty `::config`
directives all emit `Warning` entries rather than failing the parse (see
the `TestParse_GracefulDegradation_*` tests).

## Dependencies

No non-stdlib Go dependencies. `go.mod` declares only the module path
`github.com/hollis-labs/directives` and `go 1.25.0`. No framework-internal
dependencies.

## Testing

```bash
go test ./...
```

Pure unit tests, no external requirements (no database, network, or env
vars). Coverage spans the lexer, parser, config cascade, hash function, and
graceful-degradation paths.

## License

MIT License. See `LICENSE`.
