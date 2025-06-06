package main

import (
	"context"
	"fmt"
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
	if len(os.Args) > 1 && os.Args[1] == "--local" {
		c, err := client.New(context.Background(), "")
		if err != nil {
			log.Fatalf("failed to create clusterctl client: %v", err)
		}
		initOptions := client.InitOptions{
			CoreProvider:            "cluster-api",
			BootstrapProviders:      []string{"kubeadm"},
			ControlPlaneProviders:   []string{"kubeadm"},
			InfrastructureProviders: []string{"docker:v1.10.1"},
		}
		if _, err = c.Init(context.Background(), initOptions); err != nil {
			log.Fatalf("Failed to initialize providers: %v", err)
		}
		fmt.Println("✅ Providers initialized.")
		templateOptions = client.GetClusterTemplateOptions{
			ProviderRepositorySource: &client.ProviderRepositorySourceOptions{
				InfrastructureProvider: "docker",
				Flavor:                 "development",
			},
		}
	} else {
		templateOptions = client.GetClusterTemplateOptions{
			URLSource: &client.URLSourceOptions{
				URL: "https://github.com/Tomy2e/cluster-api-provider-scaleway/releases/download/v0.0.3/cluster-template.yaml",
			},
		}
	}

	kubeconfigPath := internal.GetKubeconfig()
	log.Printf("Using kubeconfig: %s", kubeconfigPath)
	clusterSvc := internal.NewClusterService(kubeconfigPath, templateOptions)
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
