package messaging

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

// ConnectionManager handles sending connection requests
type ConnectionManager struct {
	config     config.ConnectionConfig
	db         *database.DB
	logger     *logger.Logger
	timing     *stealth.TimingController
	typing     *stealth.TypingSimulator
	bezier     *stealth.BezierMouse
	mouse      *stealth.MouseHoverController
	templates  []string
}

// NewConnectionManager creates a new ConnectionManager
func NewConnectionManager(
	cfg config.ConnectionConfig,
	db *database.DB,
	log *logger.Logger,
	stealthCfg config.StealthConfig,
) *ConnectionManager {
	return &ConnectionManager{
		config:    cfg,
		db:        db,
		logger:    log.WithComponent("connection"),
		timing:    stealth.NewTimingController(stealthCfg.Timing),
		typing:    stealth.NewTypingSimulator(stealthCfg.Timing),
		bezier:    stealth.NewBezierMouse(stealthCfg.Bezier),
		mouse:     stealth.NewMouseHoverController(stealthCfg.Mouse),
		templates: cfg.Templates,
	}
}

// ConnectionRequest represents a connection request
type ConnectionRequest struct {
	ProfileURL  string
	FirstName   string
	LastName    string
	JobTitle    string
	Company     string
	Note        string
	TemplateIdx int
}

// ConnectionResult represents the result of a connection request
type ConnectionResult struct {
	Success      bool
	ProfileURL   string
	ErrorMessage string
	NeedsCaptcha bool
}

// SendConnectionRequest sends a connection request to a profile
func (cm *ConnectionManager) SendConnectionRequest(page *rod.Page, req *ConnectionRequest) (*ConnectionResult, error) {
	cm.logger.Info("sending connection request", "profile", req.ProfileURL)

	// Navigate to profile
	err := page.Navigate(req.ProfileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate: %w", err)
	}

	time.Sleep(cm.timing.GetPageLoadDelay())
	err = page.WaitLoad()
	if err != nil {
		cm.logger.LogError("page load", err, nil)
	}

	// Wait for profile to load
	time.Sleep(cm.timing.GetThinkTime())

	// Extract profile data if not provided
	if req.FirstName == "" {
		req.FirstName, req.LastName, req.JobTitle, req.Company = cm.extractProfileData(page)
	}

	// Find Connect button
	connectButton, err := cm.findConnectButton(page)
	if err != nil {
		return &ConnectionResult{
			Success:      false,
			ProfileURL:   req.ProfileURL,
			ErrorMessage: "Connect button not found - may already be connected",
		}, nil
	}

	// Click Connect button with realistic behavior
	err = cm.clickWithRealism(page, connectButton)
	if err != nil {
		return nil, fmt.Errorf("failed to click connect: %w", err)
	}

	// Wait for modal
	time.Sleep(time.Second)

	// Check if we need to add a note
	if len(cm.templates) > 0 && req.Note == "" {
		req.Note = cm.generateNote(req)
	}

	// Send with or without note
	if req.Note != "" {
		err = cm.sendWithNote(page, req.Note)
	} else {
		err = cm.sendWithoutNote(page)
	}

	if err != nil {
		return &ConnectionResult{
			Success:      false,
			ProfileURL:   req.ProfileURL,
			ErrorMessage: err.Error(),
		}, nil
	}

	// Record in database
	conn := &database.Connection{
		ID:         fmt.Sprintf("conn_%d", time.Now().UnixNano()),
		ProfileURL: req.ProfileURL,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		JobTitle:   req.JobTitle,
		Company:    req.Company,
		NoteSent:   req.Note,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}
	cm.db.SaveConnection(conn)
	cm.db.IncrementConnectionCount()
	cm.db.MarkProfileProcessed(req.ProfileURL)

	cm.logger.Info("connection request sent", "profile", req.ProfileURL)

	return &ConnectionResult{
		Success:    true,
		ProfileURL: req.ProfileURL,
	}, nil
}

// findConnectButton finds the Connect button on a profile page
func (cm *ConnectionManager) findConnectButton(page *rod.Page) (*rod.Element, error) {
	// Try various selectors
	selectors := []string{
		`button[aria-label*="Invite"]`,
		`button[aria-label*="Connect"]`,
		`button.pvs-profile-actions__action[aria-label*="connect"]`,
		`.pv-top-card-v2-ctas button:has-text("Connect")`,
		`button:has-text("Connect")`,
	}

	for _, selector := range selectors {
		btn, err := page.Timeout(3 * time.Second).Element(selector)
		if err == nil && btn != nil {
			// Verify it's visible
			visible, _ := btn.Visible()
			if visible {
				return btn, nil
			}
		}
	}

	// Check for "More" dropdown which might contain Connect
	moreBtn, err := page.Timeout(2 * time.Second).Element(`button[aria-label="More actions"]`)
	if err == nil && moreBtn != nil {
		moreBtn.Click(proto.InputMouseButtonLeft, 1)
		time.Sleep(500 * time.Millisecond)

		connectInMenu, err := page.Timeout(2 * time.Second).Element(`div[data-control-name="connect"]`)
		if err == nil && connectInMenu != nil {
			return connectInMenu, nil
		}
	}

	return nil, fmt.Errorf("connect button not found")
}

// extractProfileData extracts profile information from the current page
func (cm *ConnectionManager) extractProfileData(page *rod.Page) (firstName, lastName, jobTitle, company string) {
	// Try to get name
	nameEl, err := page.Timeout(3 * time.Second).Element(`h1.text-heading-xlarge`)
	if err == nil && nameEl != nil {
		fullName, _ := nameEl.Text()
		parts := strings.Fields(fullName)
		if len(parts) >= 1 {
			firstName = parts[0]
		}
		if len(parts) >= 2 {
			lastName = strings.Join(parts[1:], " ")
		}
	}

	// Try to get headline (job title)
	headlineEl, err := page.Timeout(2 * time.Second).Element(`.text-body-medium.break-words`)
	if err == nil && headlineEl != nil {
		jobTitle, _ = headlineEl.Text()
		jobTitle = strings.TrimSpace(jobTitle)
	}

	// Try to get company from experience section
	companyEl, err := page.Timeout(2 * time.Second).Element(`button[aria-label*="Current company"]`)
	if err == nil && companyEl != nil {
		company, _ = companyEl.Text()
		company = strings.TrimSpace(company)
	}

	return
}

// generateNote generates a personalized connection note
func (cm *ConnectionManager) generateNote(req *ConnectionRequest) string {
	if len(cm.templates) == 0 {
		return ""
	}

	// Select template (rotate through templates)
	template := cm.templates[req.TemplateIdx%len(cm.templates)]

	// Substitute variables
	vars := map[string]string{
		"firstName": req.FirstName,
		"lastName":  req.LastName,
		"jobTitle":  req.JobTitle,
		"company":   req.Company,
	}

	note := stealth.SubstituteTemplate(template, vars)

	// Enforce character limit
	if len(note) > cm.config.MaxNoteLength {
		note = note[:cm.config.MaxNoteLength-3] + "..."
	}

	return note
}

// sendWithNote sends a connection request with a personalized note
func (cm *ConnectionManager) sendWithNote(page *rod.Page, note string) error {
	// Click "Add a note" button
	addNoteBtn, err := page.Timeout(5 * time.Second).Element(`button[aria-label="Add a note"]`)
	if err != nil {
		// Try alternate approach - just find the note field
		noteField, err := page.Timeout(3 * time.Second).Element(`textarea[name="message"]`)
		if err != nil {
			// No note option available, send without note
			return cm.sendWithoutNote(page)
		}
		return cm.typeAndSend(page, noteField, note)
	}

	addNoteBtn.Click(proto.InputMouseButtonLeft, 1)
	time.Sleep(500 * time.Millisecond)

	// Find note textarea
	noteField, err := page.Timeout(3 * time.Second).Element(`textarea[name="message"], textarea#custom-message`)
	if err != nil {
		return fmt.Errorf("note field not found: %w", err)
	}

	return cm.typeAndSend(page, noteField, note)
}

// typeAndSend types the note and sends the request
func (cm *ConnectionManager) typeAndSend(page *rod.Page, noteField *rod.Element, note string) error {
	// Type note with realistic behavior
	sequence := cm.typing.GenerateTypingSequence(note)

	for _, char := range sequence {
		if char.IsBurstPause {
			time.Sleep(char.Delay)
			continue
		}

		if char.IsBackspace {
			noteField.MustInput("\b")
			time.Sleep(char.Delay)
			continue
		}

		noteField.MustInput(string(char.Char))
		time.Sleep(char.Delay)
	}

	// Think time before sending
	time.Sleep(cm.timing.GetThinkTime())

	// Click Send button
	sendBtn, err := page.Timeout(3 * time.Second).Element(`button[aria-label="Send now"], button[aria-label="Send invitation"]`)
	if err != nil {
		// Try generic send button
		sendBtn, err = page.Element(`button.ml1[aria-label*="Send"]`)
		if err != nil {
			return fmt.Errorf("send button not found: %w", err)
		}
	}

	sendBtn.Click(proto.InputMouseButtonLeft, 1)
	time.Sleep(time.Second)

	return nil
}

// sendWithoutNote sends a connection request without a note
func (cm *ConnectionManager) sendWithoutNote(page *rod.Page) error {
	// Click Send button directly
	selectors := []string{
		`button[aria-label="Send now"]`,
		`button[aria-label="Send without a note"]`,
		`button[aria-label="Send invitation"]`,
		`button.artdeco-button--primary[aria-label*="Send"]`,
	}

	for _, selector := range selectors {
		sendBtn, err := page.Timeout(2 * time.Second).Element(selector)
		if err == nil && sendBtn != nil {
			sendBtn.Click(proto.InputMouseButtonLeft, 1)
			time.Sleep(time.Second)
			return nil
		}
	}

	return fmt.Errorf("send button not found")
}

// clickWithRealism clicks an element with natural mouse movement
func (cm *ConnectionManager) clickWithRealism(page *rod.Page, element *rod.Element) error {
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

	// Pre-click hover actions
	hoverActions := cm.mouse.GeneratePreClickSequence(centerX, centerY, 1920, 1080)
	for _, action := range hoverActions {
		page.Mouse.MustMoveTo(action.X, action.Y)
		time.Sleep(action.Duration)
	}

	// Generate Bezier path to target
	path := cm.bezier.GeneratePath(0, 0, centerX, centerY)
	durations := cm.bezier.GetMovementDurations(len(path), 300*time.Millisecond)

	for i, point := range path {
		if i > 0 && i-1 < len(durations) {
			time.Sleep(durations[i-1])
		}
		page.Mouse.MustMoveTo(point.X, point.Y)
	}

	// Small hover delay
	time.Sleep(50 * time.Millisecond)

	return page.Mouse.Click(proto.InputMouseButtonLeft, 1)
}

// CanSendMoreToday checks if we can send more connection requests today
func (cm *ConnectionManager) CanSendMoreToday() (bool, int, error) {
	activity, err := cm.db.GetOrCreateDailyActivity()
	if err != nil {
		return false, 0, err
	}

	remaining := cm.config.DailyLimit - activity.ConnectionsSent
	return remaining > 0, remaining, nil
}
