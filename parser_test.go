package directives

import (
	"strings"
	"testing"
)

func TestParse_SingleAction(t *testing.T) {
	input := "::blog-draft Write about testing"
	r := Parse(input, ParserConfig{Source: "test"})

	if len(r.Directives) != 1 {
		t.Fatalf("expected 1 directive, got %d", len(r.Directives))
	}
	d := r.Directives[0]
	if d.Command != "blog-draft" {
		t.Errorf("command = %q, want %q", d.Command, "blog-draft")
	}
	if d.Prompt != "Write about testing" {
		t.Errorf("prompt = %q, want %q", d.Prompt, "Write about testing")
	}
	if d.Source != "test" {
		t.Errorf("source = %q, want %q", d.Source, "test")
	}
	if d.Category != CategoryAction {
		t.Errorf("category = %v, want %v", d.Category, CategoryAction)
	}
	if d.Hash == "" {
		t.Error("hash is empty")
	}
}

func TestParse_ContextScope(t *testing.T) {
	input := `line 1
::context_start Memory design
line 3 discussion
::blog-draft Write about memory
::context_end`

	r := Parse(input, ParserConfig{})
	if len(r.Directives) != 1 {
		t.Fatalf("expected 1 directive, got %d", len(r.Directives))
	}
	d := r.Directives[0]
	// Context range should start at the context_start line
	if d.ContextRange[0] != 2 {
		t.Errorf("context_range start = %d, want 2", d.ContextRange[0])
	}
	if d.ContextRange[1] != 4 {
		t.Errorf("context_range end = %d, want 4", d.ContextRange[1])
	}
	if len(r.Warnings) != 0 {
		t.Errorf("unexpected warnings: %v", r.Warnings)
	}
}

func TestParse_ZoomScope(t *testing.T) {
	input := `::context_start Topic
discussion
::zoom Subtopic
zoom discussion
::blog-draft Write about subtopic
::zoom out
::context_end`

	r := Parse(input, ParserConfig{})
	if len(r.Directives) != 1 {
		t.Fatalf("expected 1 directive, got %d", len(r.Directives))
	}
	d := r.Directives[0]
	// Context range should start at the zoom line
	if d.ContextRange[0] != 3 {
		t.Errorf("context_range start = %d, want 3", d.ContextRange[0])
	}
}

func TestParse_ConfigCascade(t *testing.T) {
	input := `::config voice=technical, project=carrier
::context_start Sprint
::blog-draft Post about sprint
::zoom Cron details
::config voice=casual
::note Funny bug
::zoom out
::adr Cron decision
::context_end`

	r := Parse(input, ParserConfig{})
	if len(r.Directives) != 3 {
		t.Fatalf("expected 3 directives, got %d", len(r.Directives))
	}

	// blog-draft: voice=technical, project=carrier
	bd := r.Directives[0]
	if bd.Config["voice"] != "technical" {
		t.Errorf("blog-draft voice = %q, want technical", bd.Config["voice"])
	}
	if bd.Config["project"] != "carrier" {
		t.Errorf("blog-draft project = %q, want carrier", bd.Config["project"])
	}

	// note: voice=casual (overridden in zoom), project=carrier (inherited)
	note := r.Directives[1]
	if note.Config["voice"] != "casual" {
		t.Errorf("note voice = %q, want casual", note.Config["voice"])
	}
	if note.Config["project"] != "carrier" {
		t.Errorf("note project = %q, want carrier", note.Config["project"])
	}

	// adr: voice=technical (zoom popped), project=carrier
	adr := r.Directives[2]
	if adr.Config["voice"] != "technical" {
		t.Errorf("adr voice = %q, want technical", adr.Config["voice"])
	}
}

func TestParse_GracefulDegradation_UnclosedContext(t *testing.T) {
	input := `::context_start Unclosed topic
discussion
::blog-draft Post`

	r := Parse(input, ParserConfig{})
	if len(r.Directives) != 1 {
		t.Fatalf("expected 1 directive, got %d", len(r.Directives))
	}
	// Should have a warning about auto-closing
	if len(r.Warnings) == 0 {
		t.Error("expected warning for unclosed context")
	}
	found := false
	for _, w := range r.Warnings {
		if strings.Contains(w.Message, "auto-closed") {
			found = true
		}
	}
	if !found {
		t.Error("expected auto-close warning")
	}
}

func TestParse_GracefulDegradation_UnclosedZoom(t *testing.T) {
	input := `::context_start Topic
::zoom Subtopic
::blog-draft Post
::context_end`

	r := Parse(input, ParserConfig{})
	if len(r.Directives) != 1 {
		t.Fatalf("expected 1 directive, got %d", len(r.Directives))
	}
	// Should warn about auto-closed zoom at context_end
	found := false
	for _, w := range r.Warnings {
		if strings.Contains(w.Message, "zoom") && strings.Contains(w.Message, "auto-closed") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected zoom auto-close warning, got %v", r.Warnings)
	}
}

func TestParse_GracefulDegradation_ExtraZoomOut(t *testing.T) {
	input := "::zoom out"
	r := Parse(input, ParserConfig{})
	if len(r.Directives) != 0 {
		t.Fatalf("expected 0 directives, got %d", len(r.Directives))
	}
	if len(r.Warnings) == 0 {
		t.Error("expected warning for orphan zoom out")
	}
}

func TestParse_GracefulDegradation_ExtraContextEnd(t *testing.T) {
	input := "::context_end"
	r := Parse(input, ParserConfig{})
	if len(r.Warnings) == 0 {
		t.Error("expected warning for orphan context_end")
	}
}

func TestParse_GracefulDegradation_NestedContextStart(t *testing.T) {
	input := `::context_start First
::blog-draft Post 1
::context_start Second
::blog-draft Post 2
::context_end`

	r := Parse(input, ParserConfig{})
	if len(r.Directives) != 2 {
		t.Fatalf("expected 2 directives, got %d", len(r.Directives))
	}
	// Should warn about auto-closing first context
	found := false
	for _, w := range r.Warnings {
		if strings.Contains(w.Message, "auto-closed previous context_start") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected auto-close warning for nested context_start, got %v", r.Warnings)
	}
}

func TestParse_ZoomOutAtContextLevel(t *testing.T) {
	input := `::context_start Topic
::zoom out
::context_end`

	r := Parse(input, ParserConfig{})
	found := false
	for _, w := range r.Warnings {
		if strings.Contains(w.Message, "zoom out at context level") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected warning for zoom out at context level, got %v", r.Warnings)
	}
}

func TestParse_MetaDirective(t *testing.T) {
	input := "::retry abc123"
	r := Parse(input, ParserConfig{})
	if len(r.Directives) != 1 {
		t.Fatalf("expected 1 directive, got %d", len(r.Directives))
	}
	d := r.Directives[0]
	if d.Command != "retry" {
		t.Errorf("command = %q, want retry", d.Command)
	}
	if d.Prompt != "abc123" {
		t.Errorf("prompt = %q, want abc123", d.Prompt)
	}
	if d.Category != CategoryMeta {
		t.Errorf("category = %v, want %v", d.Category, CategoryMeta)
	}
}

func TestParse_EmptyConfig(t *testing.T) {
	input := "::config\n::blog-draft Post"
	r := Parse(input, ParserConfig{})
	if len(r.Directives) != 1 {
		t.Fatalf("expected 1 directive, got %d", len(r.Directives))
	}
	// Should warn about empty config
	found := false
	for _, w := range r.Warnings {
		if strings.Contains(w.Message, "no valid key=value") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected warning for empty config, got %v", r.Warnings)
	}
}

func TestParse_FullConversationFallback(t *testing.T) {
	input := `line 1 discussion
line 2 more talk
::blog-draft Post about everything`

	r := Parse(input, ParserConfig{})
	d := r.Directives[0]
	// No structure directives, so context range is full conversation
	if d.ContextRange[0] != 1 {
		t.Errorf("context_range start = %d, want 1", d.ContextRange[0])
	}
	if d.ContextRange[1] != 3 {
		t.Errorf("context_range end = %d, want 3", d.ContextRange[1])
	}
}

func TestParse_SequentialActions(t *testing.T) {
	input := `discussion
::blog-draft First post
more discussion
::blog-draft Second post`

	r := Parse(input, ParserConfig{})
	if len(r.Directives) != 2 {
		t.Fatalf("expected 2 directives, got %d", len(r.Directives))
	}

	// First action: context from line 1 to line 2
	if r.Directives[0].ContextRange[0] != 1 {
		t.Errorf("first context start = %d, want 1", r.Directives[0].ContextRange[0])
	}

	// Second action: context from line 2 (last action) to line 4
	if r.Directives[1].ContextRange[0] != 2 {
		t.Errorf("second context start = %d, want 2", r.Directives[1].ContextRange[0])
	}
	if r.Directives[1].ContextRange[1] != 4 {
		t.Errorf("second context end = %d, want 4", r.Directives[1].ContextRange[1])
	}
}

func TestParse_HashStability(t *testing.T) {
	input := "::blog-draft Write about memory"
	r1 := Parse(input, ParserConfig{Source: "s1"})
	r2 := Parse(input, ParserConfig{Source: "s2"})

	if r1.Directives[0].Hash != r2.Directives[0].Hash {
		t.Errorf("hashes differ for same input: %s vs %s",
			r1.Directives[0].Hash, r2.Directives[0].Hash)
	}
}

func TestParse_HashDiffersWithDifferentContext(t *testing.T) {
	input1 := "context A\n::blog-draft Post"
	input2 := "context B\n::blog-draft Post"

	r1 := Parse(input1, ParserConfig{})
	r2 := Parse(input2, ParserConfig{})

	if r1.Directives[0].Hash == r2.Directives[0].Hash {
		t.Error("hashes should differ for different context")
	}
}

func TestParse_HashDiffersWithDifferentPrompt(t *testing.T) {
	input1 := "::blog-draft Post about A"
	input2 := "::blog-draft Post about B"

	r1 := Parse(input1, ParserConfig{})
	r2 := Parse(input2, ParserConfig{})

	if r1.Directives[0].Hash == r2.Directives[0].Hash {
		t.Error("hashes should differ for different prompts")
	}
}

func TestParse_AliasesInFullPipeline(t *testing.T) {
	input := "::bd Quick post\n::n Remember this"
	r := Parse(input, ParserConfig{})

	if len(r.Directives) != 2 {
		t.Fatalf("expected 2 directives, got %d", len(r.Directives))
	}
	if r.Directives[0].Command != "blog-draft" {
		t.Errorf("first command = %q, want blog-draft", r.Directives[0].Command)
	}
	if r.Directives[1].Command != "note" {
		t.Errorf("second command = %q, want note", r.Directives[1].Command)
	}
}

func TestClassifyCommand(t *testing.T) {
	tests := []struct {
		cmd  string
		want Category
	}{
		{"context_start", CategoryStructure},
		{"context_end", CategoryStructure},
		{"zoom", CategoryStructure},
		{"zoom_out", CategoryStructure},
		{"config", CategoryConfig},
		{"retry", CategoryMeta},
		{"blog-draft", CategoryAction},
		{"adr", CategoryAction},
		{"note", CategoryAction},
		{"document-feature", CategoryAction},
		{"custom-thing", CategoryAction},
	}
	for _, tt := range tests {
		t.Run(tt.cmd, func(t *testing.T) {
			got := ClassifyCommand(tt.cmd)
			if got != tt.want {
				t.Errorf("ClassifyCommand(%q) = %v, want %v", tt.cmd, got, tt.want)
			}
		})
	}
}
