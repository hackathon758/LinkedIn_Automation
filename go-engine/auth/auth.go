package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"

	"linkedin-automation/config"
	"linkedin-automation/database"
	"linkedin-automation/logger"
	"linkedin-automation/stealth"
)

// Authenticator handles LinkedIn authentication
type Authenticator struct {
	config       config.CredentialsConfig
	db           *database.DB
	logger       *logger.Logger
	timing       *stealth.TimingController
	typing       *stealth.TypingSimulator
	bezier       *stealth.BezierMouse
	fingerprint  *stealth.FingerprintMasker
}

// NewAuthenticator creates a new Authenticator
func NewAuthenticator(
	cfg config.CredentialsConfig,
	db *database.DB,
	log *logger.Logger,
	stealthCfg config.StealthConfig,
) *Authenticator {
	return &Authenticator{
		config:      cfg,
		db:          db,
		logger:      log.WithComponent("auth"),
		timing:      stealth.NewTimingController(stealthCfg.Timing),
		typing:      stealth.NewTypingSimulator(stealthCfg.Timing),
		bezier:      stealth.NewBezierMouse(stealthCfg.Bezier),
		fingerprint: stealth.NewFingerprintMasker(stealthCfg.Fingerprint),
	}
}

// LoginResult represents the result of a login attempt
type LoginResult struct {
	Success          bool
	SecurityChallenge bool
	ChallengeType    string // "2fa", "captcha", "verification"
	ErrorMessage     string
}

// Login performs LinkedIn login with realistic behavior
func (a *Authenticator) Login(browser *rod.Browser) (*rod.Page, *LoginResult, error) {
	a.logger.Info("starting login process")

	// Validate credentials
	if a.config.Email == "" || a.config.Password == "" {
		return nil, &LoginResult{Success: false, ErrorMessage: "credentials not configured"}, fmt.Errorf("credentials not configured")
	}

	// Create new page
	page, err := browser.Page(proto.TargetCreateTarget{URL: "about:blank"})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create page: %w", err)
	}

	// Apply stealth scripts
	err = a.applyStealthScripts(page)
	if err != nil {
		a.logger.LogError("apply stealth scripts", err, nil)
	}

	// Navigate to LinkedIn login
	a.logger.Info("navigating to login page")
	err = page.Navigate("https://www.linkedin.com/login")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to navigate: %w", err)
	}

	// Wait for page load
	time.Sleep(a.timing.GetPageLoadDelay())
	err = page.WaitLoad()
	if err != nil {
		return nil, nil, fmt.Errorf("page load timeout: %w", err)
	}

	// Check for security checkpoint before login
	if result := a.checkSecurityChallenge(page); result != nil {
		return page, result, nil
	}

	// Enter email
	a.logger.Info("entering email")
	emailInput, err := page.Element("#username")
	if err != nil {
		return nil, nil, fmt.Errorf("email field not found: %w", err)
	}

	err = a.typeWithRealism(emailInput, a.config.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to enter email: %w", err)
	}

	// Small delay between fields
	time.Sleep(a.timing.GetActionDelay())

	// Enter password
	a.logger.Info("entering password")
	passwordInput, err := page.Element("#password")
	if err != nil {
		return nil, nil, fmt.Errorf("password field not found: %w", err)
	}

	err = a.typeWithRealism(passwordInput, a.config.Password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to enter password: %w", err)
	}

	// Think time before clicking
	time.Sleep(a.timing.GetThinkTime())

	// Click login button
	a.logger.Info("clicking login button")
	loginButton, err := page.Element("button[type='submit']")
	if err != nil {
		return nil, nil, fmt.Errorf("login button not found: %w", err)
	}

	err = a.clickWithRealism(page, loginButton)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to click login: %w", err)
	}

	// Wait for navigation
	time.Sleep(3 * time.Second)
	err = page.WaitLoad()
	if err != nil {
		a.logger.LogError("wait after login", err, nil)
	}

	// Check for security challenges
	if result := a.checkSecurityChallenge(page); result != nil {
		return page, result, nil
	}

	// Verify successful login
	if a.isLoggedIn(page) {
		a.logger.Info("login successful")

		// Save cookies for future sessions
		err = a.saveCookies(page)
		if err != nil {
			a.logger.LogError("save cookies", err, nil)
		}

		return page, &LoginResult{Success: true}, nil
	}

	// Check for login errors
	errorMsg := a.getLoginError(page)
	return page, &LoginResult{Success: false, ErrorMessage: errorMsg}, nil
}

// TrySessionRestore attempts to restore a session from saved cookies
func (a *Authenticator) TrySessionRestore(browser *rod.Browser) (*rod.Page, bool, error) {
	a.logger.Info("attempting session restore")

	cookies, err := a.db.GetCookies()
	if err != nil || len(cookies) == 0 {
		a.logger.Info("no valid cookies found")
		return nil, false, nil
	}

	// Create page and inject cookies
	page, err := browser.Page(proto.TargetCreateTarget{URL: "about:blank"})
	if err != nil {
		return nil, false, err
	}

	// Apply stealth
	a.applyStealthScripts(page)

	// Navigate to LinkedIn first (cookies require same domain)
	err = page.Navigate("https://www.linkedin.com")
	if err != nil {
		return nil, false, err
	}
	time.Sleep(2 * time.Second)

	// Inject cookies
	for _, cookie := range cookies {
		err = page.SetCookies([]*proto.NetworkCookieParam{{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Expires:  proto.TimeSinceEpoch(cookie.ExpiresAt.Unix()),
		}})
		if err != nil {
			a.logger.LogError("inject cookie", err, map[string]interface{}{"cookie": cookie.Name})
		}
	}

	// Reload page with cookies
	err = page.Reload()
	if err != nil {
		return nil, false, err
	}
	time.Sleep(3 * time.Second)

	// Check if logged in
	if a.isLoggedIn(page) {
		a.logger.Info("session restored successfully")
		return page, true, nil
	}

	a.logger.Info("session restore failed, cookies expired")
	a.db.ClearCookies()
	return nil, false, nil
}

// applyStealthScripts injects anti-detection JavaScript
func (a *Authenticator) applyStealthScripts(page *rod.Page) error {
	script := a.fingerprint.GetAllMaskingScripts()
	if script != "" {
		_, err := page.Eval(script)
		return err
	}
	return nil
}

// typeWithRealism types text with human-like patterns
func (a *Authenticator) typeWithRealism(element *rod.Element, text string) error {
	sequence := a.typing.GenerateTypingSequence(text)

	for _, char := range sequence {
		if char.IsBurstPause {
			time.Sleep(char.Delay)
			continue
		}

		if char.IsBackspace {
			element.MustInput("\b")
			time.Sleep(char.Delay)
			continue
		}

		element.MustInput(string(char.Char))
		time.Sleep(char.Delay)
	}

	return nil
}

// clickWithRealism clicks an element with natural mouse movement
func (a *Authenticator) clickWithRealism(page *rod.Page, element *rod.Element) error {
	// Get element position
	box, err := element.Shape()
	if err != nil {
		return element.Click(proto.InputMouseButtonLeft, 1)
	}

	quads := box.Quads
	if len(quads) == 0 || len(quads[0]) < 4 {
		return element.Click(proto.InputMouseButtonLeft, 1)
	}

	// Calculate center of element
	centerX := (quads[0][0] + quads[0][2]) / 2
	centerY := (quads[0][1] + quads[0][5]) / 2

	// Get current mouse position (assume 0,0 if unknown)
	startX, startY := 0.0, 0.0

	// Generate Bezier path
	path := a.bezier.GeneratePath(startX, startY, centerX, centerY)
	durations := a.bezier.GetMovementDurations(len(path), 500*time.Millisecond)

	// Move mouse along path
	for i, point := range path {
		if i > 0 && i-1 < len(durations) {
			time.Sleep(durations[i-1])
		}
		page.Mouse.MustMoveTo(point.X, point.Y)
	}

	// Small hover delay before click
	time.Sleep(50 * time.Millisecond)

	// Click
	return page.Mouse.Click(proto.InputMouseButtonLeft, 1)
}

// checkSecurityChallenge detects security checkpoints
func (a *Authenticator) checkSecurityChallenge(page *rod.Page) *LoginResult {
	currentURL := page.MustInfo().URL

	// Check for 2FA
	if strings.Contains(currentURL, "checkpoint") || strings.Contains(currentURL, "challenge") {
		// Try to determine challenge type
		html, _ := page.HTML()

		if strings.Contains(html, "verification code") || strings.Contains(html, "two-step") {
			a.logger.Info("2FA challenge detected")
			return &LoginResult{
				Success:           false,
				SecurityChallenge: true,
				ChallengeType:     "2fa",
				ErrorMessage:      "Two-factor authentication required. Please complete manually.",
			}
		}

		if strings.Contains(html, "captcha") || strings.Contains(html, "CAPTCHA") {
			a.logger.Info("CAPTCHA challenge detected")
			return &LoginResult{
				Success:           false,
				SecurityChallenge: true,
				ChallengeType:     "captcha",
				ErrorMessage:      "CAPTCHA verification required. Please complete manually.",
			}
		}

		return &LoginResult{
			Success:           false,
			SecurityChallenge: true,
			ChallengeType:     "verification",
			ErrorMessage:      "Security verification required. Please complete manually.",
		}
	}

	return nil
}

// isLoggedIn checks if user is successfully logged in
func (a *Authenticator) isLoggedIn(page *rod.Page) bool {
	currentURL := page.MustInfo().URL

	// Check URL patterns
	if strings.Contains(currentURL, "/feed") ||
		strings.Contains(currentURL, "/mynetwork") ||
		strings.Contains(currentURL, "/messaging") {
		return true
	}

	// Check for feed elements
	feedElement, err := page.Timeout(5 * time.Second).Element(".feed-shared-update-v2")
	if err == nil && feedElement != nil {
		return true
	}

	// Check for navigation elements
	navElement, err := page.Timeout(2 * time.Second).Element(".global-nav")
	if err == nil && navElement != nil {
		return true
	}

	return false
}

// getLoginError extracts error message from failed login
func (a *Authenticator) getLoginError(page *rod.Page) string {
	// Look for error elements
	errorElement, err := page.Timeout(2 * time.Second).Element(".form__label--error, .alert-content")
	if err == nil && errorElement != nil {
		text, _ := errorElement.Text()
		return text
	}
	return "Login failed - unknown error"
}

// saveCookies saves session cookies to database
func (a *Authenticator) saveCookies(page *rod.Page) error {
	cookies, err := page.Cookies([]string{"https://www.linkedin.com"})
	if err != nil {
		return err
	}

	var dbCookies []database.SessionCookie
	for _, c := range cookies {
		dbCookies = append(dbCookies, database.SessionCookie{
			ID:        fmt.Sprintf("cookie_%s_%d", c.Name, time.Now().UnixNano()),
			Name:      c.Name,
			Value:     c.Value,
			Domain:    c.Domain,
			Path:      c.Path,
			ExpiresAt: time.Unix(int64(c.Expires), 0),
			CreatedAt: time.Now(),
		})
	}

	return a.db.SaveCookies(dbCookies)
}
