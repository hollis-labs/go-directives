package directives

import (
	"testing"
)

func TestParseConfigPairs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]string
	}{
		{"single pair", "voice=technical", map[string]string{"voice": "technical"}},
		{"multiple pairs", "voice=technical, project=carrier", map[string]string{"voice": "technical", "project": "carrier"}},
		{"with spaces", " voice = technical , project = carrier ", map[string]string{"voice": "technical", "project": "carrier"}},
		{"empty string", "", map[string]string{}},
		{"no equals", "invalid", map[string]string{}},
		{"mixed valid/invalid", "voice=casual, bad, tags=a", map[string]string{"voice": "casual", "tags": "a"}},
		{"empty key", "=value", map[string]string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseConfigPairs(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d pairs, want %d: %v", len(got), len(tt.want), got)
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("got[%q] = %q, want %q", k, got[k], v)
				}
			}
		})
	}
}

func TestConfigStack_Cascade(t *testing.T) {
	var cs configStack

	// Root scope
	cs.push()
	cs.set("voice", "technical")
	cs.set("project", "carrier")

	// Nested scope overrides voice
	cs.push()
	cs.set("voice", "casual")

	merged := cs.merged()
	if merged["voice"] != "casual" {
		t.Errorf("voice = %q, want casual", merged["voice"])
	}
	if merged["project"] != "carrier" {
		t.Errorf("project = %q, want carrier", merged["project"])
	}

	// Pop nested scope — voice reverts
	cs.pop()
	merged = cs.merged()
	if merged["voice"] != "technical" {
		t.Errorf("voice after pop = %q, want technical", merged["voice"])
	}
}

func TestConfigStack_EmptyMerged(t *testing.T) {
	var cs configStack
	merged := cs.merged()
	if len(merged) != 0 {
		t.Errorf("expected empty map, got %v", merged)
	}
}

func TestConfigStack_SetWithoutPush(t *testing.T) {
	var cs configStack
	cs.set("key", "value")
	merged := cs.merged()
	if merged["key"] != "value" {
		t.Errorf("key = %q, want value", merged["key"])
	}
}
