package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/sirjager/trueauth/cfg"
	"github.com/sirjager/trueauth/validator"
)

type GmailSender struct {
	address     string
	plainAuth   smtp.Auth
	senderName  string
	senderEmail string
}

const (
	GMAIL_SMTP_HOST = "smtp.gmail.com"
	GMAIL_SMTP_PORT = "587"
)

func NewGmailSender(config cfg.GmailSMTP) (MailSender, error) {
	if err := validator.ValidateEmail(config.SMTPUser); err != nil {
		return nil, err
	}
	if err := validator.ValidatePassword(config.SMTPPass); err != nil {
		return nil, err
	}
	return &GmailSender{
		senderEmail: config.SMTPUser,
		senderName:  config.SMTPSender,
		address:     GMAIL_SMTP_HOST + ":" + GMAIL_SMTP_PORT,
		plainAuth:   smtp.PlainAuth("", config.SMTPUser, config.SMTPPass, GMAIL_SMTP_HOST),
	}, nil
}

func (gmail *GmailSender) SendMail(mail Mail) error {
	for _, r := range mail.To {
		if err := validator.ValidateEmail(r); err != nil {
			return err
		}
	}
	for _, r := range mail.Bcc {
		if err := validator.ValidateEmail(r); err != nil {
			return err
		}
	}
	for _, r := range mail.Cc {
		if err := validator.ValidateEmail(r); err != nil {
			return err
		}
	}

	email := email.NewEmail()
	email.From = fmt.Sprintf("%s <%s>", gmail.senderName, gmail.senderEmail)
	email.Subject = mail.Subject
	email.To = mail.To
	email.Cc = mail.Cc
	email.Bcc = mail.Bcc

	for _, f := range mail.Files {
		if _, err := email.AttachFile(f); err != nil {
			return fmt.Errorf("failed to attach file: %s :%w", f, err)
		}
	}

	return email.Send(gmail.address, gmail.plainAuth)
}
