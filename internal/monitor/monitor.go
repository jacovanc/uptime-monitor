package monitor

import (
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
    StoreWebsiteStatus(website string, status string, latency time.Duration) error
}

type Monitor struct {
	ds DataStorer
	es mail.EmailSender
	isRunning bool
	
	websites []string
	interval time.Duration
	downAlertThreshold int
	alertEmails []string

	statusHistory map[string][]string
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
		statusHistory: make(map[string][]string),
		latencyHistory: make(map[string][]time.Duration),
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
                status, latency := m.pingWebsite(website)

				m.updateStatusHistory(website, status)

				if m.shouldSendDownAlert(website) {
					log.Println("Sending down alert for", website)
					err := m.es.SendEmail(m.alertEmails, "Website down", website + " is down!")
					if err != nil {
						log.Println("Error sending down alert: ", err)
					}
				}

				err := m.ds.StoreWebsiteStatus(website, status, latency)
				if err != nil {
					log.Println("Error storing website status: ", err)
				}


            }(website)
        }
        time.Sleep(m.interval)
    }
}

func (m *Monitor) Stop() {
	if(m.isRunning) {
		m.isRunning = false
	}
}

func (m *Monitor) pingWebsite(website string) (status string, latency time.Duration) {
	start := time.Now()
	
	// Ping website
	resp, err := http.Get(website)
	latency = time.Since(start)

	if(err != nil) {
		log.Println("Error pinging: ", err)

		return "down", latency // Don't return the error as it's it's not important to the caller
	}
	defer resp.Body.Close()

	return "up", latency
}

func (m *Monitor) updateStatusHistory(website string, status string) {
	// Update the status history
	m.statusHistory[website] = append(m.statusHistory[website], status)

	// If the history is longer than the downAlertThreshold * 2, remove the oldest entry (no longer relevant)
	if len(m.statusHistory[website]) > m.downAlertThreshold * 2 {
		m.statusHistory[website] = m.statusHistory[website][1:]
	}
}

func (m *Monitor) shouldSendDownAlert(website string) bool {
	history := m.statusHistory[website]

	if len(m.statusHistory[website]) < m.downAlertThreshold {
		return false
	}

	for _, status := range history[len(history) - m.downAlertThreshold:] {
		if status != "down" {
			return false
		}
	}

	return true
}