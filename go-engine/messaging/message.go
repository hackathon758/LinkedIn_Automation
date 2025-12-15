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

// MessageManager handles sending follow-up messages
type MessageManager struct {
	config    config.MessagingConfig
	db        *database.DB
	logger    *logger.Logger
	timing    *stealth.TimingController
	typing    *stealth.TypingSimulator
	templates []string
}

// NewMessageManager creates a new MessageManager
func NewMessageManager(
	cfg config.MessagingConfig,
	db *database.DB,
	log *logger.Logger,
	stealthCfg config.StealthConfig,
) *MessageManager {
	return &MessageManager{
		config:    cfg,
		db:        db,
		logger:    log.WithComponent("messaging"),
		timing:    stealth.NewTimingController(stealthCfg.Timing),
		typing:    stealth.NewTypingSimulator(stealthCfg.Timing),
		templates: cfg.Templates,
	}
}

// MessageRequest represents a message to send
type MessageRequest struct {
	ConnectionID string
	ProfileURL   string
	FirstName    string
	LastName     string
	JobTitle     string
	Company      string
	Message      string
	TemplateIdx  int
}

// MessageResult represents the result of sending a message
type MessageResult struct {
	Success      bool
	ConnectionID string
	ErrorMessage string
}

// SendMessage sends a follow-up message to an accepted connection
func (mm *MessageManager) SendMessage(page *rod.Page, req *MessageRequest) (*MessageResult, error) {
	mm.logger.Info("sending message", "connection", req.ConnectionID, "profile", req.ProfileURL)

	// Navigate to profile
	err := page.Navigate(req.ProfileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate: %w", err)
	}

	time.Sleep(mm.timing.GetPageLoadDelay())
	page.WaitLoad()

	// Find and click Message button
	messageBtn, err := mm.findMessageButton(page)
	if err != nil {
		return &MessageResult{
			Success:      false,
			ConnectionID: req.ConnectionID,
			ErrorMessage: "Message button not found - may not be connected",
		}, nil
	}

	messageBtn.Click(proto.InputMouseButtonLeft, 1)
	time.Sleep(time.Second)

	// Wait for messaging pane to open
	messageInput, err := page.Timeout(5 * time.Second).Element(`div.msg-form__contenteditable, textarea.msg-form__textarea`)
	if err != nil {
		return &MessageResult{
			Success:      false,
			ConnectionID: req.ConnectionID,
			ErrorMessage: "Message input not found",
		}, nil
	}

	// Generate message if not provided
	if req.Message == "" && len(mm.templates) > 0 {
		req.Message = mm.generateMessage(req)
	}

	if req.Message == "" {
		return &MessageResult{
			Success:      false,
			ConnectionID: req.ConnectionID,
			ErrorMessage: "No message content",
		}, nil
	}

	// Type message with realistic behavior
	err = mm.typeMessage(messageInput, req.Message)
	if err != nil {
		return &MessageResult{
			Success:      false,
			ConnectionID: req.ConnectionID,
			ErrorMessage: fmt.Sprintf("Failed to type message: %v", err),
		}, nil
	}

	// Think time before sending
	time.Sleep(mm.timing.GetThinkTime())

	// Click send
	err = mm.clickSend(page)
	if err != nil {
		return &MessageResult{
			Success:      false,
			ConnectionID: req.ConnectionID,
			ErrorMessage: fmt.Sprintf("Failed to send: %v", err),
		}, nil
	}

	// Record message in database
	msg := &database.Message{
		ID:           fmt.Sprintf("msg_%d", time.Now().UnixNano()),
		ConnectionID: req.ConnectionID,
		Content:      req.Message,
		Status:       "sent",
		SentAt:       time.Now(),
	}
	mm.db.SaveMessage(msg)
	mm.db.IncrementMessageCount()

	mm.logger.Info("message sent", "connection", req.ConnectionID)

	return &MessageResult{
		Success:      true,
		ConnectionID: req.ConnectionID,
	}, nil
}

// findMessageButton finds the Message button on a profile page
func (mm *MessageManager) findMessageButton(page *rod.Page) (*rod.Element, error) {
	selectors := []string{
		`button[aria-label*="Message"]`,
		`a[href*="/messaging/"]`,
		`button.message-anywhere-button`,
		`button:has-text("Message")`,
	}

	for _, selector := range selectors {
		btn, err := page.Timeout(3 * time.Second).Element(selector)
		if err == nil && btn != nil {
			visible, _ := btn.Visible()
			if visible {
				return btn, nil
			}
		}
	}

	return nil, fmt.Errorf("message button not found")
}

// generateMessage generates a personalized message from template
func (mm *MessageManager) generateMessage(req *MessageRequest) string {
	if len(mm.templates) == 0 {
		return ""
	}

	// Select template
	template := mm.templates[req.TemplateIdx%len(mm.templates)]

	// Substitute variables
	vars := map[string]string{
		"firstName": req.FirstName,
		"lastName":  req.LastName,
		"jobTitle":  req.JobTitle,
		"company":   req.Company,
	}

	return stealth.SubstituteTemplate(template, vars)
}

// typeMessage types a message with realistic behavior
func (mm *MessageManager) typeMessage(element *rod.Element, message string) error {
	sequence := mm.typing.GenerateTypingSequence(message)

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

// clickSend clicks the send button
func (mm *MessageManager) clickSend(page *rod.Page) error {
	selectors := []string{
		`button[type="submit"].msg-form__send-button`,
		`button.msg-form__send-button`,
		`button[aria-label="Send"]`,
	}

	for _, selector := range selectors {
		sendBtn, err := page.Timeout(2 * time.Second).Element(selector)
		if err == nil && sendBtn != nil {
			return sendBtn.Click(proto.InputMouseButtonLeft, 1)
		}
	}

	return fmt.Errorf("send button not found")
}

// DetectAcceptedConnections checks for newly accepted connections
func (mm *MessageManager) DetectAcceptedConnections(page *rod.Page) ([]database.Connection, error) {
	// Get pending connections from database
	pending, err := mm.db.GetPendingConnections()
	if err != nil {
		return nil, err
	}

	if len(pending) == 0 {
		return nil, nil
	}

	// Navigate to My Network
	err = page.Navigate("https://www.linkedin.com/mynetwork/invite-connect/connections/")
	if err != nil {
		return nil, err
	}

	time.Sleep(mm.timing.GetPageLoadDelay())
	page.WaitLoad()

	// Get connection list
	html, err := page.HTML()
	if err != nil {
		return nil, err
	}

	// Check each pending connection
	var accepted []database.Connection
	for _, conn := range pending {
		// Simple check - if profile URL appears in connections list, it's accepted
		profileID := extractProfileID(conn.ProfileURL)
		if strings.Contains(html, profileID) {
			// Update status in database
			mm.db.UpdateConnectionStatus(conn.ProfileURL, "accepted")
			conn.Status = "accepted"
			accepted = append(accepted, conn)
			mm.logger.Info("connection accepted", "profile", conn.ProfileURL)
		}
	}

	return accepted, nil
}

// extractProfileID extracts profile ID from URL
func extractProfileID(profileURL string) string {
	// Extract the ID portion from /in/profile-id/
	parts := strings.Split(profileURL, "/in/")
	if len(parts) < 2 {
		return ""
	}
	id := strings.TrimSuffix(parts[1], "/")
	return id
}

// GetConnectionsNeedingFollowUp returns accepted connections without follow-up messages
func (mm *MessageManager) GetConnectionsNeedingFollowUp() ([]database.Connection, error) {
	accepted, err := mm.db.GetAcceptedConnections()
	if err != nil {
		return nil, err
	}

	var needFollowUp []database.Connection
	for _, conn := range accepted {
		hasSent, err := mm.db.HasSentFollowUp(conn.ID)
		if err != nil {
			continue
		}
		if !hasSent {
			needFollowUp = append(needFollowUp, conn)
		}
	}

	return needFollowUp, nil
}

// CanSendMoreMessagesToday checks if we can send more messages today
func (mm *MessageManager) CanSendMoreMessagesToday() (bool, int, error) {
	activity, err := mm.db.GetOrCreateDailyActivity()
	if err != nil {
		return false, 0, err
	}

	remaining := mm.config.DailyLimit - activity.MessagesSent
	return remaining > 0, remaining, nil
}
