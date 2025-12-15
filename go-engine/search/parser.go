package search

import (
	"regexp"
	"strings"
)

// ProfileURLPattern matches LinkedIn profile URLs
var ProfileURLPattern = regexp.MustCompile(`https?://(www\.)?linkedin\.com/in/([a-zA-Z0-9\-_%]+)/?`)

// ExtractProfileURLs extracts all valid LinkedIn profile URLs from HTML content
func ExtractProfileURLs(html string) []string {
	matches := ProfileURLPattern.FindAllStringSubmatch(html, -1)

	urlSet := make(map[string]bool)
	var urls []string

	for _, match := range matches {
		if len(match) >= 3 {
			// Normalize URL
			normalized := "https://www.linkedin.com/in/" + match[2] + "/"

			if !urlSet[normalized] {
				urlSet[normalized] = true
				urls = append(urls, normalized)
			}
		}
	}

	return urls
}

// ExtractProfileID extracts the profile ID from a LinkedIn URL
func ExtractProfileID(profileURL string) string {
	match := ProfileURLPattern.FindStringSubmatch(profileURL)
	if len(match) >= 3 {
		return match[2]
	}
	return ""
}

// ParseFullName splits a full name into first and last name
func ParseFullName(fullName string) (firstName, lastName string) {
	fullName = strings.TrimSpace(fullName)
	parts := strings.Fields(fullName)

	if len(parts) == 0 {
		return "", ""
	}

	if len(parts) == 1 {
		return parts[0], ""
	}

	return parts[0], strings.Join(parts[1:], " ")
}

// CleanCompanyName cleans up company name from search results
func CleanCompanyName(company string) string {
	// Remove common suffixes
	suffixes := []string{" Inc.", " Inc", " LLC", " Ltd.", " Ltd", " Corp.", " Corp", " Co.", " Co"}

	result := strings.TrimSpace(company)
	for _, suffix := range suffixes {
		if strings.HasSuffix(result, suffix) {
			result = strings.TrimSuffix(result, suffix)
		}
	}

	return result
}

// IsValidProfileURL checks if a URL is a valid LinkedIn profile URL
func IsValidProfileURL(url string) bool {
	return ProfileURLPattern.MatchString(url)
}
