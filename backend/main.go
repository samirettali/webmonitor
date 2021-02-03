package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/samirettali/webmonitor/api"
	"github.com/samirettali/webmonitor/middlewares"
	"github.com/samirettali/webmonitor/monitor"
	"github.com/samirettali/webmonitor/notifier"
	"github.com/samirettali/webmonitor/storage"
	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

func main() {
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
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

	postgreURI, ok := os.LookupEnv("POSTGRE_URI")
	if !ok {
		log.Fatal("You must set the POSTGRE_URI environment variable.")
	}

	checksTable, ok := os.LookupEnv("POSTGRE_CHECKS_TABLE")
	if !ok {
		log.Fatal("You must set the POSTGRE_CHECKS_TABLE environment variable.")
	}

	statusesTable, ok := os.LookupEnv("POSTGRE_STATUSES_TABLE")
	if !ok {
		log.Fatal("You must set the POSTGRE_STATUES_TABLE environment variable.")
	}


	sendgridApiKey, ok := os.LookupEnv("SENDGRID_API_KEY")
	if !ok {
		log.Fatal("You must set the SENDGRID_API_KEY environment variable.")
	}

	storage := &storage.PostgreStorage{
		URI:    postgreURI,
		ChecksTable:  checksTable,
		StatusesTable:  statusesTable,
		Logger: log,
	}

	if err != nil {
		log.Fatal(err)
	}

	notifier := notifier.NewEmailNotifier(sender, sendgridApiKey, log)
	monitor := monitor.NewMonitor(storage, notifier, log)

	if err := monitor.Start(); err != nil {
		log.Fatal("Could not start monitor: ", err)
	}
	log.Println("Monitor started")

	defer monitor.Stop()

	handler := api.StorageHandler{Storage: storage, Logger: log}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/checks", handler.GetChecks).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/checks", handler.CreateCheck).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/checks/{id}", handler.GetCheck).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/checks/{id}", handler.DeleteCheck).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/checks/{id}", handler.UpdateCheck).Methods(http.MethodPatch, http.MethodOptions)
	router.HandleFunc("/checks/{id}/history", handler.GetHistory).Methods(http.MethodGet, http.MethodOptions)
	router.Use(middlewares.Logger)

	h := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH"},
	}).Handler(router)

	srv := &http.Server{
		Addr:         "0.0.0.0:8000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 15,
		Handler:      h,
	}

	go func() {
		log.Println("Starting HTTP server")
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Printf("ListenAndServe error: %s\n", err)
		} else {
			log.Printf("Server shutdown correctly")
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Println("Received SIGTERM")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Error(err)
	}
	os.Exit(0)
}
