package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tomy2e/cluster-api-provider-scaleway/internal"
	"github.com/Tomy2e/cluster-api-provider-scaleway/internal/scope"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	Router *chi.Mux
}

func NewApp() *App {
	app := &App{
		Router: chi.NewRouter(),
	}
	app.setupRoutes()
	return app
}

func (a *App) setupRoutes() {
	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.Recoverer)

	a.Router.Get("/", a.homeHandler)
	a.Router.Post("/cluster", a.clusterHandler)
}

func (a *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Clusterapi server HTTP server!"))
}

func (a *App) clusterHandler(w http.ResponseWriter, r *http.Request) {
	var cluster scope.Cluster
	if err := json.NewDecoder(r.Body).Decode(&cluster); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	kubeconfigPath := filepath.Join("kubeconfig-AMMI-CAAPH.yaml")
	err := internal.GenerateClusterConfigFromMemory(context.TODO(), cluster.Name, cluster.ControlplaneMachineCount, cluster.WorkerMachineCount, kubeconfigPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate cluster config: %v", err), http.StatusInternalServerError)
		return
	}
	response := map[string]string{"status": "success",
		"message":                  "Cluster received",
		"name":                     cluster.Name,
		"controlplaneMachineCount": fmt.Sprintf("%v", cluster.ControlplaneMachineCount),
		"workerMachineCount":       fmt.Sprintf("%v", cluster.WorkerMachineCount)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	app := NewApp()

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
