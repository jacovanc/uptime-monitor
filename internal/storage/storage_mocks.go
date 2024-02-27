package storage

import "time"

// DummyDataStorer is a dummy implementation of DataStorer interface
type DummyDataStorer struct{}

func (s *DummyDataStorer) StoreWebsiteStatus(website string, statusCode int, latency time.Duration) error {
    // Dummy implementation
    return nil
}
