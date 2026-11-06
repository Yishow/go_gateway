package workspace

import "strings"

// homeTrimLabel matches ECMAScript String.trim's WhiteSpace and LineTerminator
// set, used by the canonical frontend candidate resolver. Go TrimSpace differs
// for BOM (U+FEFF) and NEXT LINE (U+0085).
func homeTrimLabel(value string) string {
	return strings.TrimFunc(value, func(r rune) bool {
		switch r {
		case '\t', '\n', '\v', '\f', '\r', ' ', '\u00a0', '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000', '\ufeff':
			return true
		default:
			return r >= '\u2000' && r <= '\u200a'
		}
	})
}
