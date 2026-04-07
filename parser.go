package directives

import (
	"fmt"
	"strings"
)

// ParserConfig controls parser behavior.
type ParserConfig struct {
	// Aliases maps short forms to canonical command names.
	// If nil, DefaultAliases() is used.
	Aliases map[string]string

	// Source is an identifier for the input (conversation ID, file path).
	// Attached to each emitted Directive.
	Source string
}

// Parse lexes and processes directives from input text.
func Parse(input string, cfg ParserConfig) ParseResult {
	aliases := cfg.Aliases
	if aliases == nil {
		aliases = DefaultAliases()
	}

	lines := strings.Split(input, "\n")
	tokens := lex(input, aliases)

	var (
		result    ParseResult
		ctxStack  stack
		cfgStack  configStack
		lastAction int // line after the last completed action (for full-conversation fallback)
	)

	// Start with a root config scope
	cfgStack.push()

	for _, tok := range tokens {
		switch ClassifyCommand(tok.Command) {
		case CategoryStructure:
			handleStructure(tok, lines, &ctxStack, &cfgStack, &result)

		case CategoryConfig:
			handleConfig(tok, &cfgStack, &result)

		case CategoryAction:
			d := buildAction(tok, lines, &ctxStack, &cfgStack, cfg.Source, lastAction)
			result.Directives = append(result.Directives, d)
			lastAction = tok.Line

		case CategoryMeta:
			d := buildMeta(tok, &cfgStack, cfg.Source)
			result.Directives = append(result.Directives, d)
		}
	}

	// Auto-close any remaining open frames
	if n := ctxStack.clear(); n > 0 {
		result.Warnings = append(result.Warnings, Warning{
			Line:    lines_count(lines),
			Message: fmt.Sprintf("auto-closed %d unclosed scope(s) at end of input", n),
		})
	}

	return result
}

// handleStructure processes context_start, context_end, zoom, zoom_out.
func handleStructure(tok token, lines []string, s *stack, cs *configStack, r *ParseResult) {
	switch tok.Command {
	case "context_start":
		// If a context is already open, auto-close it
		if !s.empty() && hasContext(s) {
			zooms := s.popZooms()
			if zooms > 0 {
				r.Warnings = append(r.Warnings, Warning{
					Line:    tok.Line,
					Message: fmt.Sprintf("auto-closed %d zoom(s) at new context_start", zooms),
				})
			}
			// Pop the context frame itself
			if f, ok := s.peek(); ok && f.Kind == "context" {
				s.pop()
				cs.pop()
				r.Warnings = append(r.Warnings, Warning{
					Line:    tok.Line,
					Message: "auto-closed previous context_start at new context_start",
				})
			}
		}
		s.push(ContextFrame{Kind: "context", Hint: tok.Prompt, Start: tok.Line})
		cs.push()

	case "context_end":
		if s.empty() || !hasContext(s) {
			r.Warnings = append(r.Warnings, Warning{
				Line:    tok.Line,
				Message: "context_end with no open context — ignored",
			})
			return
		}
		// Pop all zooms within this context
		zooms := s.popZooms()
		if zooms > 0 {
			r.Warnings = append(r.Warnings, Warning{
				Line:    tok.Line,
				Message: fmt.Sprintf("auto-closed %d unclosed zoom(s) at context_end", zooms),
			})
		}
		// Pop the context frame
		s.pop()
		cs.pop()

	case "zoom":
		s.push(ContextFrame{Kind: "zoom", Hint: tok.Prompt, Start: tok.Line})
		cs.push()

	case "zoom_out":
		if s.empty() {
			r.Warnings = append(r.Warnings, Warning{
				Line:    tok.Line,
				Message: "zoom out with empty stack — ignored",
			})
			return
		}
		top, _ := s.peek()
		if top.Kind != "zoom" {
			r.Warnings = append(r.Warnings, Warning{
				Line:    tok.Line,
				Message: "zoom out at context level — ignored",
			})
			return
		}
		s.pop()
		cs.pop()
	}
}

// handleConfig processes ::config directives.
func handleConfig(tok token, cs *configStack, r *ParseResult) {
	pairs := parseConfigPairs(tok.Prompt)
	if len(pairs) == 0 {
		r.Warnings = append(r.Warnings, Warning{
			Line:    tok.Line,
			Message: "config directive with no valid key=value pairs",
		})
		return
	}
	for k, v := range pairs {
		cs.set(k, v)
	}
}

// buildAction constructs a Directive for an action command.
func buildAction(tok token, lines []string, s *stack, cs *configStack, source string, lastAction int) Directive {
	// Determine context range
	start, end := resolveContextRange(s, tok.Line, lastAction)

	// Extract context lines for hashing
	contextLines := extractLines(lines, start, end)

	return Directive{
		Command:      tok.Command,
		Prompt:       tok.Prompt,
		Config:       cs.merged(),
		ContextRange: [2]int{start, end},
		Hash:         computeHash(tok.Command, tok.Prompt, contextLines),
		Line:         tok.Line,
		Source:       source,
		Category:     CategoryAction,
	}
}

// buildMeta constructs a Directive for a meta command (e.g., ::retry).
func buildMeta(tok token, cs *configStack, source string) Directive {
	return Directive{
		Command:  tok.Command,
		Prompt:   tok.Prompt,
		Config:   cs.merged(),
		Hash:     computeHash(tok.Command, tok.Prompt, nil),
		Line:     tok.Line,
		Source:   source,
		Category: CategoryMeta,
	}
}

// resolveContextRange determines the start and end line numbers for an action's context.
//
// Priority:
// 1. Innermost zoom frame start → action line
// 2. Context frame start → action line
// 3. Last action line (or 1) → action line
func resolveContextRange(s *stack, actionLine int, lastAction int) (int, int) {
	if !s.empty() {
		return s.currentStart(), actionLine
	}
	// Full conversation fallback
	start := 1
	if lastAction > 0 {
		start = lastAction
	}
	return start, actionLine
}

// extractLines returns lines[start-1:end] (1-based inclusive to 0-based slice).
func extractLines(lines []string, start, end int) []string {
	if start < 1 {
		start = 1
	}
	if end > len(lines) {
		end = len(lines)
	}
	if start > end {
		return nil
	}
	return lines[start-1 : end]
}

// hasContext returns true if the stack contains at least one context frame.
func hasContext(s *stack) bool {
	for _, f := range s.frames {
		if f.Kind == "context" {
			return true
		}
	}
	return false
}

// lines_count returns the number of lines.
func lines_count(lines []string) int {
	return len(lines)
}
