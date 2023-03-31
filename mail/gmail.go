package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	_ "github.com/jordan-wright/email"
	"github.com/sirjager/trueauth/validator"
)

type GmailSender struct {
	address     string
	plainAuth   smtp.Auth
	senderName  string
	senderEmail string
}

func NewGmailSender(host, port, email, pass, senderName string) (MailSender, error) {
	if err := validator.ValidateEmail(email); err != nil {
		return nil, err
	}
	if err := validator.ValidatePassword(pass); err != nil {
		return nil, err
	}
	address := fmt.Sprintf("%s:%s", host, port)
	plainAuth := smtp.PlainAuth("", email, pass, host)
	return &GmailSender{
		address:     address,
		plainAuth:   plainAuth,
		senderEmail: email,
		senderName:  senderName,
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
