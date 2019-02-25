package deploy

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPodListWhenMissingParameters(t *testing.T) {
	clusterNamespace := &KubernetesClusterNamespace{}
	podList, _ := clusterNamespace.GetPodList()
	assert.Nil(t, podList)
}

func TestPodListOverview(t *testing.T) {
	var timestamp time.Time
	clusterNamespace := &KubernetesClusterNamespace{
		PodRetriever: &MockPodList{},
	}

	podList, _ := clusterNamespace.GetPodList()
	overview := podList.Overview()

	pod := overview[0]
	assert.Equal(t, "myapp-deployment-1376141578-9q2hx", pod.Name)
	assert.Equal(t, "Running", pod.Status)
	assert.Equal(t, "40716241027b9639db1f1067d5ea3b25087dd12e", pod.Tag)
	timestamp, _ = time.Parse(time.RFC3339, "2017-03-30T16:39:21Z")
	assert.Equal(t, timestamp, pod.Created)

	pod = overview[1]
	assert.Equal(t, "myapp-deployment-1376141578-c2sw6", pod.Name)
	assert.Equal(t, "Running", pod.Status)
	assert.Equal(t, "40716241027b9639db1f1067d5ea3b25087dd12e", pod.Tag)
	timestamp, _ = time.Parse(time.RFC3339, "2017-03-30T16:36:44Z")
	assert.Equal(t, timestamp, pod.Created)

	pod = overview[2]
	assert.Equal(t, "myapp-deployment-1376141578-s2jfg", pod.Name)
	assert.Equal(t, "Running", pod.Status)
	assert.Equal(t, "40716241027b9639db1f1067d5ea3b25087dd12e", pod.Tag)
	timestamp, _ = time.Parse(time.RFC3339, "2017-03-30T16:39:07Z")
	assert.Equal(t, timestamp, pod.Created)

	pod = overview[3]
	assert.Equal(t, "myapp-deployment-1376141578-wv7wh", pod.Name)
	assert.Equal(t, "Running", pod.Status)
	assert.Equal(t, "40716241027b9639db1f1067d5ea3b25087dd12e", pod.Tag)
	timestamp, _ = time.Parse(time.RFC3339, "2017-03-30T16:36:45Z")
	assert.Equal(t, timestamp, pod.Created)
}

func TestPodListNoContainerStatus(t *testing.T) {
	var timestamp time.Time
	clusterNamespace := &KubernetesClusterNamespace{
		PodRetriever: &MockPodListWithoutContainerStatuses{},
	}

	podList, _ := clusterNamespace.GetPodList()
	overview := podList.Overview()

	pod := overview[0]
	assert.Equal(t, "myapp-deployment-1376141578-9q2hx", pod.Name)
	assert.Equal(t, "Failed", pod.Status)
	assert.Equal(t, "", pod.Tag)
	timestamp, _ = time.Parse(time.RFC3339, "2017-03-30T16:39:21Z")
	assert.Equal(t, timestamp, pod.Created)
}

func TestPodListFilterByDeployment_WithMatch(t *testing.T) {
	clusterNamespace := &KubernetesClusterNamespace{
		PodRetriever: &MockPodListWithoutContainerStatuses{},
	}

	podList, _ := clusterNamespace.GetPodList()

	filteredPodList := podList.FilterByDeployment("myapp-deployment")
	pod := filteredPodList.Items[0]
	assert.Equal(t, "myapp-deployment-1376141578-9q2hx", pod.Metadata.Name)
}

func TestPodListFilterByDeployment_NoMatch(t *testing.T) {
	clusterNamespace := &KubernetesClusterNamespace{
		PodRetriever: &MockPodListWithoutContainerStatuses{},
	}

	podList, _ := clusterNamespace.GetPodList()

	filteredPodList := podList.FilterByDeployment("unknown-deployment")
	assert.Equal(t, 0, len(filteredPodList.Items))
}

//
// MOCK DATA
//

type MockPodList struct{}

func (p *MockPodList) PodInformation() (*PodList, error) {
	podList := &PodList{}
	err := json.Unmarshal([]byte(examplePodList), podList)
	return podList, err
}

type MockPodListWithoutContainerStatuses struct{}

func (p *MockPodListWithoutContainerStatuses) PodInformation() (*PodList, error) {
	podList := &PodList{}
	err := json.Unmarshal([]byte(examplePodListWithoutContainerStatuses), podList)
	return podList, err
}

const examplePodList = `
{
	"kind": "PodList",
	"apiVersion": "v1",
	"items": [
	{
		"metadata": {
		"name": "myapp-deployment-1376141578-9q2hx",
		"generateName": "myapp-deployment-1376141578-",
		"namespace": "myapp-development",
		"creationTimestamp": "2017-03-30T16:39:21Z",
		"ownerReferences": []
		},
		"spec": {
		"volumes": [],
		"containers": [
			{
			"name": "myapp-container",
			"image": "artifactory.myorg.com:5010/myapp-docker-image:40716241027b9639db1f1067d5ea3b25087dd12e",
			"env": [],
			"volumeMounts": [],
			"terminationMessagePath": "/dev/termination-log",
			"imagePullPolicy": "Always"
			}
		],
		"restartPolicy": "Always",
		"terminationGracePeriodSeconds": 30,
		"dnsPolicy": "ClusterFirst",
		"serviceAccountName": "default",
		"serviceAccount": "default",
		"nodeName": "ip-1-2-3-4.internal",
		"securityContext": {},
		"imagePullSecrets": [
			{
			"name": "artifactory"
			}
		]
		},
		"status": {
		"phase": "Running",
		"conditions": [
			{
			"type": "Initialized",
			"status": "True",
			"lastProbeTime": null,
			"lastTransitionTime": "2017-03-30T16:39:22Z"
			},
			{
			"type": "Ready",
			"status": "True",
			"lastProbeTime": null,
			"lastTransitionTime": "2017-03-30T16:39:25Z"
			},
			{
			"type": "PodScheduled",
			"status": "True",
			"lastProbeTime": null,
			"lastTransitionTime": "2017-03-30T16:39:21Z"
			}
		],
		"hostIP": "2.3.4.5",
		"podIP": "4.5.6.7",
		"startTime": "2017-03-30T16:39:22Z",
		"containerStatuses": [
			{
			"name": "myapp-container",
			"state": {
				"running": {
				"startedAt": "2017-03-30T16:39:25Z"
				}
			},
			"lastState": {},
			"ready": true,
			"restartCount": 0,
			"image": "artifactory.myorg.com:5010/myapp-docker-image:40716241027b9639db1f1067d5ea3b25087dd12e",
			"imageID": "docker-pullable://artifactory.myorg.com:5010/myapp-docker-image@sha256:056daaa91dbe40676b2a65095ba603d2f1789774b0963546aaa244a79c397683",
			"containerID": "docker://f6d31b15f72d1323718daacb3b4bec2a895c1471d3931b78ade4572f00ef9076"
			}
		]
		}
	},
	{
		"metadata": {
		"name": "myapp-deployment-1376141578-c2sw6",
		"generateName": "myapp-deployment-1376141578-",
		"namespace": "myapp-development",
		"creationTimestamp": "2017-03-30T16:36:44Z",
		"ownerReferences": []
		},
		"spec": {
		"volumes": [],
		"containers": [
			{
			"name": "myapp-container",
			"image": "artifactory.myorg.com:5010/myapp-docker-image:40716241027b9639db1f1067d5ea3b25087dd12e",
			"env": [],
			"volumeMounts": [],
			"terminationMessagePath": "/dev/termination-log",
			"imagePullPolicy": "Always"
			}
		],
		"restartPolicy": "Always",
		"terminationGracePeriodSeconds": 30,
		"dnsPolicy": "ClusterFirst",
		"serviceAccountName": "default",
		"serviceAccount": "default",
		"nodeName": "ip-1-2-3-4.internal",
		"securityContext": {},
		"imagePullSecrets": [
			{
			"name": "artifactory"
			}
		]
		},
		"status": {
		"phase": "Running",
		"conditions": [
			{
			"type": "Initialized",
			"status": "True",
			"lastProbeTime": null,
			"lastTransitionTime": "2017-03-30T16:36:44Z"
			},
			{
			"type": "Ready",
			"status": "True",
			"lastProbeTime": null,
			"lastTransitionTime": "2017-03-30T16:39:21Z"
			},
			{
			"type": "PodScheduled",
			"status": "True",
			"lastProbeTime": null,
			"lastTransitionTime": "2017-03-30T16:36:44Z"
			}
		],
		"hostIP": "2.3.4.5",
		"podIP": "4.5.6.7",
		"startTime": "2017-03-30T16:36:44Z",
		"containerStatuses": [
			{
			"name": "myapp-container",
			"state": {
				"waiting": {
				"reason": "ImagePullBackOff",
				"message": "Back-off pulling image \"artifactory.myorg.com:5010/myapp-docker-image:3362ff29b425bb9fdc1d97039a2d9a2fa2d7d454\""
				}
			},
			"lastState": {},
			"ready": true,
			"restartCount": 0,
			"image": "artifactory.myorg.com:5010/myapp-docker-image:40716241027b9639db1f1067d5ea3b25087dd12e",
			"imageID": "docker-pullable://artifactory.myorg.com:5010/myapp-docker-image@sha256:056daaa91dbe40676b2a65095ba603d2f1789774b0963546aaa244a79c397683",
			"containerID": "docker://20942cf97db86f546e35a96a0e8acf585ef8bc8f6e8ca94426449d7adae93159"
			}
		]
		}
	},
	{
		"metadata": {
		"name": "myapp-deployment-1376141578-s2jfg",
		"generateName": "myapp-deployment-1376141578-",
		"namespace": "myapp-development",
		"creationTimestamp": "2017-03-30T16:39:07Z",
		"ownerReferences": []
		},
		"spec": {
		"volumes": [],
		"containers": [
			{
			"name": "myapp-container",
			"image": "artifactory.myorg.com:5010/myapp-docker-image:40716241027b9639db1f1067d5ea3b25087dd12e",
			"env": [],
			"volumeMounts": [],
			"terminationMessagePath": "/dev/termination-log",
			"imagePullPolicy": "Always"
			}
		],
		"restartPolicy": "Always",
		"terminationGracePeriodSeconds": 30,
		"dnsPolicy": "ClusterFirst",
		"serviceAccountName": "default",
		"serviceAccount": "default",
		"nodeName": "ip-1-2-3-4.internal",
		"securityContext": {},
		"imagePullSecrets": [
			{
			"name": "artifactory"
			}
		]
		},
		"status": {
		"phase": "Running",
		"conditions": [
			{
			"type": "Initialized",
			"status": "True",
			"lastProbeTime": null,
			"lastTransitionTime": "2017-03-30T16:39:08Z"
			},
			{
			"type": "Ready",
			"status": "True",
			"lastProbeTime": null,
			"lastTransitionTime": "2017-03-30T16:41:15Z"
			},
			{
			"type": "PodScheduled",
			"status": "True",
			"lastProbeTime": null,
			"lastTransitionTime": "2017-03-30T16:39:07Z"
			}
		],
		"hostIP": "2.3.4.5",
		"podIP": "4.5.6.7",
		"startTime": "2017-03-30T16:39:08Z",
		"containerStatuses": [
			{
			"name": "myapp-container",
			"state": {
				"running": {
				"startedAt": "2017-03-30T16:41:14Z"
				}
			},
			"lastState": {},
			"ready": true,
			"restartCount": 0,
			"image": "artifactory.myorg.com:5010/myapp-docker-image:40716241027b9639db1f1067d5ea3b25087dd12e",
			"imageID": "docker-pullable://artifactory.myorg.com:5010/myapp-docker-image@sha256:056daaa91dbe40676b2a65095ba603d2f1789774b0963546aaa244a79c397683",
			"containerID": "docker://4f821bd4f3ef4b42233b7d6465737a4581a4bfdd8b97cd0408cd850b3fba7f0c"
			}
		]
		}
	},
	{
		"metadata": {
		"name": "myapp-deployment-1376141578-wv7wh",
		"generateName": "myapp-deployment-1376141578-",
		"namespace": "myapp-development",
		"creationTimestamp": "2017-03-30T16:36:45Z",
		"ownerReferences": []
		},
		"spec": {
		"volumes": [],
		"containers": [
			{
			"name": "myapp-container",
			"image": "artifactory.myorg.com:5010/myapp-docker-image:40716241027b9639db1f1067d5ea3b25087dd12e",
			"env": [],
			"volumeMounts": [],
			"terminationMessagePath": "/dev/termination-log",
			"imagePullPolicy": "Always"
			}
		],
		"restartPolicy": "Always",
		"terminationGracePeriodSeconds": 30,
		"dnsPolicy": "ClusterFirst",
		"serviceAccountName": "default",
		"serviceAccount": "default",
		"nodeName": "ip-1-2-3-4.internal",
		"securityContext": {},
		"imagePullSecrets": [
			{
			"name": "artifactory"
			}
		]
		},
		"status": {
		"phase": "Running",
		"conditions": [
			{
			"type": "Initialized",
			"status": "True",
			"lastProbeTime": null,
			"lastTransitionTime": "2017-03-30T16:36:45Z"
			},
			{
			"type": "Ready",
			"status": "True",
			"lastProbeTime": null,
			"lastTransitionTime": "2017-03-30T16:39:07Z"
			},
			{
			"type": "PodScheduled",
			"status": "True",
			"lastProbeTime": null,
			"lastTransitionTime": "2017-03-30T16:36:45Z"
			}
		],
		"hostIP": "2.3.4.5",
		"podIP": "4.5.6.7",
		"startTime": "2017-03-30T16:36:45Z",
		"containerStatuses": [
			{
			"name": "myapp-container",
			"state": {
				"running": {
				"startedAt": "2017-03-30T16:39:07Z"
				}
			},
			"lastState": {},
			"ready": true,
			"restartCount": 0,
			"image": "artifactory.myorg.com:5010/myapp-docker-image:40716241027b9639db1f1067d5ea3b25087dd12e",
			"imageID": "docker-pullable://artifactory.myorg.com:5010/myapp-docker-image@sha256:056daaa91dbe40676b2a65095ba603d2f1789774b0963546aaa244a79c397683",
			"containerID": "docker://17b73965115e4220a36348b8fad3243e77b9599a4abc0417dccc12bd59e01b38"
			}
		]
		}
	}
	]
}
`

const examplePodListWithoutContainerStatuses = `
{
	"kind": "PodList",
	"apiVersion": "v1",
	"items": [
	{
		"metadata": {
		"name": "myapp-deployment-1376141578-9q2hx",
		"generateName": "myapp-deployment-1376141578-",
		"namespace": "myapp-development",
		"creationTimestamp": "2017-03-30T16:39:21Z",
		"ownerReferences": []
		},
		"spec": {
		"volumes": [],
		"containers": [
			{
			"name": "myapp-container",
			"image": "artifactory.myorg.com:5010/myapp-docker-image:40716241027b9639db1f1067d5ea3b25087dd12e",
			"env": [],
			"volumeMounts": [],
			"terminationMessagePath": "/dev/termination-log",
			"imagePullPolicy": "Always"
			}
		],
		"restartPolicy": "Always",
		"terminationGracePeriodSeconds": 30,
		"dnsPolicy": "ClusterFirst",
		"serviceAccountName": "default",
		"serviceAccount": "default",
		"nodeName": "ip-1-2-3-4.internal",
		"securityContext": {},
		"imagePullSecrets": [
			{
			"name": "artifactory"
			}
		]
		},
		"status": {
		"phase": "Failed",
		"reason": "Evicted",
		"startTime": "2017-03-30T16:39:22Z"
		}
	}
	]
}
`
