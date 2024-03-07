package mail

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

// MailgunSender implements EmailSender for Mailgun
type MailgunSender struct {
	mailgun mailgun.Mailgun
	sender string
}

func CreateMailgunSender(domain string, apiKey string, sender string) *MailgunSender {
	mg := mailgun.NewMailgun(domain, apiKey)
	mg.SetAPIBase(mailgun.APIBaseEU)

	return &MailgunSender{
		mailgun: mg,
		sender: sender,
	}
}

func (m *MailgunSender) SendEmail(to []string, subject string, body string) error {
    message := m.mailgun.NewMessage(m.sender, subject, body, to...)
    
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
    defer cancel()

    _, _, err := m.mailgun.Send(ctx, message)
    return err
}
