package directives

import (
	"testing"
)

func TestLex_BasicDirectives(t *testing.T) {
	input := "hello\n::blog-draft Write about X\ngoodbye"
	tokens := lex(input, DefaultAliases())

	if len(tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(tokens))
	}
	if tokens[0].Command != "blog-draft" {
		t.Errorf("command = %q, want %q", tokens[0].Command, "blog-draft")
	}
	if tokens[0].Prompt != "Write about X" {
		t.Errorf("prompt = %q, want %q", tokens[0].Prompt, "Write about X")
	}
	if tokens[0].Line != 2 {
		t.Errorf("line = %d, want 2", tokens[0].Line)
	}
}

func TestLex_AliasResolution(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantCmd string
	}{
		{"ctx -> context_start", "::ctx Topic", "context_start"},
		{"ctx_end -> context_end", "::ctx_end", "context_end"},
		{"/ctx -> context_end", "::/ctx", "context_end"},
		{"z -> zoom", "::z Subtopic", "zoom"},
		{"zo -> zoom_out", "::zo", "zoom_out"},
		{"bd -> blog-draft", "::bd My post", "blog-draft"},
		{"df -> document-feature", "::df Feature X", "document-feature"},
		{"n -> note", "::n Remember this", "note"},
		{"no alias", "::adr Some decision", "adr"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lex(tt.input, DefaultAliases())
			if len(tokens) != 1 {
				t.Fatalf("expected 1 token, got %d", len(tokens))
			}
			if tokens[0].Command != tt.wantCmd {
				t.Errorf("command = %q, want %q", tokens[0].Command, tt.wantCmd)
			}
		})
	}
}

func TestLex_ZoomOut(t *testing.T) {
	input := "::zoom out"
	tokens := lex(input, DefaultAliases())
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(tokens))
	}
	if tokens[0].Command != "zoom_out" {
		t.Errorf("command = %q, want %q", tokens[0].Command, "zoom_out")
	}
}

func TestLex_EmptyDirective(t *testing.T) {
	input := "::\nsome text"
	tokens := lex(input, DefaultAliases())
	if len(tokens) != 0 {
		t.Fatalf("expected 0 tokens for empty ::, got %d", len(tokens))
	}
}

func TestLex_IndentedDirective(t *testing.T) {
	input := "  ::note Indented note"
	tokens := lex(input, DefaultAliases())
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(tokens))
	}
	if tokens[0].Command != "note" {
		t.Errorf("command = %q, want %q", tokens[0].Command, "note")
	}
}

func TestLex_MultipleDirectives(t *testing.T) {
	input := "::context_start Topic\nSome text\n::blog-draft Post\n::context_end"
	tokens := lex(input, DefaultAliases())
	if len(tokens) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(tokens))
	}
	want := []string{"context_start", "blog-draft", "context_end"}
	for i, w := range want {
		if tokens[i].Command != w {
			t.Errorf("token[%d].Command = %q, want %q", i, tokens[i].Command, w)
		}
	}
}

func TestLex_CustomAliases(t *testing.T) {
	aliases := map[string]string{
		"b": "blog-draft",
	}
	input := "::b Quick post"
	tokens := lex(input, aliases)
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(tokens))
	}
	if tokens[0].Command != "blog-draft" {
		t.Errorf("command = %q, want %q", tokens[0].Command, "blog-draft")
	}
}

func TestLex_NoPrompt(t *testing.T) {
	input := "::context_end"
	tokens := lex(input, DefaultAliases())
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(tokens))
	}
	if tokens[0].Prompt != "" {
		t.Errorf("prompt = %q, want empty", tokens[0].Prompt)
	}
}
