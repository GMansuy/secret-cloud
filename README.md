# Secret Cloud - Rocket App

This repository contains a Go application (Rocket App) and Kubernetes deployment manifests for deploying the application to a Kubernetes cluster.

## Project Structure

- `/cmd`: Contains the main Go application code
- `/api`: API definitions
- `/internal`: Internal packages
- `/config`: Configuration files
- `/kubernetes`: Kubernetes deployment manifests

## Deploying to Kubernetes

The application can be deployed to Kubernetes using the manifests in the `kubernetes` directory. See [Kubernetes Deployment Documentation](kubernetes/README.md) for detailed instructions.

## Building the Docker Image

The application can be built into a Docker image using the Dockerfile in the root directory:

```bash
docker build -t your-registry/rocket-app:latest .
```

## Development

To run the application locally:

```bash
go run cmd/main.go
```

## License

[License information]