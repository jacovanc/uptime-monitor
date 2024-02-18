package main

import (
	"log"
	"uptime-monitor/internal/monitor"
	"uptime-monitor/internal/storage"
	"uptime-monitor/internal/webserver"
)

func main() {
	sqliteStorer, err := storage.NewSQLiteStorer("")
	if err != nil {
		log.Fatal("Error creating SQLite storer: ", err)
	}
	monitor := monitor.NewMonitor(sqliteStorer)

	go monitor.Start()

	webserver.Start()
}