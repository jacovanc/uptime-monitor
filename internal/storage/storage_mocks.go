package storage

import "time"

type DummyDataStorer struct{}

func (s *DummyDataStorer) StoreWebsiteStatus(website string, statusCode int, latency time.Duration) error {
    // Dummy implementation
    return nil
}
