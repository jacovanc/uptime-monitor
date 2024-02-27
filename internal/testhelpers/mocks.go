package testhelpers

import "time"

// DummyDataStorer is a dummy implementation of DataStorer interface
type DummyDataStorer struct{}

func (d *DummyDataStorer) StoreWebsiteStatus(website string, statusCode int, latency time.Duration) error {
    // Dummy implementation
    return nil
}

// DummyEmailSender is a dummy implementation of EmailSender interface
type DummyEmailSender struct{}

func (d *DummyEmailSender) SendEmail(recipients []string, subject, body string) error {
    // Dummy implementation
    return nil
}