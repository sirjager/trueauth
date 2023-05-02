package mail

type Config struct {
	SMTPSender string `mapstructure:"GMAIL_NAME"`
	SMTPUser   string `mapstructure:"GMAIL_USER"`
	SMTPPass   string `mapstructure:"GMAIL_PASS"`
}

// Mail represents an email message.
type Mail struct {
	To      []string // List of recipients' email addresses.
	Cc      []string // List of carbon copy recipients' email addresses.
	Bcc     []string // List of blind carbon copy recipients' email addresses.
	Subject string   // Email subject.
	Body    string   // Email body.
	Files   []string // List of file paths to attach to the email.
}

// MailSender is an interface for sending emails.
type MailSender interface {
	SendMail(mail Mail) error
}
