# Cluster configuration
CLUSTER_NAME := capi-management
KUBECONFIG_PATH := $(PWD)/.kubeconfig

.PHONY: all
all: cluster load-env

# Create the kind cluster
.PHONY: cluster
cluster:
	@echo "Creating $(CLUSTER_NAME) cluster..."
	@if ! kind get clusters | grep -q "^$(CLUSTER_NAME)$$"; then \
		kind create cluster --name $(CLUSTER_NAME); \
	else \
		echo "Cluster $(CLUSTER_NAME) already exists"; \
	fi
	@kind get kubeconfig --name $(CLUSTER_NAME) > $(KUBECONFIG_PATH)
	@echo "Kubeconfig saved to $(KUBECONFIG_PATH)"

# Load the kubeconfig environment variable
.PHONY: load-env
load-env:
	@echo "Setting KUBECONFIG environment variable..."
	@echo "export KUBECONFIG=$(KUBECONFIG_PATH)"

# Run the application with KUBECONFIG set
.PHONY: run
run: cluster
	KUBECONFIG=$(KUBECONFIG_PATH) go run cmd/main.go

# Initialize CAPI providers
.PHONY: init-providers
run-local: cluster
	CLUSTER_TOPOLOGY=true KUBECONFIG=$(KUBECONFIG_PATH) go run cmd/main.go --local

# Clean up resources
.PHONY: clean
clean:
	@echo "Cleaning up resources..."
	@if kind get clusters | grep -q "^$(CLUSTER_NAME)$$"; then \
		kind delete cluster --name $(CLUSTER_NAME); \
	fi
	@rm -f $(KUBECONFIG_PATH)
	@rm -f *-cluster.yaml

# Help command
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make              - Create cluster and set kubeconfig"
	@echo "  make cluster      - Create kind cluster"
	@echo "  make load-env     - Print command to export KUBECONFIG"
	@echo "  make init-providers - Initialize CAPI providers"
	@echo "  make run          - Run application with KUBECONFIG set"
	@echo "  make clean        - Delete cluster and temporary files"