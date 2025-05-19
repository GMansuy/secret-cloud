package internal

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"os"

	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	cr "sigs.k8s.io/controller-runtime/pkg/client"
	crcfg "sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	scheme = runtime.NewScheme()
	c      client.Client
)

func GenerateClusterConfigFromMemory(ctx context.Context, clusterName string, controlPlaneCount, workerCount int64, kubeConfigPath string) error {

	capiClient, err := client.New(ctx, kubeConfigPath)
	if err != nil {
		return fmt.Errorf("failed to create clusterctl client: %w", err)
	}

	// Generate the cluster configuration
	templateOptions := client.GetClusterTemplateOptions{
		Kubeconfig: client.Kubeconfig{
			Path: kubeConfigPath,
		},
		URLSource: &client.URLSourceOptions{
			URL: "https://github.com/Tomy2e/cluster-api-provider-scaleway/releases/download/v0.0.3/infrastructure-components.yaml",
		},
		ClusterName:              clusterName,
		TargetNamespace:          "default",
		KubernetesVersion:        "1.32.2",
		ControlPlaneMachineCount: &controlPlaneCount,
		WorkerMachineCount:       &workerCount,
	}
	// Get the cluster template

	template, err := capiClient.GetClusterTemplate(ctx, templateOptions)
	if err != nil {
		return fmt.Errorf("failed to get cluster template: %w", err)
	}

	restConfig, err := crcfg.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get rest config: %w", err)
	}

	ControllerRuntimeClient, err := cr.New(restConfig, cr.Options{Scheme: scheme})

	ControllerRuntimeClient.Create(ctx, template)

	yamlData, err := template.Yaml()
	if err != nil {
		return err
	}
	// Save the generated YAML to a file
	fileName := fmt.Sprintf("%s-cluster.yaml", clusterName)
	if err := os.WriteFile(fileName, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write cluster configuration to file: %w", err)
	}

	fmt.Printf("Cluster configuration saved to %s\n", fileName)
	return nil
}
