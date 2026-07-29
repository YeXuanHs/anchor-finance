package email

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"gorm.io/gorm"
)

// Config holds SMTP configuration read from message_configs table.
type Config struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	From       string `json:"from"`        // Sender address
	FromName   string `json:"from_name"`   // Sender display name
	UseTLS     bool   `json:"use_tls"`     // Use STARTTLS
	UseSSL     bool   `json:"use_ssl"`     // Use SSL/TLS
	SkipVerify bool   `json:"skip_verify"` // Skip TLS certificate verification
}

// Sender sends emails via SMTP.
type Sender struct {
	db *gorm.DB
}

// NewSender creates a new email Sender.
func NewSender(db *gorm.DB) *Sender {
	return &Sender{db: db}
}

// loadConfig reads email config from the message_configs table (channel = "email").
func (s *Sender) loadConfig() (*Config, error) {
	var row struct {
		Provider  string `gorm:"column:provider"`
		Config    []byte `gorm:"column:config"`
		SenderName string `gorm:"column:sender_name"`
		SenderAddr string `gorm:"column:sender_addr"`
		IsEnabled bool   `gorm:"column:is_enabled"`
		Status    int16  `gorm:"column:status"`
	}

	err := s.db.Table("message_configs").
		Select("provider, config, sender_name, sender_addr, is_enabled, status").
		Where("channel = ?", "email").
		First(&row).Error
	if err != nil {
		return nil, fmt.Errorf("email config not found: %w", err)
	}

	if !row.IsEnabled || row.Status != 1 {
		return nil, errors.New("email channel is disabled")
	}

	cfg := &Config{}

	if len(row.Config) > 0 {
		if err := json.Unmarshal(row.Config, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse email config: %w", err)
		}
	}

	// Sender address from column takes precedence
	if row.SenderAddr != "" {
		cfg.From = row.SenderAddr
	}
	if row.SenderName != "" {
		cfg.FromName = row.SenderName
	}

	return cfg, nil
}

// Send sends an email with the given subject and HTML body.
func (s *Sender) Send(to, subject, htmlBody string) error {
	cfg, err := s.loadConfig()
	if err != nil {
		return err
	}

	if cfg.Host == "" {
		return errors.New("smtp host is required")
	}
	if cfg.Port == 0 {
		cfg.Port = 587
	}

	// Build message
	from := cfg.From
	if from == "" {
		from = cfg.Username
	}

	fromHeader := from
	if cfg.FromName != "" {
		fromHeader = fmt.Sprintf("%s <%s>", cfg.FromName, from)
	}

	msg := buildMessage(fromHeader, to, subject, htmlBody)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	if cfg.UseSSL {
		return s.sendWithTLS(addr, auth, from, to, msg, cfg.SkipVerify)
	}

	// Try STARTTLS first, fall back to plain
	err = s.sendWithSTARTTLS(addr, auth, from, to, msg, cfg.SkipVerify)
	if err != nil {
		// Fall back to plain SMTP
		return smtp.SendMail(addr, auth, from, []string{to}, msg)
	}

	return nil
}

// sendWithSTARTTLS sends email using STARTTLS.
func (s *Sender) sendWithSTARTTLS(addr string, auth smtp.Auth, from, to string, msg []byte, skipVerify bool) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	host, _, _ := net.SplitHostPort(addr)
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Start TLS if supported
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: skipVerify,
		}
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("STARTTLS failed: %w", err)
		}
	}

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP auth failed: %w", err)
		}
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("SMTP MAIL FROM failed: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("SMTP RCPT TO failed: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA failed: %w", err)
	}

	_, err = writer.Write(msg)
	if err != nil {
		return fmt.Errorf("SMTP write failed: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("SMTP close writer failed: %w", err)
	}

	return client.Quit()
}

// sendWithTLS sends email using implicit TLS (port 465).
func (s *Sender) sendWithTLS(addr string, auth smtp.Auth, from, to string, msg []byte, skipVerify bool) error {
	host, _, _ := net.SplitHostPort(addr)

	tlsConfig := &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: skipVerify,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server with TLS: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP auth failed: %w", err)
		}
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("SMTP MAIL FROM failed: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("SMTP RCPT TO failed: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA failed: %w", err)
	}

	_, err = writer.Write(msg)
	if err != nil {
		return fmt.Errorf("SMTP write failed: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("SMTP close writer failed: %w", err)
	}

	return client.Quit()
}

// buildMessage constructs a MIME email message.
func buildMessage(from, to, subject, htmlBody string) []byte {
	var msg strings.Builder

	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	return []byte(msg.String())
}
