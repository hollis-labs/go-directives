package directives

import (
	"testing"
)

func TestComputeHash_Deterministic(t *testing.T) {
	h1 := computeHash("blog-draft", "Write about X", []string{"context line 1", "context line 2"})
	h2 := computeHash("blog-draft", "Write about X", []string{"context line 1", "context line 2"})
	if h1 != h2 {
		t.Errorf("hashes differ: %s vs %s", h1, h2)
	}
}

func TestComputeHash_DiffersOnCommand(t *testing.T) {
	h1 := computeHash("blog-draft", "prompt", []string{"ctx"})
	h2 := computeHash("note", "prompt", []string{"ctx"})
	if h1 == h2 {
		t.Error("hashes should differ for different commands")
	}
}

func TestComputeHash_DiffersOnPrompt(t *testing.T) {
	h1 := computeHash("blog-draft", "prompt A", []string{"ctx"})
	h2 := computeHash("blog-draft", "prompt B", []string{"ctx"})
	if h1 == h2 {
		t.Error("hashes should differ for different prompts")
	}
}

func TestComputeHash_DiffersOnContext(t *testing.T) {
	h1 := computeHash("blog-draft", "prompt", []string{"context A"})
	h2 := computeHash("blog-draft", "prompt", []string{"context B"})
	if h1 == h2 {
		t.Error("hashes should differ for different context")
	}
}

func TestComputeHash_NilContext(t *testing.T) {
	h := computeHash("retry", "abc123", nil)
	if h == "" {
		t.Error("hash should not be empty")
	}
	if len(h) != 64 {
		t.Errorf("hash length = %d, want 64", len(h))
	}
}

func TestComputeHash_Length(t *testing.T) {
	h := computeHash("cmd", "prompt", []string{"line"})
	if len(h) != 64 {
		t.Errorf("SHA-256 hex hash length = %d, want 64", len(h))
	}
}
