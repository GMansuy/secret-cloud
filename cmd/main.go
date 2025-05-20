package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Tomy2e/cluster-api-provider-scaleway/api"
	"github.com/Tomy2e/cluster-api-provider-scaleway/internal"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
)

func main() {
	var templateOptions client.GetClusterTemplateOptions
	if os.Args[1] == "--local" {
		templateOptions = client.GetClusterTemplateOptions{
			ProviderRepositorySource: &client.ProviderRepositorySourceOptions{
				InfrastructureProvider: "docker",
				Flavor:                 "development"},
		}
	} else {
		templateOptions = client.GetClusterTemplateOptions{
			URLSource: &client.URLSourceOptions{
				URL: "https://github.com/Tomy2e/cluster-api-provider-scaleway/releases/download/v0.0.3/cluster-template.yaml",
			},
		}
	}

	clusterSvc := internal.NewClusterService(os.Getenv("KUBECONFIG"), templateOptions)
	app := api.NewApp(clusterSvc)
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
