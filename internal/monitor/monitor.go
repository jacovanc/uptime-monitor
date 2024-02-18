package monitor

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// DataStorer is the interface that wraps methods to store monitoring data.
type DataStorer interface {
    StoreWebsiteStatus(website string, status string, latency time.Duration) error
}
type Monitor struct {
	ds DataStorer
	isRunning bool
	websites []string
	interval time.Duration
}

// Config - TODO move to config
const defaultInterval time.Duration = 5 * time.Second

func NewMonitor(ds DataStorer) (*Monitor, error) {
	// Config
	var websitesString = os.Getenv("WEBSITES")
	websites := strings.Split(websitesString, ",")

	if len(websites) == 0 {
		// return error
		return nil, errors.New("no websites to monitor")
	}

	var intervalString = os.Getenv("INTERVAL_SECONDS")
	var interval time.Duration

	if intervalString == "" {
		interval = defaultInterval
	} else {
		var err error
		interval, err = time.ParseDuration(intervalString + "s")
		if err != nil {
			return nil, errors.New("invalid interval_seconds configuration")
		}
	}

	return &Monitor{
		ds: ds,
		isRunning: false,
		websites: websites,
		interval: interval,
	}, nil
}

func (m *Monitor) Start() {
	fmt.Println(m.interval, m.websites)
    m.isRunning = true
    for m.isRunning {
        for _, website := range m.websites {
			// Ping inside a goroutine to prevent blocking the loop
            go func(website string) {
                log.Println("Pinging website", website)
                m.pingWebsite(website)
            }(website)
        }
        time.Sleep(m.interval * time.Second)
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