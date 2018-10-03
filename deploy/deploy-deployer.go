package deploy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// KubernetesDeployer updates a deployment via Kubernetes API
type KubernetesDeployer struct {
	Client             *http.Client
	Endpoint           string
	Namespace          string
	DeploymentName     string
	ContainerName      string
	ContainerImage     string
	BearerTokenService BearerTokenRetriever
}

// Deploy a container via Kubernetes API
func (d *KubernetesDeployer) Deploy(containerTag string) error {
	if d.Endpoint == "" || d.Namespace == "" {
		return fmt.Errorf("missing Endpoint or Namespace information")
	}

	url := fmt.Sprintf("https://%s/apis/extensions/v1beta1/namespaces/%s/deployments/%s", d.Endpoint, d.Namespace, d.DeploymentName)
	image := fmt.Sprintf(`{"spec":{"template":{"spec":{"containers":[{"name":"%s","image":"%s:%s"}]}}}}`, d.ContainerName, d.ContainerImage, containerTag)
	payload := bytes.NewBuffer([]byte(image))

	req, err := http.NewRequest(http.MethodPatch, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/strategic-merge-patch+json")

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", d.BearerTokenService.RetrieveToken()))
	res, err := d.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	podDeployResponse := &PodDeployResponse{}
	err = json.Unmarshal([]byte(body), &podDeployResponse)
	if err != nil {
		return err
	}
	return nil
}

// PodDeployResponse is part of the deployment response coming back from Kubernetes
type PodDeployResponse struct {
	Status PodDeployStatus `json:"status"`
}

// PodDeployStatus has details about what state the Pod is in.
type PodDeployStatus struct {
	AvailableReplicas int `json:"availableReplicas"`
}
