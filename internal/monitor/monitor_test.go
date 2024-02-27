package monitor

import (
	"os"
	"testing"
	"uptime-monitor/internal/mail"
	"uptime-monitor/internal/storage"
)

func TestCreateMonitor(t *testing.T) {
	os.Setenv("WEBSITES", "http://example.com,http://example2.com")
	os.Setenv("INTERVAL_SECONDS", "30")
	os.Setenv("DOWN_ALERT_THRESHOLD", "3")
	os.Setenv("ALERT_EMAILS", "alert@example.com")

	mockDataStorer := &storage.DummyDataStorer{}
	mockEmailSender := &mail.DummyEmailSender{}

	monitor, err := NewMonitor(mockDataStorer, mockEmailSender)
	if err != nil {
		t.Error("Expected no error, got", err)
	}

	// Check monitor.interval of type time.Duration has 30s
	if monitor.interval.Seconds() != 30 {
		t.Error("Expected interval to be 30s, got", monitor.interval)
	}
	
	// Check monitor.downAlertThreshold is 3
	if monitor.downAlertThreshold != 3 {
		t.Error("Expected downAlertThreshold to be 3, got", monitor.downAlertThreshold)
	}

	// Check monitor.alertEmails
	if len(monitor.alertEmails) != 1 || monitor.alertEmails[0] != "alert@example.com" {
		t.Error("Expected alertEmails to be array with one element of value 'alert@example.com', got", monitor.alertEmails)
	}

	// Check monitor.websites
	if len(monitor.websites) != 2 || monitor.websites[0] != "http://example.com" || monitor.websites[1] != "http://example2.com" {
		t.Error("Expected websites to be array with two elements of value 'http://example.com' and 'http://example2.com', got", monitor.websites)
	}
}

func TestMonitorStartStop(t *testing.T) {
	os.Setenv("WEBSITES", "https://example.com")
	os.Setenv("INTERVAL_SECONDS", "1")
	os.Setenv("DOWN_ALERT_THRESHOLD", "3")
	os.Setenv("ALERT_EMAILS", "alert@example.com")

	mockDataStorer := &storage.DummyDataStorer{}
	mockEmailSender := &mail.DummyEmailSender{}

	monitor, err := NewMonitor(mockDataStorer, mockEmailSender)
	if err != nil {
		t.Error("Expected no error, got", err)
	}

	monitor.Start()

	if !monitor.isRunning {
		t.Error("Expected monitor to be running")
	}

	if monitor.cancelFunc == nil {
		t.Error("Expected cancelFunc to be set")
	}

	monitor.Stop()

	if monitor.isRunning {
		t.Error("Expected monitor to be stopped")
	}

	if monitor.cancelFunc != nil {
		t.Error("Expected cancelFunc to be nil")
	}
}

func TestAppendStatusHistory(t *testing.T) {
	website := "https://example.com"

	os.Setenv("WEBSITES", website)
	os.Setenv("INTERVAL_SECONDS", "1")
	os.Setenv("DOWN_ALERT_THRESHOLD", "3")
	os.Setenv("ALERT_EMAILS", "alert@example.com")

	mockDataStorer := &storage.DummyDataStorer{}
	mockEmailSender := &mail.DummyEmailSender{}

	monitor, err := NewMonitor(mockDataStorer, mockEmailSender)
	if err != nil {
		t.Error("Expected no error, got", err)
	}

	monitor.appendStatusHistory(website, 200)
	monitor.appendStatusHistory(website, 500)
	monitor.appendStatusHistory(website, 200)

	if len(monitor.statusHistory[website]) != 3 {
		t.Error("Expected statusHistory to have 3 elements, got", len(monitor.statusHistory[website]))
	}

	if monitor.statusHistory[website][0] != 200 || monitor.statusHistory[website][1] != 500 || monitor.statusHistory[website][2] != 200 {
		t.Error("Expected statusHistory to be [200, 500, 200], got", monitor.statusHistory[website])
	}

	// Append until we are past 2*downAlertThreshold to ensure that it trims the history
	for i := 0; i < 10; i++ {
		monitor.appendStatusHistory(website, 200)
	}
	// Check that the history length is 2*downAlertThreshold
	if len(monitor.statusHistory[website]) != 6 {
		t.Error("Expected statusHistory to have 6 elements, got", len(monitor.statusHistory[website]))
	}

	// Check that all elements currently have status code 200
	for _, status := range monitor.statusHistory[website] {
		if status != 200 {
			t.Error("Expected all elements to be 200, got", status)
		}
	}
}

func TestShouldSendDownAlert(t *testing.T) {
	website := "https://example.com"

	os.Setenv("WEBSITES", website)
	os.Setenv("INTERVAL_SECONDS", "1")
	os.Setenv("DOWN_ALERT_THRESHOLD", "3")
	os.Setenv("ALERT_EMAILS", "alert@example.com")

	mockDataStorer := &storage.DummyDataStorer{}
	mockEmailSender := &mail.DummyEmailSender{}

	monitor, err := NewMonitor(mockDataStorer, mockEmailSender)
	if err != nil {
		t.Error("Expected no error, got", err)
	}

	// Append 3 status codes 500
	monitor.appendStatusHistory(website, 500)
	monitor.appendStatusHistory(website, 500)
	monitor.appendStatusHistory(website, 500)

	// Should return true
	if !monitor.shouldSendDownAlert(website) {
		t.Error("Expected shouldSendDownAlert to return true")
	}
}

func TestPingWebsite(t *testing.T) {
	// Placeholder for this test.
	// Probably want to use a mock for the http get request, instead of actually pinging a website
	// Otherwise our test could fail if the website is down, even if our code is correct.
	// Also we can't test various status codes, as we can't control the response from the website
}

// Also, once we have mocked the http get request, we can test the Monitor.tick function (and determine if all the correct dependency mocks were called correctly for different scenarios)
