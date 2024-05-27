package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"

	"github.com/sirjager/trueauth/pkg/validator"
)

// GmailSender is a struct that represents a Gmail email sender.
type GmailSender struct {
	address     string
	plainAuth   smtp.Auth
	senderName  string
	senderEmail string
}

const (
	// GmailSMTPHost is the SMTP host for Gmail.
	GmailSMTPHost = "smtp.gmail.com"
	// GmailSMTPPort is the SMTP port for Gmail.
	GmailSMTPPort = "587"
)

// NewGmailSender creates a new GmailSender instance with the provided GmailSMTP configuration.
func NewGmailSender(config Config) (Sender, error) {
	// Validate the SMTP user email.
	if err := validator.ValidateEmail(config.SMTPUser); err != nil {
		return nil, err
	}
	// Validate the SMTP password.
	if err := validator.ValidatePassword(config.SMTPPass); err != nil {
		return nil, err
	}

	return &GmailSender{
		senderEmail: config.SMTPUser,
		senderName:  config.SMTPSender,
		address:     GmailSMTPHost + ":" + GmailSMTPPort,
		plainAuth:   smtp.PlainAuth("", config.SMTPUser, config.SMTPPass, GmailSMTPHost),
	}, nil
}

// SendMail sends an email using the GmailSender.
func (gmail *GmailSender) SendMail(mail Mail) error {
	// Validate the recipient email addresses.
	for _, r := range mail.To {
		if err := validator.ValidateEmail(r); err != nil {
			return err
		}
	}
	// Validate the Bcc email addresses.
	for _, r := range mail.Bcc {
		if err := validator.ValidateEmail(r); err != nil {
			return err
		}
	}
	// Validate the Cc email addresses.
	for _, r := range mail.Cc {
		if err := validator.ValidateEmail(r); err != nil {
			return err
		}
	}

	// Create a new email instance.
	email := email.NewEmail()
	email.From = fmt.Sprintf("%s <%s>", gmail.senderName, gmail.senderEmail)
	email.To = mail.To
	email.Cc = mail.Cc
	email.Bcc = mail.Bcc
	email.Subject = mail.Subject
	email.HTML = []byte(mail.Body)

	// Attach files to the email.
	for _, f := range mail.Files {
		if _, err := email.AttachFile(f); err != nil {
			return fmt.Errorf("failed to attach file: %s :%w", f, err)
		}
	}

	// Send the email using the Gmail SMTP server and authentication.
	return email.Send(gmail.address, gmail.plainAuth)
}
