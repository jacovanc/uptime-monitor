package monitor

import (
	"os"
	"testing"
	"uptime-monitor/internal/testhelpers"
)

func TestCreateMonitor(t *testing.T) {
	os.Setenv("WEBSITES", "http://example.com")
	os.Setenv("INTERVAL_SECONDS", "30")
	os.Setenv("DOWN_ALERT_THRESHOLD", "3")
	os.Setenv("ALERT_EMAILS", "alert@example.com")

	mockDataStorer := testhelpers.DummyEmailSender{}
	mockEmailSender := testhelpers.DummyEmailSender{}

	monitor, err := NewMonitor(mockDataStorer, mockEmailSender)
	if err != nil {
		t.Error("Expected no error, got", err)
	}

	if monitor.interval != 30 {
		t.Error("Expected interval to be 30, got", monitor.interval)
	}
	

}
