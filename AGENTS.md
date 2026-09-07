# go-directives

A pure parser for chat directives — inline `::` commands embedded in
conversation text. It lexes directives, maintains a context-scope stack,
resolves cascading config and computes a deterministic hash per directive. It
performs no side effects and executes nothing: text in, structured `Directive`
values out, routing left to the caller.

## Start Here

- `README.md` summarizes the public API; godoc is the authoritative reference.
- `directive.go` declares `Directive`, `Category`, `Warning` and `ParseResult`.
- `lexer.go` recognizes directive lines; `parser.go` owns the parse pipeline
  and `ParserConfig`.
- `stack.go` and `config.go` own context scoping and the config cascade.
- `hash.go` owns the idempotency hash.
- `examples/parse/main.go` exercises the cascade and zoom scopes end to end.

## Commands

```bash
gofmt -l .
go vet ./...
go test -race -count=1 ./...
go run ./examples/parse
```

There is no CI workflow in this repo.

## Boundaries

The module has no dependencies beyond the standard library. Keeping it that way
is most of why it is a separate library.

The hash is a contract, not an implementation detail. It is
`sha256(command + "\n" + prompt + "\n" + sha256(context lines))`, and callers
use it for idempotency — so changing the payload shape, the separator or the
ordering silently invalidates every stored hash downstream. `TestComputeHash_Deterministic`
and the three `DiffersOn*` tests pin which inputs must change it.

Malformed input degrades gracefully rather than failing the parse: an unmatched
`::context_end` or an extra zoom-out becomes a `Warning`, not an error
(`TestParse_GracefulDegradation_ExtraContextEnd`,
`TestParse_GracefulDegradation_ExtraZoomOut`). Directives arrive inside human
conversation text, so a strict parser would reject real messages.
