package config

import "testing"

func TestMatchModelPattern(t *testing.T) {
	tests := []struct {
		pattern string
		value   string
		want    bool
	}{
		// Exact match
		{"gpt-5-codex", "gpt-5-codex", true},
		{"gpt-5-codex", "gpt-5-codex-other", false},
		{"gpt-5-codex", "gpt-5", false},

		// Prefix match
		{"gemini-*", "gemini-2.5-flash", true},
		{"gemini-*", "gemini-3-pro-preview", true},
		{"gemini-*", "claude-opus-4", false},
		{"gemini-*", "gemini", false},

		// Suffix match
		{"*-preview", "gemini-3-pro-preview", true},
		{"*-preview", "claude-opus-4-preview", true},
		{"*-preview", "gemini-3-pro", false},

		// Substring match
		{"*flash*", "gemini-2.5-flash", true},
		{"*flash*", "gemini-2.5-flash-lite", true},
		{"*flash*", "gemini-2.5-pro", false},

		// Multi-segment with wildcards
		{"gemini-*-flash", "gemini-2.5-flash", true},
		{"gemini-*-flash", "gemini-3-pro-flash", true},
		{"gemini-*-flash", "claude-2.5-flash", false},

		// Match all
		{"*", "anything", true},
		{"*", "", true},

		// Edge cases
		{"", "", true}, // empty pattern matches empty value exactly
		{"", "x", false},

		// Case sensitivity (callers must lower-case before calling)
		{"Gemini-*", "gemini-2.5-flash", false},
		{"gemini-*", "GEMINI-2.5-FLASH", false},

		// Multiple wildcards
		{"*gemini*flash*", "foo-gemini-bar-flash-baz", true},
		{"*gemini*flash*", "foo-gemini-bar-pro-baz", false},
	}

	for _, tc := range tests {
		got := MatchModelPattern(tc.pattern, tc.value)
		if got != tc.want {
			t.Errorf("MatchModelPattern(%q, %q) = %v, want %v",
				tc.pattern, tc.value, got, tc.want)
		}
	}
}
