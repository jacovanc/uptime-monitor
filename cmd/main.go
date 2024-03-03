package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"uptime-monitor/internal/mail"
	"uptime-monitor/internal/monitor"
	"uptime-monitor/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	// Load config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbPath := os.Getenv("DB_PATH") // If this isn't set, the empty string will be replaced with a default in storage module
	
	sqliteStorer, err := storage.NewSQLiteStorer(dbPath)
	if err != nil {
		log.Fatal("Error creating SQLite storer: ", err)
	}
	
	emailSender := mail.CreateMailgunSender(
		os.Getenv("MAILGUN_DOMAIN"),
		os.Getenv("MAILGUN_API_KEY"),
		os.Getenv("MAILGUN_SENDER"),
	)

	monitor, err := monitor.NewMonitor(sqliteStorer, emailSender)
	if err != nil {
		log.Fatal("Error creating monitor: ", err)
	}

	go monitor.Start()

	// Block waiting for shutdown signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) // Listen for SIGINT and SIGTERM 
	<-sigs // Block until channel receives a signal
	
	fmt.Println("Shutting down")
	monitor.Stop()
	fmt.Println("Shut down");
}