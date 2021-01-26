package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"webmonitor/monitor"
	"webmonitor/notifier"
	"webmonitor/storage"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//Checking that an environment variable is present or not.
	// webhook, ok := os.LookupEnv("WEBHOOK")
	// if !ok {
	// 	log.Fatal("You must set the WEBHOOK environment variable.")
	// }

	sender, ok := os.LookupEnv("SENDER_EMAIL")
	if !ok {
		log.Fatal("You must set the SENDER_EMAIL environment variable.")
	}

	sendgridApiKey, ok := os.LookupEnv("SENDGRID_API_KEY")
	if !ok {
		log.Fatal("You must set the SENDGRID_API_KEY environment variable.")
	}

	storage := storage.NewMemoryStorage()
	notifier := notifier.NewEmailNotifier(sender, sendgridApiKey)
	monitor := monitor.NewMonitor(storage, notifier)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	monitor.Start()

	<-c
	log.Println("Received SIGTERM, stopping server")
	monitor.Stop()
}
