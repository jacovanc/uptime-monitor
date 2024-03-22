package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uptime-monitor/internal/mail"
	"uptime-monitor/internal/monitor"
	"uptime-monitor/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	// Load development configuration from .env file
	if os.Getenv("ENV") != "prod" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	sqliteStorer, err := storage.NewSQLiteStorer(os.Getenv("DB_PATH"))
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

	// Start a webserver that simply returns 200 OK for all requests
	// This is for deployment status checks, but will be extended later to provide more functionality
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{Addr: ":8080"}

	go func() {
		fmt.Println("Server is starting...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v\n", err)
		}
	}()
	
	log.Println("Setting up signal handlers...")
	// Blocking main and waiting for shutdown signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Waiting for shutdown signal...")
	<-sigs

	fmt.Println("Shutting down server...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	fmt.Println("Server gracefully stopped")

	// Stop your monitor or any other services below
	monitor.Stop()
	fmt.Println("Monitor stopped. Exiting now.")
}