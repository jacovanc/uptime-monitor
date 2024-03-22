package mail

// DummyEmailSender is a dummy implementation of EmailSender interface
type DummyEmailSender struct{
    SendEmailCalled int // Number of times SendEmail was called
}

func (d *DummyEmailSender) SendEmail(recipients []string, subject, body string) error {
    d.SendEmailCalled++

    return nil
}