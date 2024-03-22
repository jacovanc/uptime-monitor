package monitor

import (
	"net/http"
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
	if monitor.alertThreshold != 3 {
		t.Error("Expected downAlertThreshold to be 3, got", monitor.alertThreshold)
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

func TestShouldSendAlert(t *testing.T) {
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

	monitor.appendStatusHistory(website, 500)
	monitor.appendStatusHistory(website, 500)

	// Should return false
	if monitor.shouldSendAlert(website) {
		t.Error("Expected shouldSendDownAlert to return true")
	}

	monitor.appendStatusHistory(website, 500)

	// Should return true
	if !monitor.shouldSendAlert(website) {
		t.Error("Expected shouldSendDownAlert to return true")
	}
}

func TestSendAlerts(t *testing.T) {
	website := "https://example.com"
	website2 := "https://example2.com"

	os.Setenv("WEBSITES", website + "," + website2)
	os.Setenv("INTERVAL_SECONDS", "1")
	os.Setenv("DOWN_ALERT_THRESHOLD", "3")
	os.Setenv("ALERT_EMAILS", "test@example.com")

	mockDataStorer := &storage.DummyDataStorer{}
	mockEmailSender := &mail.DummyEmailSender{}

	monitor, err := NewMonitor(mockDataStorer, mockEmailSender)
	if err != nil {
		t.Error("Expected no error, got", err)
	}

	// Append 3 status codes 500 for website 1
	monitor.appendStatusHistory(website, 500)
	monitor.appendStatusHistory(website, 500)
	monitor.appendStatusHistory(website, 500)

	monitor.sendAlerts()

	// Expect that the email sender was called once
	if mockEmailSender.SendEmailCalled != 1 {
		t.Error("Expected sendEmail to be called once, got", mockEmailSender.SendEmailCalled)
	}

	// Append 3 status codes 500 for website 2
	monitor.appendStatusHistory(website2, 500)
	monitor.appendStatusHistory(website2, 500)
	monitor.appendStatusHistory(website2, 500)

	// Expect that the email sender was called once (previous status history not cleared)
	if mockEmailSender.SendEmailCalled != 1 {
		t.Error("Expected sendEmail to be called once, got", mockEmailSender.SendEmailCalled)
	}
}


func TestPingWebsite(t *testing.T) {
	// Create a local server that returns 200 OK
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	go http.ListenAndServe(":8080", nil)

	activeWebsite := "http://localhost:8080" // Use our local server that we know is up and will return a 200 status code
	inactiveWebsite := "http://localhost:8081" // Use our local server on a port we know is not running and will fail

	statusCode, latency := pingWebsite(activeWebsite)

	if statusCode != 200 {
		t.Error("Expected status code to be 200, got", statusCode)
	}

	if latency.Milliseconds() <= 0 {
		t.Error("Expected latency to be greater than 0, got", latency)
	}

	statusCode, _ = pingWebsite(inactiveWebsite)

	if statusCode != 0 {
		t.Error("Expected status code to be 0, got", statusCode)
	}
}