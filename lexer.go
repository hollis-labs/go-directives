package directives

import (
	"strings"
)

// token represents a raw lexed directive before semantic processing.
type token struct {
	Command string // raw command text (before alias resolution)
	Prompt  string // freeform text after command
	Line    int    // 1-based line number
}

// DefaultAliases returns the MVP alias map.
func DefaultAliases() map[string]string {
	return map[string]string{
		"ctx":     "context_start",
		"ctx_end": "context_end",
		"/ctx":    "context_end",
		"z":       "zoom",
		"zo":      "zoom_out",
		"bd":      "blog-draft",
		"df":      "document-feature",
		"n":       "note",
	}
}

// lex scans input text for :: prefixed lines and returns tokens.
// Aliases are resolved during lexing so downstream code only sees canonical names.
func lex(input string, aliases map[string]string) []token {
	lines := strings.Split(input, "\n")
	var tokens []token

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		if !strings.HasPrefix(trimmed, "::") {
			continue
		}

		// Strip the :: prefix
		rest := trimmed[2:]
		if rest == "" {
			continue
		}

		// Handle "zoom out" as a special two-word command
		if strings.HasPrefix(rest, "zoom out") || strings.HasPrefix(rest, "zo ") {
			// Check if it's actually "zoom out" (not "zoom Something")
			if strings.HasPrefix(rest, "zoom out") {
				prompt := strings.TrimSpace(rest[len("zoom out"):])
				cmd := resolveAlias("zoom_out", aliases)
				tokens = append(tokens, token{Command: cmd, Prompt: prompt, Line: lineNum})
				continue
			}
		}

		// Split into command and prompt
		cmd, prompt := splitCommandPrompt(rest)

		// Resolve alias
		cmd = resolveAlias(cmd, aliases)

		tokens = append(tokens, token{Command: cmd, Prompt: prompt, Line: lineNum})
	}

	return tokens
}

// splitCommandPrompt splits "command rest of line" into (command, prompt).
func splitCommandPrompt(s string) (string, string) {
	// Find the first space
	idx := strings.IndexByte(s, ' ')
	if idx == -1 {
		return s, ""
	}
	return s[:idx], strings.TrimSpace(s[idx+1:])
}

// resolveAlias maps an alias to its canonical form. If no alias exists,
// the original command is returned unchanged.
func resolveAlias(cmd string, aliases map[string]string) string {
	if canonical, ok := aliases[cmd]; ok {
		return canonical
	}
	return cmd
}
