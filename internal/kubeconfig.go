package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// ScalewayResponse represents the response from Scaleway API
// We only need the clusterID from the first cluster
type ScalewayResponse struct {
	Clusters []struct {
		ID string `json:"id"`
	} `json:"clusters"`
}

// KubeconfigResponse represents the response from the kubeconfig endpoint
type KubeconfigResponse struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// makeKubeconfigRequest makes a request to the kubeconfig endpoint, saves it to a file, and returns the file path
func makeKubeconfigRequest(client *http.Client, url, authToken string) (string, error) {
	kubeconfigReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating kubeconfig request: %v\n", err)
		return "", err
	}

	kubeconfigReq.Header.Add("X-Auth-Token", authToken)
	kubeconfigResp, err := client.Do(kubeconfigReq)
	if err != nil {
		log.Printf("Error making kubeconfig request: %v\n", err)
		return "", err
	}

	defer kubeconfigResp.Body.Close()
	kubeconfigBody, err := ioutil.ReadAll(kubeconfigResp.Body)
	if err != nil {
		log.Printf("Error reading kubeconfig response body: %v\n", err)
		return "", err
	}

	// Parse the JSON response
	var kubeconfigResponse KubeconfigResponse
	if err := json.Unmarshal(kubeconfigBody, &kubeconfigResponse); err != nil {
		log.Printf("Error parsing kubeconfig JSON: %v\n", err)
		return "", err
	}

	// Decode the base64-encoded content
	decodedContent, err := base64.StdEncoding.DecodeString(kubeconfigResponse.Content)
	if err != nil {
		log.Printf("Error decoding base64 content: %v\n", err)
		return "", err
	}

	// Write the decoded content to the file
	kubeconfigPath := "kubeconfig.yaml"
	err = ioutil.WriteFile(kubeconfigPath, decodedContent, 0644)
	if err != nil {
		log.Printf("Error writing kubeconfig to file: %v\n", err)
		return "", err
	}

	// Get absolute path
	absPath, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %v\n", err)
		return kubeconfigPath, nil
	}

	fullPath := fmt.Sprintf("%s/%s", absPath, kubeconfigPath)
	log.Printf("Kubeconfig written to: %s\n", fullPath)

	return fullPath, nil
}

// GetKubeconfig retrieves the kubeconfig from Scaleway API and returns the path to the saved kubeconfig file
func GetKubeconfig() string {
	// Define structs to parse the Scaleway API response
	req, err := http.NewRequest("GET", "https://api.scaleway.com/k8s/v1/regions/fr-par/clusters", nil)
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return ""
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", os.Getenv("SCW_SECRET_KEY"))
	q := req.URL.Query()
	q.Add("param", "capi-cloud")
	req.URL.RawQuery = q.Encode()
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error making request: %v\n", err)
		return ""
	}

	defer resp.Body.Close()
	log.Printf("Request returned status: %s\n", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return ""
	}
	var response ScalewayResponse
	if err := json.Unmarshal(body, &response); err == nil && len(response.Clusters) > 0 {
		log.Printf("Found cluster ID from struct: %s\n", response.Clusters[0].ID)
	}

	// If we found a cluster ID, get the kubeconfig
	if response.Clusters[0].ID != "" {
		kubeconfigURL := fmt.Sprintf("https://api.scaleway.com/k8s/v1/regions/fr-par/clusters/%s/kubeconfig", response.Clusters[0].ID)
		log.Printf("Kubeconfig URL: %s\n", kubeconfigURL)
		kubeconfigPath, err := makeKubeconfigRequest(httpClient, kubeconfigURL, os.Getenv("SCW_SECRET_KEY"))
		if err != nil {
			log.Printf("Error getting kubeconfig: %v\n", err)
			return ""
		}
		return kubeconfigPath
	} else {
		log.Println("No cluster ID found in the response")
		return ""
	}
}
