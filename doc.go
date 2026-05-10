// Package directives implements a parser for chat directives — inline ::
// commands embedded in conversation text. It lexes directives, maintains a
// context-scope stack, resolves cascading config, and computes deterministic
// SHA-256 hashes for idempotency.
//
// The parser is a pure library with no side effects. It takes text in and
// emits structured Directive values; all execution logic lives in the
// consumer.
//
// # Quick start
//
//	result := directives.Parse(input, directives.ParserConfig{Source: "session-1"})
//	for _, d := range result.Directives {
//	    // route d.Command / d.Prompt / d.Config / d.ContextRange to your runtime
//	}
//
// See the examples_test.go file in this package and the runnable program
// under examples/parse for end-to-end usage.
package directives
