// Package main is a runnable example that demonstrates parsing a small block
// of chat directives with go-directives, then printing each parsed Directive
// and any non-fatal Warnings.
//
// Run from the module root:
//
//	go run ./examples/parse
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
::zoom Cron details
::note Funny bug
::zoom out
::adr Cron decision
::context_end`

	result := directives.Parse(input, directives.ParserConfig{
		Source: "session-example",
	})

	for _, w := range result.Warnings {
		fmt.Fprintf(os.Stderr, "warning: line %d: %s\n", w.Line, w.Message)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result.Directives); err != nil {
		fmt.Fprintf(os.Stderr, "encode: %v\n", err)
		os.Exit(1)
	}
}
