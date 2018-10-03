package deploy

import (
	"fmt"
	"strings"
	"time"
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

// PodStatusForFirstContainer returns a state for an individual Pod. Leave desiredImageHash as an
// empty string to avoid warning check if desired container is not running.
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

// FormatPodStatusForFirstContainer returns a nicely formatting Pod status string including a warning
// if image doesn't yet match one that is being deployed. Formatted in markdown for Slack.
//
// Examples:
// "• `77d0ea51fdc30234918f2726d26479c66b7f7777` image has been *Running* for 45.7 hours."
// "• `77d0ea51fdc30234918f2726d26479c66b7f7777` image has been *WrongHash* for 45.8 hours."
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
