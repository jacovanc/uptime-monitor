package monitor

import (
	"log"
	"net/http"
	"time"
)

// DataStorer is the interface that wraps methods to store monitoring data.
type DataStorer interface {
    StoreWebsiteStatus(website string, status string, latency time.Duration) error
}
type Monitor struct {
	ds DataStorer
	isRunning bool
}

// Config - TODO move to config
const intervalSeconds time.Duration = 5
var websites = []string{
	"https://www.google.com",
	"https://www.facebook.com",
}

func NewMonitor(ds DataStorer) *Monitor {
	return &Monitor{
		ds: ds,
		isRunning: false,
	}
}

func (m *Monitor) Start() {
    m.isRunning = true
    for m.isRunning {
        for _, website := range websites {
			// Ping inside a goroutine to prevent blocking the loop
            go func(website string) {
                log.Println("Pinging website", website)
                m.pingWebsite(website)
            }(website)
        }
        time.Sleep(intervalSeconds * time.Second)
    }
}

func (m *Monitor) Stop() {
	if(m.isRunning) {
		m.isRunning = false
	}
}

func (m *Monitor) pingWebsite(website string) {
	start := time.Now()
	
	// Ping website
	resp, err := http.Get(website)

	latency := time.Since(start)

	if(err != nil) {
		// There was an issue pinging the website, it is probably down
		m.ds.StoreWebsiteStatus(website, "down", latency)
		return
	}
	defer resp.Body.Close()

	m.ds.StoreWebsiteStatus(website, "up", latency)
}