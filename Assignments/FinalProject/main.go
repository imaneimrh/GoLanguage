package main

import (
	. "FinalProject/logging"
	. "FinalProject/reports"
	. "FinalProject/routes"
	. "FinalProject/utils"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	logFile, err := SetupLogging()
	if err != nil {
		log.Fatalf("Failed to set up logging: %v", err)
	}
	defer logFile.Close()

	router, bookStore, authorStore, customerStore, orderStore := InitializeRoutes()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go StartSalesReportBackgroundJob(ctx, orderStore, bookStore, 3*time.Hour) //24*time.Hour

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Server running on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :8080: %v\n", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	SaveAllData(ctx, bookStore, authorStore, customerStore, orderStore)

	log.Println("Server exited cleanly")
}
