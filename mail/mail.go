package mail

type Mail struct {
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body    string
	Files   []string
}

type MailSender interface {
	SendMail(mail Mail) error
}
