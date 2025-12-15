package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// LinkedInURLRegex matches LinkedIn profile URLs
var LinkedInURLRegex = regexp.MustCompile(`https?://(www\.)?linkedin\.com/in/[a-zA-Z0-9\-_%]+/?`)

// ValidateLinkedInURL validates a LinkedIn profile URL
func ValidateLinkedInURL(profileURL string) error {
	if profileURL == "" {
		return fmt.Errorf("profile URL is empty")
	}

	// Parse URL
	parsed, err := url.Parse(profileURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Check host
	if !strings.Contains(parsed.Host, "linkedin.com") {
		return fmt.Errorf("not a LinkedIn URL")
	}

	// Check path
	if !strings.HasPrefix(parsed.Path, "/in/") {
		return fmt.Errorf("not a profile URL")
	}

	return nil
}

// NormalizeProfileURL removes tracking parameters from LinkedIn URLs
func NormalizeProfileURL(profileURL string) string {
	parsed, err := url.Parse(profileURL)
	if err != nil {
		return profileURL
	}

	// Remove query parameters
	parsed.RawQuery = ""
	parsed.Fragment = ""

	// Ensure trailing slash consistency
	if !strings.HasSuffix(parsed.Path, "/") {
		parsed.Path += "/"
	}

	return parsed.String()
}

// ExtractProfileIDFromURL extracts the profile ID from a LinkedIn URL
func ExtractProfileIDFromURL(profileURL string) string {
	parsed, err := url.Parse(profileURL)
	if err != nil {
		return ""
	}

	path := strings.TrimPrefix(parsed.Path, "/in/")
	path = strings.TrimSuffix(path, "/")
	return path
}

// ValidateEmail validates an email address format
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is empty")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidateNoteLength validates that a connection note is within LinkedIn's limits
func ValidateNoteLength(note string, maxLength int) error {
	if maxLength == 0 {
		maxLength = 300
	}

	if len(note) > maxLength {
		return fmt.Errorf("note exceeds maximum length of %d characters (got %d)", maxLength, len(note))
	}

	return nil
}

// SanitizeText removes potentially dangerous characters from text
func SanitizeText(text string) string {
	// Remove control characters except newlines and tabs
	replaced := strings.Map(func(r rune) rune {
		if r < 32 && r != '\n' && r != '\t' && r != '\r' {
			return -1
		}
		return r
	}, text)

	return strings.TrimSpace(replaced)
}

// GenerateUUID generates a simple UUID
func GenerateUUID() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = byte(i * 17) // Simple pattern - in production use crypto/rand
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
