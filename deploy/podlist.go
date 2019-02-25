package deploy

import (
	"strings"
	"time"
)

// PodList holds a Kubernetes PodList.
type PodList struct {
	Items []PodMetadataContainer `json:"items"`
}

// PodItem is an app specific summary status coming back from pods
type PodItem struct {
	Name    string
	Status  string
	Created time.Time
	Tag     string
}

// Overview for a group of pods in a deployment
func (p *PodList) Overview() []PodItem {
	var metadata []PodItem

	for _, item := range p.Items {
		tag := ""
		// an Evicted pod for example, will not have any ContainerStatuses
		if item.Status.ContainerStatuses != nil {
			tag = formatPodImage(item.Status.ContainerStatuses[0].Image)
		}

		metadata = append(metadata, PodItem{
			Name:    item.Metadata.Name,
			Status:  item.Status.Phase,
			Created: item.Metadata.CreationTimestamp,
			Tag:     tag,
		})
	}
	return metadata
}

// FilterByDeployment returns a new PodList with only deployment names that match prefix.
// Handy for only retrieving specific deployment when running multiple deployments in same namespace.
func (p *PodList) FilterByDeployment(namePrefix string) *PodList {
	pods := []PodMetadataContainer{}

	for _, item := range p.Items {
		if strings.HasPrefix(item.Metadata.Name, namePrefix) {
			pods = append(pods, item)
		}
	}
	return &PodList{
		Items: pods,
	}
}

// formatPodImage converts a full pod image name into only the docker image and commit hash
func formatPodImage(raw string) (result string) {
	s := strings.Split(raw, ":")
	if len(s) == 3 {
		result = s[2]
	}
	return
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

// PodContainerStatusesStateRunning is part of PodList response coming back from Kubernetes
type PodContainerStatusesStateRunning struct {
	StartedAt time.Time `json:"startedAt"`
}

// PodContainerStatusesStateWaiting is part of PodList response coming back from Kubernetes
type PodContainerStatusesStateWaiting struct {
	Reason string `json:"reason"`
}

// PodContainerStatusesState is part of PodList response coming back from Kubernetes
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
