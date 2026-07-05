package config

import "strings"

// MatchModelPattern reports whether value matches pattern, where pattern
// may use '*' as a wildcard (prefix, suffix, substring, or exact match).
//
// Examples:
//
//	"gemini-*"   matches "gemini-2.5-flash"           (prefix)
//	"*-preview"  matches "gemini-3-pro-preview"       (suffix)
//	"*flash*"    matches "gemini-2.5-flash-lite"       (substring)
//	"gpt-5-codex" matches "gpt-5-codex" only          (exact)
//	"*"          matches everything                    (any)
func MatchModelPattern(pattern, value string) bool {
	if pattern == "*" {
		return true
	}
	// No wildcard: exact match.
	if !strings.Contains(pattern, "*") {
		return pattern == value
	}

	parts := strings.Split(pattern, "*")
	// parts is never empty because pattern is non-empty.

	switch len(parts) {
	case 2:
		left, right := parts[0], parts[1]
		// Prefix match: "abc*"
		if right == "" {
			return strings.HasPrefix(value, left)
		}
		// Suffix match: "*abc"
		if left == "" {
			return strings.HasSuffix(value, right)
		}
		// Substring match: "abc*def" (one * in the middle)
		return strings.HasPrefix(value, left) &&
			strings.HasSuffix(value, right) &&
			len(value) >= len(left)+len(right)
	case 3:
		// "abc*def*ghi" — two wildcards.
		// Only the case "*abc*" (left=="", right!="", middle!="") makes sense for
		// our use, but we handle all forms generically.
		// Generic: find left prefix, then find middle substring, then verify right suffix.
		left, middle, right := parts[0], parts[1], parts[2]
		if left != "" {
			if !strings.HasPrefix(value, left) {
				return false
			}
			value = value[len(left):]
		}
		if middle != "" {
			idx := strings.Index(value, middle)
			if idx < 0 {
				return false
			}
			value = value[idx+len(middle):]
		}
		if right != "" {
			return strings.HasSuffix(value, right)
		}
		return true
	default:
		// 4+ parts: general case with multiple wildcards.
		// Match segments in order.
		rest := value
		for i, part := range parts {
			if part == "" {
				continue
			}
			if i == 0 {
				// First non-empty part must match as prefix.
				if !strings.HasPrefix(rest, part) {
					return false
				}
				rest = rest[len(part):]
				continue
			}
			if i == len(parts)-1 {
				// Last non-empty part must match as suffix.
				return strings.HasSuffix(rest, part)
			}
			// Middle part must be found somewhere in the remaining string.
			idx := strings.Index(rest, part)
			if idx < 0 {
				return false
			}
			rest = rest[idx+len(part):]
		}
		return true
	}
}
