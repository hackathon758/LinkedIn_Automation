package search

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/go-rod/rod"

	"linkedin-automation/config"
	"linkedin-automation/database"
	"linkedin-automation/logger"
	"linkedin-automation/stealth"
)

// Searcher handles LinkedIn user search
type Searcher struct {
	config    config.SearchConfig
	db        *database.DB
	logger    *logger.Logger
	timing    *stealth.TimingController
	scrolling *stealth.ScrollController
	mouse     *stealth.MouseHoverController
}

// NewSearcher creates a new Searcher
func NewSearcher(
	cfg config.SearchConfig,
	db *database.DB,
	log *logger.Logger,
	stealthCfg config.StealthConfig,
) *Searcher {
	return &Searcher{
		config:    cfg,
		db:        db,
		logger:    log.WithComponent("search"),
		timing:    stealth.NewTimingController(stealthCfg.Timing),
		scrolling: stealth.NewScrollController(stealthCfg.Scrolling),
		mouse:     stealth.NewMouseHoverController(stealthCfg.Mouse),
	}
}

// ProfileInfo contains extracted profile information
type ProfileInfo struct {
	ProfileURL string
	FirstName  string
	LastName   string
	JobTitle   string
	Company    string
	Location   string
}

// SearchResult contains the results of a search operation
type SearchResult struct {
	Profiles     []ProfileInfo
	TotalFound   int
	PagesScraped int
	Duplicates   int
	Errors       []string
}

// BuildSearchURL constructs a LinkedIn search URL with filters
func (s *Searcher) BuildSearchURL() string {
	baseURL := "https://www.linkedin.com/search/results/people/"

	params := url.Values{}

	// Build keywords from job titles and keywords
	var keywords []string
	keywords = append(keywords, s.config.JobTitles...)
	keywords = append(keywords, s.config.Keywords...)
	if len(keywords) > 0 {
		params.Set("keywords", strings.Join(keywords, " "))
	}

	// Add location filter
	if len(s.config.Locations) > 0 {
		params.Set("geoUrn", s.getGeoUrn(s.config.Locations[0]))
	}

	// Add company filter (if available)
	// Note: Company filtering requires company LinkedIn IDs

	params.Set("origin", "GLOBAL_SEARCH_HEADER")

	return baseURL + "?" + params.Encode()
}

// getGeoUrn returns LinkedIn geo URN for common locations
func (s *Searcher) getGeoUrn(location string) string {
	// Common location URNs (in production, this would be a larger mapping)
	locationURNs := map[string]string{
		"san francisco":        "90009496",
		"san francisco bay area": "90009496",
		"new york":             "90009563",
		"new york city":        "90009563",
		"los angeles":          "90009494",
		"chicago":              "90009457",
		"seattle":              "90009483",
		"boston":               "90009611",
		"austin":               "90009493",
		"denver":               "90009481",
		"united states":        "103644278",
	}

	lower := strings.ToLower(location)
	if urn, ok := locationURNs[lower]; ok {
		return "[\"" + urn + "\"]"
	}
	return ""
}

// Search performs a search and extracts profile URLs
func (s *Searcher) Search(page *rod.Page) (*SearchResult, error) {
	result := &SearchResult{}

	searchURL := s.BuildSearchURL()
	s.logger.Info("starting search", "url", searchURL)

	// Navigate to search
	err := page.Navigate(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate: %w", err)
	}

	time.Sleep(s.timing.GetPageLoadDelay())
	err = page.WaitLoad()
	if err != nil {
		s.logger.LogError("page load", err, nil)
	}

	// Process pages
	for pageNum := 1; pageNum <= s.config.MaxPages; pageNum++ {
		s.logger.Info("processing page", "page", pageNum)

		// Extract profiles from current page
		profiles, err := s.extractProfiles(page)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("page %d: %v", pageNum, err))
			continue
		}

		// Filter duplicates
		for _, profile := range profiles {
			processed, _ := s.db.IsProfileProcessed(profile.ProfileURL)
			if processed {
				result.Duplicates++
				continue
			}
			result.Profiles = append(result.Profiles, profile)
		}

		result.PagesScraped++
		result.TotalFound += len(profiles)

		// Try to go to next page
		if pageNum < s.config.MaxPages {
			hasNext := s.goToNextPage(page)
			if !hasNext {
				s.logger.Info("no more pages available")
				break
			}
			// Wait between pages
			time.Sleep(s.timing.GetThinkTime())
		}
	}

	s.logger.Info("search complete",
		"total_found", result.TotalFound,
		"unique", len(result.Profiles),
		"duplicates", result.Duplicates,
		"pages", result.PagesScraped)

	return result, nil
}

// extractProfiles extracts profile information from the current page
func (s *Searcher) extractProfiles(page *rod.Page) ([]ProfileInfo, error) {
	var profiles []ProfileInfo

	// Scroll to load all results
	s.scrollToLoadResults(page)

	// Find profile links
	profileLinks, err := page.Elements(`a[href*="/in/"]`)
	if err != nil {
		return nil, err
	}

	profileURLRegex := regexp.MustCompile(`/in/([a-zA-Z0-9\-_%]+)`)
	seenURLs := make(map[string]bool)

	for _, link := range profileLinks {
		href, err := link.Attribute("href")
		if err != nil || href == nil {
			continue
		}

		// Extract and normalize profile URL
		matches := profileURLRegex.FindStringSubmatch(*href)
		if len(matches) < 2 {
			continue
		}

		profileURL := fmt.Sprintf("https://www.linkedin.com/in/%s/", matches[1])

		// Skip if already seen in this batch
		if seenURLs[profileURL] {
			continue
		}
		seenURLs[profileURL] = true

		// Try to extract additional info from the search result card
		profile := ProfileInfo{ProfileURL: profileURL}

		// Try to get name and other details from parent elements
		parent := link
		for i := 0; i < 5; i++ {
			parent, err = parent.Parent()
			if err != nil {
				break
			}

			// Look for name span
			nameEl, err := parent.Element(".entity-result__title-text")
			if err == nil && nameEl != nil {
				name, _ := nameEl.Text()
				parts := strings.Fields(name)
				if len(parts) >= 1 {
					profile.FirstName = parts[0]
				}
				if len(parts) >= 2 {
					profile.LastName = strings.Join(parts[1:], " ")
				}
			}

			// Look for headline (job title)
			headlineEl, err := parent.Element(".entity-result__primary-subtitle")
			if err == nil && headlineEl != nil {
				headline, _ := headlineEl.Text()
				profile.JobTitle = strings.TrimSpace(headline)
			}

			// Look for location
			locationEl, err := parent.Element(".entity-result__secondary-subtitle")
			if err == nil && locationEl != nil {
				location, _ := locationEl.Text()
				profile.Location = strings.TrimSpace(location)
			}
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// scrollToLoadResults scrolls through the page to load all dynamic results
func (s *Searcher) scrollToLoadResults(page *rod.Page) {
	// Get page height
	height, err := page.Eval(`() => document.body.scrollHeight`)
	if err != nil {
		return
	}

	totalHeight := int(height.Value.Num())
	currentY := 0
	scrollStep := 400

	for currentY < totalHeight {
		// Generate natural scroll
		steps := s.scrolling.GenerateScrollSequence(currentY+scrollStep, currentY)

		for _, step := range steps {
			page.Eval(fmt.Sprintf(`() => window.scrollBy(0, %d)`, step.DeltaY))
			time.Sleep(step.Duration)
		}

		currentY += scrollStep

		// Random pause while scrolling
		if s.scrolling.ShouldPauseWhileScrolling() {
			time.Sleep(s.scrolling.GetRandomScrollPause())
		}

		// Update total height (might have changed with lazy loading)
		height, err = page.Eval(`() => document.body.scrollHeight`)
		if err == nil {
			totalHeight = int(height.Value.Num())
		}
	}
}

// goToNextPage navigates to the next page of results
func (s *Searcher) goToNextPage(page *rod.Page) bool {
	// Look for next button
	nextButton, err := page.Element(`button[aria-label="Next"]`)
	if err != nil {
		// Try alternate selector
		nextButton, err = page.Element(`.artdeco-pagination__button--next`)
		if err != nil {
			return false
		}
	}

	// Check if button is disabled
	disabled, _ := nextButton.Attribute("disabled")
	if disabled != nil {
		return false
	}

	// Click with natural movement
	nextButton.Click("left", 1)

	// Wait for page load
	time.Sleep(s.timing.GetPageLoadDelay())
	page.WaitLoad()

	return true
}
