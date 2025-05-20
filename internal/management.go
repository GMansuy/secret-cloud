package internal

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/Tomy2e/cluster-api-provider-scaleway/internal/scope"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
)

type ClusterService struct {
	ManagementKubeconfigPath string
	TemplateOptions          client.GetClusterTemplateOptions
}

func NewClusterService(kubecfg string, t client.GetClusterTemplateOptions) *ClusterService {
	return &ClusterService{
		ManagementKubeconfigPath: kubecfg,
		TemplateOptions:          t,
	}
}

func (c *ClusterService) CreateCluster(ctx context.Context, clusterName string, opts client.GetClusterTemplateOptions) error {
	capiClient, err := client.New(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to create clusterctl client: %w", err)
	}

	template, err := capiClient.GetClusterTemplate(ctx, opts)
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

func (c *ClusterService) DeleteCluster(ctx context.Context, clusterName string) error {
	cmd := exec.Command("kubectl", "delete", "cluster", clusterName)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "Deleting Cluster", slog.String("clusterName", clusterName))
	return nil
}

func (c *ClusterService) SetCreationTemplateOptions(cluster *scope.Cluster) client.GetClusterTemplateOptions {
	template := c.TemplateOptions
	template.Kubeconfig = client.Kubeconfig{Path: c.ManagementKubeconfigPath}
	template.ClusterName = cluster.Name
	template.KubernetesVersion = "1.32.2"
	template.ControlPlaneMachineCount = &cluster.ControlplaneMachineCount
	template.WorkerMachineCount = &cluster.WorkerMachineCount
	return template
}
