package mail

// Config represents the email configuration.
type Config struct {
	SMTPSender string `mapstructure:"SMTP_NAME"`
	SMTPUser   string `mapstructure:"SMTP_USER"`
	SMTPPass   string `mapstructure:"SMTP_PASS"`
	Template   string
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

// Sender is an interface for sending emails.
type Sender interface {
	SendMail(mail Mail) error
}
