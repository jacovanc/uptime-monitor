package mail

type EmailSender interface {
	SendEmail(to []string, subject string, body string) error
}
