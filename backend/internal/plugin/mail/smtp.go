package mail

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"mime"
	"net"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SmtpConfig holds the configuration for SMTP mail sender.
type SmtpConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	FromName    string `json:"from_name"`
	SystemEmail string `json:"system_email"`
	SmtpSecure  string `json:"smtp_secure"` // ssl, tls, or empty for plain
	Charset     string `json:"charset"`
}

// SmtpSender implements MailSender using Go's net/smtp package.
type SmtpSender struct {
	config SmtpConfig
}

// NewSmtpSender creates a new SmtpSender with the given configuration.
func NewSmtpSender(config SmtpConfig) *SmtpSender {
	if config.Charset == "" {
		config.Charset = "UTF-8"
	}
	if config.SmtpSecure == "" {
		config.SmtpSecure = "none"
	}
	return &SmtpSender{config: config}
}

// Send sends an email via SMTP with optional attachments.
func (s *SmtpSender) Send(to, subject, content string, attachments []string) error {
	if to == "" {
		return fmt.Errorf("recipient address is empty")
	}

	from := s.config.SystemEmail
	if from == "" {
		from = s.config.Username
	}

	fromHeader := from
	if s.config.FromName != "" {
		fromHeader = fmt.Sprintf("%s <%s>", mime.QEncoding.Encode(s.config.Charset, s.config.FromName), from)
	}

	// Build MIME message
	boundary := fmt.Sprintf("_boundary_%d", time.Now().UnixNano())
	var msg bytes.Buffer

	// Headers
	msg.WriteString(fmt.Sprintf("From: %s\r\n", fromHeader))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", mime.QEncoding.Encode(s.config.Charset, subject)))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))

	if len(attachments) > 0 {
		msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
		msg.WriteString("\r\n")
		msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	}

	msg.WriteString(fmt.Sprintf("Content-Type: text/html; charset=%s\r\n", s.config.Charset))
	msg.WriteString("Content-Transfer-Encoding: base64\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(base64.StdEncoding.EncodeToString([]byte(content)))
	msg.WriteString("\r\n")

	// Attachments
	for _, attachPath := range attachments {
		if attachPath == "" {
			continue
		}
		data, err := os.ReadFile(attachPath)
		if err != nil {
			return fmt.Errorf("failed to read attachment %s: %w", attachPath, err)
		}
		filename := filepath.Base(attachPath)
		msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		msg.WriteString(fmt.Sprintf("Content-Type: application/octet-stream; name=%s\r\n", mime.QEncoding.Encode(s.config.Charset, filename)))
		msg.WriteString("Content-Transfer-Encoding: base64\r\n")
		msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\r\n", mime.QEncoding.Encode(s.config.Charset, filename)))
		msg.WriteString("\r\n")
		msg.WriteString(base64.StdEncoding.EncodeToString(data))
		msg.WriteString("\r\n")
	}

	if len(attachments) > 0 {
		msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	}

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	recipients := []string{to}

	switch strings.ToLower(s.config.SmtpSecure) {
	case "ssl":
		return s.sendWithTLS(addr, from, recipients, msg.Bytes())
	case "tls":
		return s.sendWithSTARTTLS(addr, from, recipients, msg.Bytes())
	default:
		if s.config.Username != "" && s.config.Password != "" {
			auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
			return smtp.SendMail(addr, auth, from, recipients, msg.Bytes())
		}
		return smtp.SendMail(addr, nil, from, recipients, msg.Bytes())
	}
}

// sendWithTLS sends email over a direct TLS connection (typically port 465).
func (s *SmtpSender) sendWithTLS(addr, from string, recipients []string, msg []byte) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("invalid address %s: %w", addr, err)
	}

	tlsConfig := &tls.Config{
		ServerName: host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS dial failed: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("SMTP client creation failed: %w", err)
	}
	defer client.Close()

	if s.config.Username != "" && s.config.Password != "" {
		auth := smtp.PlainAuth("", s.config.Username, s.config.Password, host)
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP auth failed: %w", err)
		}
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("SMTP MAIL FROM failed: %w", err)
	}
	for _, r := range recipients {
		if err = client.Rcpt(r); err != nil {
			return fmt.Errorf("SMTP RCPT TO failed for %s: %w", r, err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA failed: %w", err)
	}
	if _, err = writer.Write(msg); err != nil {
		return fmt.Errorf("SMTP write failed: %w", err)
	}
	if err = writer.Close(); err != nil {
		return fmt.Errorf("SMTP close writer failed: %w", err)
	}

	return client.Quit()
}

// sendWithSTARTTLS sends email using STARTTLS (typically port 587).
func (s *SmtpSender) sendWithSTARTTLS(addr, from string, recipients []string, msg []byte) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("invalid address %s: %w", addr, err)
	}

	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("SMTP dial failed: %w", err)
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName: host,
		}
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("STARTTLS failed: %w", err)
		}
	}

	if s.config.Username != "" && s.config.Password != "" {
		auth := smtp.PlainAuth("", s.config.Username, s.config.Password, host)
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP auth failed: %w", err)
		}
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("SMTP MAIL FROM failed: %w", err)
	}
	for _, r := range recipients {
		if err = client.Rcpt(r); err != nil {
			return fmt.Errorf("SMTP RCPT TO failed for %s: %w", r, err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA failed: %w", err)
	}
	if _, err = writer.Write(msg); err != nil {
		return fmt.Errorf("SMTP write failed: %w", err)
	}
	if err = writer.Close(); err != nil {
		return fmt.Errorf("SMTP close writer failed: %w", err)
	}

	return client.Quit()
}
