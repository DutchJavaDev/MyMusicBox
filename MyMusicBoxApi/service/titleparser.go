package service

import (
	"fmt"
	"regexp"
	"strings"
)

// FixJSONStringValues scans all JSON string values and escapes invalid characters.
func FixJSONStringValues(input string) string {
	// Regex pattern to match string values like: "key": "value"
	re := regexp.MustCompile(`"(?:[^"\\]|\\.)*?"\s*:\s*"(.*?)"`)

	return strings.TrimSpace(re.ReplaceAllStringFunc(input, func(match string) string {
		parts := strings.SplitN(match, ":", 2)
		if len(parts) != 2 {
			return match
		}
		key := strings.TrimSpace(parts[0])
		rawValue := strings.TrimSpace(parts[1])

		// Remove surrounding quotes
		unquoted := strings.Trim(rawValue, `"`)

		// Escape invalid characters inside the value
		escaped := escapeUnsafeJSONCharacters(unquoted)

		return fmt.Sprintf(`%s: "%s"`, key, escaped)
	}))
}

// escapeUnsafeJSONCharacters escapes characters that would break JSON string parsing.
func escapeUnsafeJSONCharacters(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			if r < 0x20 {
				// Escape control characters as unicode
				fmt.Fprintf(&b, `\u%04x`, r)
			} else {
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}
