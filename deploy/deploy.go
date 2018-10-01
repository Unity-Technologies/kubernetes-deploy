package deploy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type BearerTokenRetriever interface {
	RetrieveToken() string
}

// KubernetesCluster is a struct used to connect to a Kubernetes cluster.
type KubernetesClusterNamespace struct {
	Client             *http.Client
	Description        string
	Endpoint           string
	Namespace          string
	DeploymentName     string
	ContainerName      string
	ContainerImage     string
	BearerTokenService BearerTokenRetriever
}

func (n *KubernetesClusterNamespace) GetPodList() (*PodList, error) {
	url := fmt.Sprintf("https://%s/api/v1/namespaces/%s/pods", n.Endpoint, n.Namespace)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", n.BearerTokenService.RetrieveToken()))

	res, err := n.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("received %v", res.StatusCode)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	podList, err := convertJSON(body)
	if err != nil {
		return nil, err
	}

	return podList, nil
}

func (n *KubernetesClusterNamespace) Deploy(containerTag string) error {
	url := fmt.Sprintf("https://%s/apis/extensions/v1beta1/namespaces/%s/deployments/%s", n.Endpoint, n.Namespace, n.DeploymentName)
	image := fmt.Sprintf(`{"spec":{"template":{"spec":{"containers":[{"name":"%s","image":"%s:%s"}]}}}}`, n.ContainerName, n.ContainerImage, containerTag)
	payload := bytes.NewBuffer([]byte(image))

	req, err := http.NewRequest(http.MethodPatch, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/strategic-merge-patch+json")

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", n.BearerTokenService.RetrieveToken()))
	res, err := n.Client.Do(req)
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

	log.Printf("--- DEPLOYED %s:%s ---", n.ContainerImage, containerTag)
	// log.Printf("podDeployResponse.AvailableReplicas %+v", podDeployResponse.Status.AvailableReplicas)

	return nil
}

func convertJSON(jsonBytes []byte) (podList *PodList, err error) {
	podList = &PodList{}
	err = json.Unmarshal([]byte(jsonBytes), podList)
	return
}

// PodStatusForFirstContainer returns a state for an individual Pod
func PodStatusForFirstContainer(pod *PodMetadataContainer, desiredImageHash string) (status string) {

	// only check hash if desiredImageHash is not an empty string
	// so that status checks (which don't know hash) still work
	if len(desiredImageHash) > 0 && desiredImageHash != formatPodImage(pod.Status.ContainerStatuses[0].Image) {
		status = "WrongHash"
		return
	}

	running := pod.Status.ContainerStatuses[0].State.Running
	waiting := pod.Status.ContainerStatuses[0].State.Waiting

	status = "Unknown state"

	if running != nil {
		status = "Running"
	}

	if waiting != nil {
		status = pod.Status.ContainerStatuses[0].State.Waiting.Reason
	}
	return
}

// FormatPodStatusForFirstContainer returns a nicely formatting Pod status string.
func FormatPodStatusForFirstContainer(pod *PodMetadataContainer, now time.Time, desiredImageHash string) (result string) {
	if len(pod.Status.ContainerStatuses) > 0 {
		result = fmt.Sprintf("\n• `%s` image has been *%s* for %.1f hours.",
			formatPodImage(pod.Status.ContainerStatuses[0].Image),
			PodStatusForFirstContainer(pod, desiredImageHash),
			now.Sub(pod.Metadata.CreationTimestamp).Hours())
		return
	}
	result = fmt.Sprintf("\n• *No containers found for %q deployment.* Phase `%s`.", pod.Metadata.Name, pod.Status.Phase)
	return
}

// formatPodImage converts a full pod image name into only the docker image and commit hash
func formatPodImage(raw string) (result string) {
	s := strings.Split(raw, ":")
	if len(s) == 3 {
		result = s[2]
	}
	return
}

// PodList holds a Kubernetes PodList.
type PodList struct {
	Items []PodMetadataContainer `json:"items"`
}

// PodMetadataContainer houses PodMetadataDetail values.
type PodMetadataContainer struct {
	Metadata PodMetadataDetail `json:"metadata"`
	Status   PodStatus         `json:"status"`
}

// PodMetadataDetail has details about an individual Kubernetes Pod.
type PodMetadataDetail struct {
	Name              string    `json:"name"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}

type PodContainerStatusesStateRunning struct {
	StartedAt time.Time `json:"startedAt"`
}

type PodContainerStatusesStateWaiting struct {
	Reason string `json:"reason"`
}

type PodContainerStatusesState struct {
	Running *PodContainerStatusesStateRunning `json:"running,omitempty"`
	Waiting *PodContainerStatusesStateWaiting `json:"waiting,omitempty"`
}

//PodContainerStatuses used for determining what state the container is in.
type PodContainerStatuses struct {
	State PodContainerStatusesState `json:"state"`
	Image string                    `json:"image,omitempty"`
}

// PodStatus has details about what state the Pod is in.
type PodStatus struct {
	Phase             string                 `json:"phase"`
	ContainerStatuses []PodContainerStatuses `json:"containerStatuses"`
}

// FOR DEPLOY

type PodDeployResponse struct {
	Status PodDeployStatus `json:"status"`
}

// PodDeployStatus has details about what state the Pod is in.
type PodDeployStatus struct {
	AvailableReplicas int `json:"availableReplicas"`
}
