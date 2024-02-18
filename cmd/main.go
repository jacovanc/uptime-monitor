package main

import (
	"log"
	"os"
	"uptime-monitor/internal/monitor"
	"uptime-monitor/internal/storage"
	"uptime-monitor/internal/webserver"

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
	monitor, err := monitor.NewMonitor(sqliteStorer)
	if err != nil {
		log.Fatal("Error creating monitor: ", err)
	}

	go monitor.Start()

	webserver.Start()
}