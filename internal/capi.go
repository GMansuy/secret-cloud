package internal

import (
	"context"
	"fmt"
	"os"

	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
)

func GenerateClusterConfig(ctx context.Context, clusterName string, controlPlaneCount, workerCount int, kubeConfigPath string) error {
	c, err := client.New(ctx, kubeConfigPath)
	if err != nil {
		return err
	}

	initOptions = client.InitOptions{
		Kubeconfig: client.Kubeconfig{
			Path:    kubeConfigPath,
			Context: "",
		},
		CoreProvider:            "cluster-api",
		BootstrapProviders:      []string{"kubadm"},
		InfrastructureProviders: nil, // add your infrastructure provider here (how to add scaleway)
		ControlPlaneProviders:   []string{"kubadm"},
		AddonProviders:          nil, // Can whe put caaph here?
		WaitProviders:           true,
		WaitProviderTimeout:     0,
		IgnoreValidationErrors:  false,
	}

	mangagmentInited, err = c.Init(ctx)

	// Save the generated YAML to a file
	fileName := fmt.Sprintf("%s-cluster.yaml", clusterName)
	if err := os.WriteFile(fileName, yaml, 0644); err != nil {
		return fmt.Errorf("failed to write cluster configuration to file: %w", err)
	}

	fmt.Printf("Cluster configuration saved to %s\n", fileName)
	return nil
}
