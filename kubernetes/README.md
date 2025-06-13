# Rocket App Kubernetes Deployment

This directory contains Kubernetes manifests for deploying the Rocket App to a Kubernetes cluster.

## Prerequisites

- Docker installed locally
- Access to a Kubernetes cluster
- kubectl configured to communicate with your cluster
- (Optional) kustomize installed

## Deployment Steps

### 1. Build and Push the Docker Image

From the root of the repository:

```bash
# Build the Docker image
docker build -t your-registry/rocket-app:latest .

# Push the image to your registry
docker push your-registry/rocket-app:latest
```

### 2. Update the Deployment Configuration

Edit `kubernetes/deployment.yaml` to use your Docker image:

```yaml
image: your-registry/rocket-app:latest
```

Alternatively, if using kustomize, uncomment and update the images section in `kustomization.yaml`:

```yaml
images:
- name: ${DOCKER_REGISTRY}/rocket-app
  newName: your-registry/rocket-app
  newTag: latest
```

### 3. Update the Ingress Configuration (if needed)

Edit `kubernetes/ingress.yaml` to use your domain:

```yaml
host: your-domain.com
```

### 4. Deploy to Kubernetes

Using kubectl directly:

```bash
kubectl apply -f kubernetes/deployment.yaml
kubectl apply -f kubernetes/service.yaml
kubectl apply -f kubernetes/ingress.yaml
```

Or using kustomize:

```bash
kubectl apply -k kubernetes/
```

### 5. Verify the Deployment

```bash
# Check the deployment status
kubectl get deployments

# Check the pods
kubectl get pods

# Check the service
kubectl get services

# Check the ingress
kubectl get ingress
```

## Configuration

The application is configured with:

- 2 replicas for high availability
- Resource limits to prevent resource starvation
- Health checks to ensure the application is running correctly

## Scaling

To scale the application:

```bash
kubectl scale deployment rocket-app --replicas=3
```

## Troubleshooting

If you encounter issues:

1. Check the pod logs:
   ```bash
   kubectl logs -l app=rocket-app
   ```

2. Check the pod status:
   ```bash
   kubectl describe pods -l app=rocket-app
   ```

3. Check the service:
   ```bash
   kubectl describe service rocket-app
   ```

4. Check the ingress:
   ```bash
   kubectl describe ingress rocket-app
   ```