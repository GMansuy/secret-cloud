package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	b64 "encoding/base64"
	"github.com/Tomy2e/cluster-api-provider-scaleway/internal"
	"github.com/Tomy2e/cluster-api-provider-scaleway/internal/scope"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type App struct {
	Router     *chi.Mux
	ClusterSvc *internal.ClusterService
}

func NewApp(c *internal.ClusterService) *App {
	app := &App{
		Router:     chi.NewRouter(),
		ClusterSvc: c,
	}
	app.setupRoutes()
	return app
}

func (a *App) setupRoutes() {
	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.Recoverer)
	a.Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	a.Router.Get("/", a.homeHandler)
	a.Router.Post("/cluster", a.createClusterHandler)
	a.Router.Delete("/cluster", a.deleteClusterHandler)
	a.Router.Get("/list", a.listClusters)
	a.Router.Get("/cluster", a.getCluster)
	a.Router.Get("/cluster/{name}/kubeconfig", a.getClusterKubeconfigHandler)
}

func (a *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Clusterapi server HTTP server!"))
}

func (a *App) createClusterHandler(w http.ResponseWriter, r *http.Request) {
	var cluster scope.Cluster
	if err := json.NewDecoder(r.Body).Decode(&cluster); err != nil {
		log.Printf("Error decoding JSON in createClusterHandler: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	opts := a.ClusterSvc.SetCreationTemplateOptions(&cluster)

	err := a.ClusterSvc.CreateCluster(context.Background(), cluster.Name, opts)
	if err != nil {
		log.Printf("Error creating cluster %s: %v", cluster.Name, err)
		http.Error(w, fmt.Sprintf("Failed to generate cluster config: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"status": "success",
		"message":                  "Cluster received",
		"name":                     cluster.Name,
		"controlplaneMachineCount": fmt.Sprintf("%v", cluster.ControlplaneMachineCount),
		"workerMachineCount":       fmt.Sprintf("%v", cluster.WorkerMachineCount)}
	json.NewEncoder(w).Encode(response)
}

func (a *App) deleteClusterHandler(w http.ResponseWriter, r *http.Request) {
	var cluster scope.Cluster
	if err := json.NewDecoder(r.Body).Decode(&cluster); err != nil {
		log.Printf("Error decoding JSON in deleteClusterHandler: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := a.ClusterSvc.DeleteCluster(context.Background(), cluster.Name)
	if err != nil {
		log.Printf("Error deleting cluster %s: %v", cluster.Name, err)
		http.Error(w, fmt.Sprintf("Failed to generate cluster config: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"status": "success",
		"message": "Cluster deleting",
		"name":    cluster.Name,
	}
	json.NewEncoder(w).Encode(response)
}

func (a *App) listClusters(w http.ResponseWriter, r *http.Request) {

	cmd := exec.Command("kubectl", "get", "cluster", `-o=jsonpath='{range .items[*]}{.metadata.name}{.status.conditions}{end}`)
	ouput, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error listing clusters: %v", err)
		http.Error(w, fmt.Sprintf("Failed to list clusters: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status":   "success",
		"clusters": string(ouput),
	}
	json.NewEncoder(w).Encode(response)
}

func (a *App) getCluster(w http.ResponseWriter, r *http.Request) {

	clusterName := "mang2"

	cmd := exec.Command("kubectl", "get", "cluster", clusterName, `-o=jsonpath={.status.conditions}`)
	ouput, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error getting cluster %s: %v", clusterName, err)
		http.Error(w, fmt.Sprintf("Failed to list clusters: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"name":   clusterName,
		"status": string(ouput),
	}
	json.NewEncoder(w).Encode(response)
}

func (a *App) getClusterKubeconfigHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the name parameter from the URL
	clusterName := chi.URLParam(r, "name")

	if clusterName == "" {
		log.Printf("Error: Empty cluster name in getClusterKubeconfigHandler")
		http.Error(w, "Cluster name is required", http.StatusBadRequest)
		return
	}

	// Use the extracted cluster name to get the kubeconfig

	cmd := exec.Command("kubectl", "get", "secret", fmt.Sprintf("%s-kubeconfig", clusterName), "-o", "jsonpath={.data.value}")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error getting kubeconfig for cluster %s: %v", clusterName, err)
		http.Error(w, fmt.Sprintf("Failed to get kubeconfig: %v", err), http.StatusInternalServerError)
		return
	}
	sDec, _ := b64.StdEncoding.DecodeString(string(output))

	response := map[string]string{
		"name":       clusterName,
		"kubeconfig": string(sDec),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
