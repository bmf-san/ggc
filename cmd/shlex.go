package cmd

import "strings"

// tokenize splits a command string into tokens using POSIX-style word splitting.
// It handles single-quoted strings, double-quoted strings, and backslash escapes
// inside double quotes. This is used for alias command expansion and avoids the
// incorrect splitting that results from the naive strings.Split approach.
//
// Examples:
//
//	tokenize(`commit -m "fix bug"`)         → ["commit", "-m", "fix bug"]
//	tokenize(`commit -m 'fix the bug'`)     → ["commit", "-m", "fix the bug"]
//	tokenize(`echo "it's a \"test\""`)      → ["echo", `it's a "test"`]
//
//nolint:revive
func tokenize(s string) []string {
	var tokens []string
	var cur strings.Builder
	inSingle := false
	inDouble := false
	inToken := false // true once we have started accumulating a token

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\'' && !inDouble:
			// Toggle single-quote mode; the quotes themselves are stripped.
			inSingle = !inSingle
			inToken = true
		case c == '"' && !inSingle:
			// Toggle double-quote mode; the quotes themselves are stripped.
			inDouble = !inDouble
			inToken = true
		case c == '\\' && inDouble && i+1 < len(s):
			// Inside double quotes, a backslash escapes the next character.
			i++
			cur.WriteByte(s[i])
			inToken = true
		case (c == ' ' || c == '\t') && !inSingle && !inDouble:
			// Unquoted whitespace flushes the current token (which may be "").
			if inToken {
				tokens = append(tokens, cur.String())
				cur.Reset()
				inToken = false
			}
		default:
			cur.WriteByte(c)
			inToken = true
		}
	}

	if inToken {
		tokens = append(tokens, cur.String())
	}
	return tokens
}
