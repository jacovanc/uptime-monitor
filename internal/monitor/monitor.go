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
	"uptime-monitor/internal/mail"
)

// DataStorer is the interface that wraps methods to store monitoring data.
type DataStorer interface {
    StoreWebsiteStatus(website string, statusCode int, latency time.Duration) error
}

type Monitor struct {
	ds DataStorer
	es mail.EmailSender
	isRunning bool
	cancelFunc context.CancelFunc
	
	websites []string
	interval time.Duration
	downAlertThreshold int
	alertEmails []string

	statusHistory map[string][]int
	latencyHistory map[string][]time.Duration
}

// Config - TODO move to config
const defaultInterval time.Duration = 60 * time.Second

func NewMonitor(ds DataStorer, es mail.EmailSender) (*Monitor, error) {
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

	var downAlertThresholdString = os.Getenv("DOWN_ALERT_THRESHOLD")
	var downAlertThreshold int
	downAlertThreshold, err := strconv.Atoi(downAlertThresholdString)
	if err != nil {
		return nil, errors.New("invalid down_alert_threshold configuration")
	}

	var alertEmailsString = os.Getenv("ALERT_EMAILS")
	alertEmails := strings.Split(alertEmailsString, ",")
	if len(alertEmails) == 0 {
		return nil, errors.New("no alert emails configured")
	}

	log.Println("Creating monitor with interval", interval, "and websites", websites)
	log.Println("Down alert threshold is", downAlertThreshold, "and alert emails are", alertEmails)

	// TODO: Set up the status history using the DB historical data rather than assuming empty

	return &Monitor{
		ds: ds,
		es: es,
		isRunning: false,

		// Config
		websites: websites,
		interval: interval,
		downAlertThreshold: downAlertThreshold,
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
		go func(website string) {
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
		}(website)
	}
}

// Stop the monitoring goroutines for each website.
// uses the cancel function to stop the goroutines.
func (m *Monitor) Stop() {
	if(m.isRunning && m.cancelFunc != nil) {
		// Call the cancel function to stop the goroutines
		m.cancelFunc()
		m.isRunning = false
	}
}

// Handles the event inside the monitoring goroutine loop for a website.
// Includes pinging the website, updating the status history, and sending alerts if necessary.
func (m *Monitor) tick(website string) {
	log.Println("Pinging website", website)
	statusCode, latency := m.pingWebsite(website)

	m.appendStatusHistory(website, statusCode)

	err := m.ds.StoreWebsiteStatus(website, statusCode, latency)
	if err != nil {
		log.Println("Error storing website status: ", err)
	}

	if m.shouldSendDownAlert(website) {
		log.Println("Sending down alert for", website)
		err := m.es.SendEmail(m.alertEmails, "Website down", website + " is down!")
		if err != nil {
			log.Println("Error sending down alert: ", err)
		}

		// Reset the status history to prevent sending multiple alerts
		m.statusHistory[website] = []int{}
	}
}

// Pings a website and returns the status and latency.
func (m *Monitor) pingWebsite(website string) (status int, latency time.Duration) {
	log.Println("hit", website)

	start := time.Now()
	var statusCode int

	// Ping website
	resp, err := http.Get(website)
	latency = time.Since(start)

	log.Println("Pinged", website, "in", latency)
	if(err != nil) {
		log.Println("Error pinging: ", err)
		statusCode = 0
		return statusCode, latency
	}
	defer resp.Body.Close()
	
	// Get status code
	statusCode = resp.StatusCode

	return statusCode, latency
}

// Appends the status to the status history and removes the oldest entry if far enough out of bounds for the required logic.
func (m *Monitor) appendStatusHistory(website string, statusCode int) {
	// Update the status history
	m.statusHistory[website] = append(m.statusHistory[website], statusCode)

	// If the history is longer than the downAlertThreshold * 2, remove the oldest entry (no longer relevant)
	if len(m.statusHistory[website]) > m.downAlertThreshold * 2 {
		m.statusHistory[website] = m.statusHistory[website][1:]
	}
}

// Checks if the website has been down for the downAlertThreshold and should send an alert.
func (m *Monitor) shouldSendDownAlert(website string) bool {
	history := m.statusHistory[website]

	if len(m.statusHistory[website]) < m.downAlertThreshold {
		return false
	}

	for _, statusCode := range history[len(history) - m.downAlertThreshold:] {
		if statusCode > 499 || statusCode == 0 { // 0 represents an error on the ping (e.g. if the website does not exist)
			return true
		}
	}

	return false
}