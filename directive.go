package directives

// Category classifies a directive.
type Category int

const (
	CategoryStructure Category = iota
	CategoryAction
	CategoryConfig
	CategoryMeta
)

func (c Category) String() string {
	switch c {
	case CategoryStructure:
		return "structure"
	case CategoryAction:
		return "action"
	case CategoryConfig:
		return "config"
	case CategoryMeta:
		return "meta"
	default:
		return "unknown"
	}
}

// Directive is the parsed output for a single :: command.
type Directive struct {
	Command      string            `json:"command"`       // canonical name (aliases resolved)
	Prompt       string            `json:"prompt"`        // freeform text after command
	Config       map[string]string `json:"config"`        // merged config (scope cascade applied)
	ContextRange [2]int            `json:"context_range"` // start/end line indices
	Hash         string            `json:"hash"`          // deterministic SHA-256
	Line         int               `json:"line"`          // source line number (1-based)
	Source       string            `json:"source"`        // source identifier
	Category     Category          `json:"category"`      // directive category
}

// Warning represents a non-fatal issue found during parsing.
type Warning struct {
	Line    int    `json:"line"`
	Message string `json:"message"`
}

// ParseResult holds all output from a parse pass.
type ParseResult struct {
	Directives []Directive `json:"directives"`
	Warnings   []Warning   `json:"warnings"`
}

// ContextFrame represents a scope on the context stack.
type ContextFrame struct {
	Kind  string // "context" or "zoom"
	Hint  string // topic hint from the directive prompt
	Start int    // line index where this frame opened
}

// structureCommands are directives that control parser scope.
var structureCommands = map[string]bool{
	"context_start": true,
	"context_end":   true,
	"zoom":          true,
	"zoom_out":      true,
}

// configCommand is the config directive.
const configCommand = "config"

// metaCommands are directives that control parser behavior.
var metaCommands = map[string]bool{
	"retry": true,
}

// ClassifyCommand returns the category for a canonical command name.
func ClassifyCommand(cmd string) Category {
	if structureCommands[cmd] {
		return CategoryStructure
	}
	if cmd == configCommand {
		return CategoryConfig
	}
	if metaCommands[cmd] {
		return CategoryMeta
	}
	return CategoryAction
}
