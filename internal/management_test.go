package internal_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/Tomy2e/cluster-api-provider-scaleway/internal"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
)

func TestCreateClusterE2E(t *testing.T) {
	//given
	workerCount := int64(1)
	controlPlaneCount := int64(1)
	clusterName := "test-cluster"
	fakeCluster := internal.ClusterService{
		ManagementKubeconfigPath: "kubeconfig-test",
		TemplateOptions:          client.GetClusterTemplateOptions{},
	}

	urlSource := &client.URLSourceOptions{
		URL: "https://github.com/Tomy2e/cluster-api-provider-scaleway/releases/download/v0.0.3/cluster-template.yaml",
	}
	ClusterCreationConxtext := client.GetClusterTemplateOptions{
		Kubeconfig: client.Kubeconfig{
			Path: fakeCluster.ManagementKubeconfigPath,
		},
		URLSource:                urlSource,
		ClusterName:              clusterName,
		KubernetesVersion:        "1.32.2",
		ControlPlaneMachineCount: &workerCount,
		WorkerMachineCount:       &controlPlaneCount,
		TargetNamespace:          "default",
	}

	err := fakeCluster.CreateCluster(context.TODO(), clusterName, ClusterCreationConxtext)
	require.NoError(t, err)

	capiClient, err := client.New(context.Background(), "")
	require.NoError(t, err)

	kubeConfingOptions := client.GetKubeconfigOptions{
		WorkloadClusterName: clusterName,
		Kubeconfig: client.Kubeconfig{
			Path: "kubeconfig-test",
		},
		Namespace: "default",
	}

	//when
	var retries int
	for {
		err := ClusterReadinessProbe(clusterName, capiClient, kubeConfingOptions)
		if err != nil {
			fmt.Printf("Cluster is not ready : %s\n", err)
		} else {
			fmt.Printf("Cluster %s is ready after %d attempts\n", clusterName, retries)
			break
		}
		if retries >= 100 {
			t.Fatalf("Cluster %s is not ready after 10 attempts %s", clusterName, err)
			break
		}
		retries++
		time.Sleep(time.Second * 10)
	}
}

func ClusterReadinessProbe(clusterName string, capiClient client.Client, kubeConfigOptions client.GetKubeconfigOptions) error {
	//assert that a cluster has been created
	kubeconfig, err := capiClient.GetKubeconfig(context.Background(), kubeConfigOptions)
	if err != nil {
		return fmt.Errorf("failed to get kubeconfig: %w", err)
	}
	err = os.WriteFile("output-kubeconfig.yaml", []byte(kubeconfig), 0644)
	if err != nil {
		return fmt.Errorf("failed to write kubeconfig: %w", err)
	}
	cmd := exec.Command("kubectl", "get", "no", "--kubeconfig=output-kubeconfig.yaml")
	_, err = cmd.CombinedOutput()
	if err == nil {
		fmt.Printf("Cluster %s is ready\n", clusterName)
		return nil
	} else {
		fmt.Printf("Cluster %s is not ready yet, retrying...\n", clusterName)
		return err
	}
}
