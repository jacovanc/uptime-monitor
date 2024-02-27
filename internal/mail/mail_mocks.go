package mail

// DummyEmailSender is a dummy implementation of EmailSender interface
type DummyEmailSender struct{}

func (d *DummyEmailSender) SendEmail(recipients []string, subject, body string) error {
    // Dummy implementation
    return nil
}