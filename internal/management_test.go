package internal_test

import (
	"context"
	"fmt"
	"github.com/Tomy2e/cluster-api-provider-scaleway/internal"
	"os"
	"os/exec"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	"testing"
)

func TestCreateClusterE2E(t *testing.T) {
	//given
	workerCount := int64(1)
	controlPlaneCount := int64(1)
	clusterName := "test-cluster"
	urlSource := &client.URLSourceOptions{
		URL: "https://github.com/Tomy2e/cluster-api-provider-scaleway/releases/download/v0.0.3/cluster-template.yaml",
	}
	ClusterCreationConxtext := client.GetClusterTemplateOptions{
		Kubeconfig: client.Kubeconfig{
			Path: "var/kubeconfig-le-cluster-r-et-d.yaml",
		},
		URLSource:                urlSource,
		ClusterName:              clusterName,
		KubernetesVersion:        "1.32.2",
		ControlPlaneMachineCount: &workerCount,
		WorkerMachineCount:       &controlPlaneCount,
		TargetNamespace:          "default",
	}

	//when
	for i := 0; i < 10; i++ {
		err := ClusterReadinessProbe(clusterName, ClusterCreationConxtext)
		if err != nil {
			fmt.Println("Cluster is not ready")
		} else {
			fmt.Printf("Cluster %s is ready after %d attempts\n", clusterName, i+1)
			break
		}

		if i == 9 {
			t.Fatalf("Cluster %s is not ready after 10 attempts %s", clusterName, err)
			break
		}
	}

	//then

	//assert.NoError(t, err, "no nodes found in the cluster, please check if the cluster is created successfully")
}

func ClusterReadinessProbe(clusterName string, ClusterCreationConxtext client.GetClusterTemplateOptions) error {
	fakeCluster := internal.ClusterService{
		ManagementKubeconfigPath: "var/kubeconfig-le-cluster-r-et-d.yaml",
		TemplateOptions:          client.GetClusterTemplateOptions{},
	}
	err := fakeCluster.CreateCluster(context.TODO(), clusterName, ClusterCreationConxtext)
	if err != nil {
		return fmt.Errorf("failed to create cluster: %w", err)
	}

	//assert that a cluster has been created
	capiClient, err := client.New(context.Background(), "")
	if err != nil {
		return fmt.Errorf("failed to create clusterctl client: %w", err)
	}
	kubeConfingOptions := client.GetKubeconfigOptions{
		WorkloadClusterName: clusterName,
		Kubeconfig: client.Kubeconfig{
			Path: "/Users/michel.amoussou/perso/work-and-workrelated/secret-cloud/var/kubeconfig-le-cluster-r-et-d.yaml",
		},
	}
	kubeconfig, err := capiClient.GetKubeconfig(context.Background(), kubeConfingOptions)
	if err != nil {
		return fmt.Errorf("failed to get kubeconfig: %w", err)

	}
	os.WriteFile(fmt.Sprintf("%s-kubeconfig.yaml", kubeconfig), []byte(kubeconfig), 0644)

	cmd := exec.Command("kubectl", "get", "no", "--kubeconfig=", fmt.Sprintf("%s-kubeconfig.yaml", kubeconfig))
	_, err = cmd.CombinedOutput()
	if err == nil {
		fmt.Printf("Cluster %s is ready\n", clusterName)
		return nil
	} else {
		fmt.Printf("Cluster %s is not ready yet, retrying...\n", clusterName)
		return err
	}
}
