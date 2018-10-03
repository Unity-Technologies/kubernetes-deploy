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
		metadata = append(metadata, PodItem{
			Name:    item.Metadata.Name,
			Status:  item.Status.Phase,
			Created: item.Metadata.CreationTimestamp,
			Tag:     formatPodImage(item.Status.ContainerStatuses[0].Image),
		})
	}
	return metadata
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
