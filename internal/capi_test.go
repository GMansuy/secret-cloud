package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/config"
)

func TestProviderAddedCorrectly(t *testing.T) {
	// Create a memory reader
	memoryReader := config.NewMemoryReader()

	// Add the provider
	providerName := "scaleway"
	providerType := clusterctlv1.InfrastructureProviderType
	providerURL := "https://github.com/Tomy2e/cluster-api-provider-scaleway/releases/latest/infrastructure-components.yaml"

	provider, err := memoryReader.AddProvider(providerName, providerType, providerURL)

	// Verify provider was added successfully
	assert.NoError(t, err)
	assert.NotNil(t, provider)

	// Verify we can retrieve the provider from the memory reader
	providers, err := memoryReader.Get("providers")
	assert.NoError(t, err)
	assert.Equal(t, providerName, providers)

}
