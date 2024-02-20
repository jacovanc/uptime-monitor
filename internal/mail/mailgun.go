package mail

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

// MailgunSender implements EmailSender for Mailgun
type MailgunSender struct {
    Domain string
    APIKey string
}

func CreateMailgunSender(domain string, apiKey string) EmailSender {
	return &MailgunSender{
		Domain: domain,
		APIKey: apiKey,
	}
}

func (m *MailgunSender) SendEmail(to []string, subject string, body string) error {
    mg := mailgun.NewMailgun(m.Domain, m.APIKey)
    message := mg.NewMessage("sender@example.com", subject, body, to...)
    
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
    defer cancel()

    _, _, err := mg.Send(ctx, message)
    return err
}
