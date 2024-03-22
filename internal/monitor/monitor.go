package monitor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type DataStorer interface {
    StoreWebsiteStatus(website string, statusCode int, latency time.Duration) error
}

type EmailSender interface {
	SendEmail(to []string, subject string, body string) error
}

// Pings a website and returns the status and latency.
func pingWebsite(website string) (status int, latency time.Duration) {
	start := time.Now()
	var statusCode int

	// Ping website
	resp, err := http.Head(website)
	latency = time.Since(start)

	if(err != nil) {
		log.Println("Error pinging: ", err)
		statusCode = 0
		return statusCode, latency
	}
	log.Println("Pinged", website, "in", latency)

	defer resp.Body.Close()
	
	// Get status code
	statusCode = resp.StatusCode

	return statusCode, latency
}

type Monitor struct {
	ds DataStorer
	es EmailSender
	isRunning bool
	cancelFunc context.CancelFunc
	
	websites []string
	interval time.Duration
	alertThreshold int
	alertEmails []string

	statusHistory map[string][]int
	latencyHistory map[string][]time.Duration
}

// Config - TODO move to config
const defaultInterval time.Duration = 60 * time.Second

func NewMonitor(ds DataStorer, es EmailSender) (*Monitor, error) {
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

	alertThreshold, err := strconv.Atoi(os.Getenv("DOWN_ALERT_THRESHOLD"))
	if err != nil {
		return nil, errors.New("invalid down_alert_threshold configuration")
	}

	alertEmails := strings.Split(os.Getenv("ALERT_EMAILS"), ",")
	if len(alertEmails) == 0 {
		return nil, errors.New("no alert emails configured")
	}

	log.Println("Creating monitor with interval", interval, "and websites", websites)
	log.Println("Alert threshold is", alertThreshold, "and alert emails are", alertEmails)

	return &Monitor{
		ds: ds,
		es: es,
		isRunning: false,

		// Config
		websites: websites,
		interval: interval,
		alertThreshold: alertThreshold,
		alertEmails: alertEmails,

		// State
		statusHistory: make(map[string][]int),
		latencyHistory: make(map[string][]time.Duration),
	}, nil
}

// Start the monitoring goroutines for each website.
func (m *Monitor) Start() {
	fmt.Println(m.interval, m.websites)
    m.isRunning = true

	ctx, cancel := context.WithCancel(context.Background())
    m.cancelFunc = cancel // Store the cancel function to call it later

	for _, website := range m.websites {
		go m.monitorWebsite(ctx, website)
	}

	go m.enableAlerts(ctx)
}

// Stop the monitoring goroutines for each website.
// uses the cancel function to stop the goroutines.
func (m *Monitor) Stop() {
	if(m.isRunning && m.cancelFunc != nil) {
		// Call the cancel function to stop the goroutines
		m.cancelFunc()
		m.isRunning = false
		m.cancelFunc = nil
	}
}

// The monitoring goroutine for a website. Runs until the context says to cancel.
func (m *Monitor) monitorWebsite(ctx context.Context, website string) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping monitor for", website)
			return
		default:
			m.tick(website)
			time.Sleep(m.interval)
		}
	}
}

// The alerting goroutine. Runs until the context says to cancel.
// Fires alerts for any websites that have been down for the alertThreshold.
func (m *Monitor) enableAlerts(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			m.sendAlerts()
			time.Sleep(m.interval * 2) // Allow a bit of extra time before firing alerts to batch them
		}
	}
}

// Handles the event inside the monitoring goroutine loop for a website.
// Includes pinging the website, updating the status history, and sending alerts if necessary.
func (m *Monitor) tick(website string) {
	log.Println("Pinging website", website)
	statusCode, latency := pingWebsite(website)

	m.appendStatusHistory(website, statusCode)

	err := m.ds.StoreWebsiteStatus(website, statusCode, latency)
	if err != nil {
		log.Println("Error storing website status: ", err)
	}
}

// Appends the status to the status history and removes the oldest entry if far enough out of bounds for the required logic.
func (m *Monitor) appendStatusHistory(website string, statusCode int) {
	// Update the status history
	m.statusHistory[website] = append(m.statusHistory[website], statusCode)

	// If the history is longer than the alertThreshold * 2, remove the oldest entry (no longer relevant)
	if len(m.statusHistory[website]) > m.alertThreshold * 2 {
		m.statusHistory[website] = m.statusHistory[website][1:]
	}
}

// Checks if the website has been down for the alertThreshold and should send an alert.
func (m *Monitor) shouldSendAlert(website string) bool {
	history := m.statusHistory[website]

	if len(m.statusHistory[website]) < m.alertThreshold {
		return false
	}

	shouldSendAlert := true
	for _, statusCode := range history[len(history) - m.alertThreshold:] {
		if statusCode < 500 && statusCode != 0 { // If there is any non-500 status code, don't send an alert
			shouldSendAlert = false
		}
	}

	return shouldSendAlert
}

// Check all websites for down alerts and send an email if necessary.
// Batches the alerts to send all in one email
func (m *Monitor) sendAlerts() {
	var alerts []string

	for _, website := range m.websites {
		if m.shouldSendAlert(website) {
			alerts = append(alerts, website)

			// Reset the status history to prevent sending multiple alerts
			m.statusHistory[website] = []int{}
		}
	}

	if len(alerts) > 0 {
		subject := "Website alert"
		body := "The following websites are down:\n" + strings.Join(alerts, "\n")
		err := m.es.SendEmail(m.alertEmails, subject, body)
		if err != nil {
			log.Println("Error sending down alert email: ", err)
		}
	}
}
