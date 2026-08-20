package mail

// MailSender defines the interface for sending emails.
// All mail plugins (smtp, submail, alimail, etc.) must implement this interface.
type MailSender interface {
	Send(to, subject, content string, attachments []string) error
}
