package deploy

import (
	"fmt"
)

// BearerTokenRetriever represents any struct that can return a Bearer token.
// This supports both long-lived tokens (say from an environment variable) or
// short-lived tokens that need to be refreshed regularly.
type BearerTokenRetriever interface {
	RetrieveToken() string
}

// PodListRetriever represents any struct that returns a PodList
type PodListRetriever interface {
	PodInformation() (*PodList, error)
}

// Deployer represents any struct that can deploy a container
type Deployer interface {
	Deploy(containerTag string) error
}

// KubernetesClusterNamespace is a struct used to connect to a Kubernetes cluster.
type KubernetesClusterNamespace struct {
	Description  string
	PodRetriever PodListRetriever
	DeployMaker  Deployer
}

// GetPodList retrieves all the pods running in a deployment
func (n *KubernetesClusterNamespace) GetPodList() (*PodList, error) {
	if n.PodRetriever == nil {
		return nil, fmt.Errorf("missing PodListRetriever")
	}
	return n.PodRetriever.PodInformation()
}

// Deploy changes the image for an existing deployment and Kubernetes rebuilds the pods
func (n *KubernetesClusterNamespace) Deploy(containerTag string) error {
	if n.DeployMaker == nil {
		return fmt.Errorf("missing DeployMaker")
	}
	return n.DeployMaker.Deploy(containerTag)
}
