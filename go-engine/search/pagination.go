package pagination

import (
	"time"

	"github.com/go-rod/rod"
)

// Paginator handles pagination through search results
type Paginator struct {
	currentPage int
	maxPages    int
	hasMore     bool
}

// NewPaginator creates a new paginator
func NewPaginator(maxPages int) *Paginator {
	return &Paginator{
		currentPage: 1,
		maxPages:    maxPages,
		hasMore:     true,
	}
}

// CurrentPage returns the current page number
func (p *Paginator) CurrentPage() int {
	return p.currentPage
}

// HasMore returns true if there are more pages to process
func (p *Paginator) HasMore() bool {
	return p.hasMore && p.currentPage <= p.maxPages
}

// NextPage attempts to navigate to the next page
func (p *Paginator) NextPage(page *rod.Page) bool {
	// Look for pagination controls
	nextSelectors := []string{
		`button[aria-label="Next"]`,
		`.artdeco-pagination__button--next`,
		`button.artdeco-pagination__button--next`,
	}

	var nextButton *rod.Element
	var err error

	for _, selector := range nextSelectors {
		nextButton, err = page.Timeout(5 * time.Second).Element(selector)
		if err == nil && nextButton != nil {
			break
		}
	}

	if nextButton == nil {
		p.hasMore = false
		return false
	}

	// Check if disabled
	disabled, _ := nextButton.Attribute("disabled")
	if disabled != nil {
		p.hasMore = false
		return false
	}

	// Click next
	err = nextButton.Click("left", 1)
	if err != nil {
		p.hasMore = false
		return false
	}

	p.currentPage++
	return true
}

// Reset resets the paginator
func (p *Paginator) Reset() {
	p.currentPage = 1
	p.hasMore = true
}

// SetMaxPages updates the maximum pages to process
func (p *Paginator) SetMaxPages(max int) {
	p.maxPages = max
}

// Progress returns pagination progress as a percentage
func (p *Paginator) Progress() float64 {
	if p.maxPages == 0 {
		return 100.0
	}
	return float64(p.currentPage) / float64(p.maxPages) * 100.0
}
