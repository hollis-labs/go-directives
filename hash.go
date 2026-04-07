package directives

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// computeHash produces a deterministic SHA-256 hex string for a directive.
//
// Hash = sha256(command + "\n" + prompt + "\n" + contextWindowHash)
//
// contextWindowHash is sha256 of the lines within the directive's context range.
func computeHash(command, prompt string, contextLines []string) string {
	contextHash := hashLines(contextLines)
	payload := command + "\n" + prompt + "\n" + contextHash
	h := sha256.Sum256([]byte(payload))
	return fmt.Sprintf("%x", h)
}

// hashLines produces a SHA-256 hex string for a slice of lines.
func hashLines(lines []string) string {
	joined := strings.Join(lines, "\n")
	h := sha256.Sum256([]byte(joined))
	return fmt.Sprintf("%x", h)
}
