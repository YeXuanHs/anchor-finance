package sms

import "errors"

// Common errors returned by SMS plugins.
var (
	ErrInvalidPhone     = errors.New("sms: invalid phone number")
	ErrEmptyContent     = errors.New("sms: message content is empty")
	ErrMissingConfig    = errors.New("sms: missing required configuration")
	ErrSendFailed       = errors.New("sms: send failed")
	ErrTemplateNotFound = errors.New("sms: template not found")
)

// PluginInfo describes an SMS plugin's metadata.
type PluginInfo struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Version     string `json:"version"`
	HelpURL     string `json:"help_url"`
}

// TemplateResult is the result of a template operation.
type TemplateResult struct {
	Status         string `json:"status"` // "success" or "error"
	TemplateID     string `json:"template_id,omitempty"`
	TemplateStatus int    `json:"template_status,omitempty"` // 1=pending, 2=approved
	Msg            string `json:"msg,omitempty"`
}

// SendResult is the result of a send operation.
type SendResult struct {
	Status  string `json:"status"` // "success" or "error"
	Content string `json:"content,omitempty"`
	Msg     string `json:"msg,omitempty"`
}

// SmsSender is the interface that all SMS plugins must implement.
type SmsSender interface {
	// Info returns plugin metadata.
	Info() PluginInfo

	// Send sends a plain text SMS message to the given mobile number.
	Send(mobile, content string) (*SendResult, error)

	// SendTemplate sends an SMS using a template with parameters.
	SendTemplate(mobile, templateID string, params map[string]string) (*SendResult, error)

	// GetTemplate queries the status of a template on the provider side.
	GetTemplate(templateID string) (*TemplateResult, error)

	// CreateTemplate creates a new template on the provider.
	CreateTemplate(title, content string) (*TemplateResult, error)

	// UpdateTemplate updates an existing template on the provider.
	UpdateTemplate(templateID, title, content string) (*TemplateResult, error)

	// DeleteTemplate deletes a template from the provider.
	DeleteTemplate(templateID string) (*TemplateResult, error)
}
