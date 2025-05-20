package internal

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
)

func CreateCluster(ctx context.Context, clusterName string, controlPlaneCount, workerCount int64, kubeConfigPath string) error {
	capiClient, err := client.New(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to create clusterctl client: %w", err)
	}

	templateOptions := client.GetClusterTemplateOptions{
		Kubeconfig: client.Kubeconfig{
			Path: kubeConfigPath,
		},
		URLSource: &client.URLSourceOptions{
			URL: "https://github.com/Tomy2e/cluster-api-provider-scaleway/releases/download/v0.0.3/cluster-template.yaml",
		},
		ClusterName:              clusterName,
		TargetNamespace:          "default",
		KubernetesVersion:        "1.32.2",
		ControlPlaneMachineCount: &controlPlaneCount,
		WorkerMachineCount:       &workerCount,
	}

	template, err := capiClient.GetClusterTemplate(ctx, templateOptions)
	if err != nil {
		return fmt.Errorf("failed to get cluster template: %w", err)
	}

	yamlData, err := template.Yaml()
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s-cluster.yaml", clusterName)
	if err := os.WriteFile(fileName, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write cluster configuration to file: %w", err)
	}

	cmd := exec.Command("kubectl", "apply", "-f", fmt.Sprintf("%s-cluster.yaml", clusterName))
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "Creating Cluster", slog.String("clusterName", clusterName))
	return nil
}

func DeleteCluster(ctx context.Context, clusterName string, kubeConfigPath string) error {
	cmd := exec.Command("kubectl", "delete", "cluster", clusterName)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "Deleting Cluster", slog.String("clusterName", clusterName))
	return nil
}
