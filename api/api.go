package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Tomy2e/cluster-api-provider-scaleway/internal"
	"github.com/Tomy2e/cluster-api-provider-scaleway/internal/scope"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type App struct {
	Router                   *chi.Mux
	ManagementKubeconfigPath string
}

func NewApp() *App {
	app := &App{
		Router:                   chi.NewRouter(),
		ManagementKubeconfigPath: os.Getenv("KUBECONFIG"),
	}
	app.setupRoutes()
	return app
}

func (a *App) setupRoutes() {
	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.Recoverer)

	a.Router.Get("/", a.homeHandler)
	a.Router.Post("/cluster", a.createClusterHandler)
	a.Router.Delete("/cluster", a.deleteClusterHandler)
}

func (a *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Clusterapi server HTTP server!"))
}

func (a *App) createClusterHandler(w http.ResponseWriter, r *http.Request) {
	var cluster scope.Cluster
	if err := json.NewDecoder(r.Body).Decode(&cluster); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := internal.CreateCluster(context.Background(), cluster.Name, cluster.ControlplaneMachineCount, cluster.WorkerMachineCount, a.ManagementKubeconfigPath)
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

func (a *App) deleteClusterHandler(w http.ResponseWriter, r *http.Request) {
	var cluster scope.Cluster
	if err := json.NewDecoder(r.Body).Decode(&cluster); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := internal.DeleteCluster(context.Background(), cluster.Name, a.ManagementKubeconfigPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate cluster config: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"status": "success",
		"message": "Cluster deleting",
		"name":    cluster.Name,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
