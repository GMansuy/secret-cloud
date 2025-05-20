package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Tomy2e/cluster-api-provider-scaleway/api"
)

func main() {
	app := api.NewApp()

	server := &http.Server{
		Addr:    ":8080",
		Handler: app.Router,
	}

	go func() {
		log.Println("Starting server on :8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :8080: %v\n", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
