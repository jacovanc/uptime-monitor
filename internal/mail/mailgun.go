package mail

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

// MailgunSender implements EmailSender for Mailgun
type MailgunSender struct {
    domain string
    apiKey string
	sender string
}

func CreateMailgunSender(domain string, apiKey string, sender string) EmailSender {
	return &MailgunSender{
		domain: domain,
		apiKey: apiKey,
		sender: sender,
	}
}

func (m *MailgunSender) SendEmail(to []string, subject string, body string) error {
    mg := mailgun.NewMailgun(m.domain, m.apiKey)

	mg.SetAPIBase(mailgun.APIBaseEU)

    message := mg.NewMessage(m.sender, subject, body, to...)
    
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
    defer cancel()

    _, _, err := mg.Send(ctx, message)
    return err
}
