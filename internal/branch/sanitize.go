package branch

import (
	"regexp"
	"strings"
	"unicode"
)

// Sanitize creates a valid git branch name from a Linear identifier and title.
// Rules:
// - Lowercase the Linear identifier (DEV-123 → dev-123)
// - Slugify title: lowercase, replace spaces with hyphens
// - Remove all chars except [a-z0-9-]
// - Collapse multiple hyphens to single
// - Strip leading/trailing hyphens
// - Max total length: 32 chars (truncate title part if needed)
// - If title becomes empty after sanitization, return just identifier
func Sanitize(identifier, title string) string {
	// Lowercase identifier
	identifier = strings.ToLower(identifier)

	// Slugify title: lowercase and replace spaces with hyphens
	title = strings.ToLower(title)
	title = strings.ReplaceAll(title, " ", "-")

	// Remove non-ASCII characters (emoji, unicode)
	title = removeNonASCII(title)

	// Remove all chars except [a-z0-9-]
	title = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(title, "")

	// Collapse multiple hyphens to single
	title = regexp.MustCompile(`-+`).ReplaceAllString(title, "-")

	// Strip leading/trailing hyphens
	title = strings.Trim(title, "-")

	// If title is empty, return just identifier
	if title == "" {
		return identifier
	}

	// Combine identifier and title
	result := identifier + "-" + title

	// Truncate to max 32 chars if needed
	if len(result) > 32 {
		// Calculate how much space we have for the title
		maxTitleLen := 32 - len(identifier) - 1 // -1 for the hyphen between identifier and title
		if maxTitleLen < 1 {
			// If identifier itself is too long, just return it truncated
			return identifier[:32]
		}
		// Truncate title and remove trailing hyphen if present
		title = title[:maxTitleLen]
		title = strings.TrimRight(title, "-")
		result = identifier + "-" + title
	}

	return result
}

// removeNonASCII removes all non-ASCII characters from a string
func removeNonASCII(s string) string {
	return strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII {
			return -1 // Remove the character
		}
		return r
	}, s)
}
