package messaging

import (
	"math/rand"
	"strings"
	"time"
)

// TemplateManager handles message template operations
type TemplateManager struct {
	connectionTemplates []string
	followUpTemplates   []string
	rng                 *rand.Rand
}

// NewTemplateManager creates a new TemplateManager
func NewTemplateManager(connectionTemplates, followUpTemplates []string) *TemplateManager {
	return &TemplateManager{
		connectionTemplates: connectionTemplates,
		followUpTemplates:   followUpTemplates,
		rng:                 rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// TemplateVariables holds variables for template substitution
type TemplateVariables struct {
	FirstName string
	LastName  string
	JobTitle  string
	Company   string
	Location  string
}

// GetRandomConnectionTemplate returns a random connection request template
func (tm *TemplateManager) GetRandomConnectionTemplate() string {
	if len(tm.connectionTemplates) == 0 {
		return ""
	}
	return tm.connectionTemplates[tm.rng.Intn(len(tm.connectionTemplates))]
}

// GetRandomFollowUpTemplate returns a random follow-up message template
func (tm *TemplateManager) GetRandomFollowUpTemplate() string {
	if len(tm.followUpTemplates) == 0 {
		return ""
	}
	return tm.followUpTemplates[tm.rng.Intn(len(tm.followUpTemplates))]
}

// Render renders a template with the given variables
func (tm *TemplateManager) Render(template string, vars TemplateVariables) string {
	result := template

	replacements := map[string]string{
		"{{firstName}}": vars.FirstName,
		"{{lastName}}":  vars.LastName,
		"{{jobTitle}}":  vars.JobTitle,
		"{{company}}":   vars.Company,
		"{{location}}":  vars.Location,
	}

	for placeholder, value := range replacements {
		if value != "" {
			result = strings.ReplaceAll(result, placeholder, value)
		} else {
			// Handle empty values gracefully
			result = strings.ReplaceAll(result, placeholder, "")
		}
	}

	// Clean up any double spaces from empty variables
	for strings.Contains(result, "  ") {
		result = strings.ReplaceAll(result, "  ", " ")
	}

	return strings.TrimSpace(result)
}

// RenderConnectionNote generates a personalized connection note
func (tm *TemplateManager) RenderConnectionNote(vars TemplateVariables, maxLength int) string {
	template := tm.GetRandomConnectionTemplate()
	if template == "" {
		return ""
	}

	result := tm.Render(template, vars)

	// Enforce character limit
	if maxLength > 0 && len(result) > maxLength {
		result = result[:maxLength-3] + "..."
	}

	return result
}

// RenderFollowUpMessage generates a personalized follow-up message
func (tm *TemplateManager) RenderFollowUpMessage(vars TemplateVariables) string {
	template := tm.GetRandomFollowUpTemplate()
	if template == "" {
		return ""
	}

	return tm.Render(template, vars)
}

// AddConnectionTemplate adds a new connection request template
func (tm *TemplateManager) AddConnectionTemplate(template string) {
	tm.connectionTemplates = append(tm.connectionTemplates, template)
}

// AddFollowUpTemplate adds a new follow-up message template
func (tm *TemplateManager) AddFollowUpTemplate(template string) {
	tm.followUpTemplates = append(tm.followUpTemplates, template)
}

// SetConnectionTemplates replaces all connection templates
func (tm *TemplateManager) SetConnectionTemplates(templates []string) {
	tm.connectionTemplates = templates
}

// SetFollowUpTemplates replaces all follow-up templates
func (tm *TemplateManager) SetFollowUpTemplates(templates []string) {
	tm.followUpTemplates = templates
}

// GetConnectionTemplateCount returns the number of connection templates
func (tm *TemplateManager) GetConnectionTemplateCount() int {
	return len(tm.connectionTemplates)
}

// GetFollowUpTemplateCount returns the number of follow-up templates
func (tm *TemplateManager) GetFollowUpTemplateCount() int {
	return len(tm.followUpTemplates)
}

// ValidateTemplate checks if a template has valid variable placeholders
func ValidateTemplate(template string) []string {
	var errors []string

	validVars := []string{"{{firstName}}", "{{lastName}}", "{{jobTitle}}", "{{company}}", "{{location}}"}

	// Find all placeholders in template
	for {
		start := strings.Index(template, "{{")
		if start == -1 {
			break
		}
		end := strings.Index(template[start:], "}}")
		if end == -1 {
			errors = append(errors, "Unclosed placeholder starting at position "+string(rune(start)))
			break
		}

		placeholder := template[start : start+end+2]
		isValid := false
		for _, valid := range validVars {
			if placeholder == valid {
				isValid = true
				break
			}
		}

		if !isValid {
			errors = append(errors, "Unknown variable: "+placeholder)
		}

		template = template[start+end+2:]
	}

	return errors
}
