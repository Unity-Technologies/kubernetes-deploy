package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Unity-Technologies/kubernetes-deploy/deploy"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	client := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	tokenProvider := new(sampleBearerTokenProvider)

	cluster := deploy.KubernetesClusterNamespace{
		Description: os.Getenv("DESCRIPTION"),
		PodRetriever: &deploy.KubernetesPodListRetriever{
			Client:             client,
			Endpoint:           os.Getenv("KUBERNETES_ENDPOINT"),
			Namespace:          os.Getenv("KUBERNETES_NAMESPACE"),
			BearerTokenService: tokenProvider,
		},
		DeployMaker: &deploy.KubernetesDeployer{
			Client:             client,
			Endpoint:           os.Getenv("KUBERNETES_ENDPOINT"),
			Namespace:          os.Getenv("KUBERNETES_NAMESPACE"),
			BearerTokenService: tokenProvider,
			DeploymentName:     os.Getenv("KUBERNETES_DEPLOYMENT_NAME"),
			ContainerName:      os.Getenv("KUBERNETES_DEPLOYMENT_CONTAINERNAME"),
			ContainerImage:     os.Getenv("KUBERNETES_DEPLOYMENT_IMAGE_PREFIX"),
		},
	}

	command, containerTag := pickCommand(os.Args)

	if command == "ls" {
		// Get a list of pods and their status
		podList, err := cluster.GetPodList()
		if err != nil {
			fmt.Printf("Unable to retrieve pod list due to %s", err.Error())
			return
		}

		printStatus(podList, "")
	}

	if command == "deploy" {
		// Deploy container named KUBERNETES_DEPLOYMENT_IMAGE_PREFIX:tag
		err := cluster.Deploy(containerTag)
		if err != nil {
			fmt.Printf("Unable to deploy %q due to %s", containerTag, err.Error())
			return
		}

		podList, _ := cluster.GetPodList()
		printStatus(podList, containerTag)
	}
}

// Sample BearerTokenRetriever
type sampleBearerTokenProvider struct{}

func (*sampleBearerTokenProvider) RetrieveToken() string {
	return os.Getenv("KUBERNETES_ENDPOINT_BEARER_TOKEN")
}

//
// Helpers
//

// printStatus runs through all pods in a deployment, and displays their status.
func printStatus(podList *deploy.PodList, desiredImageTag string) {
	now := time.Now()

	overview := podList.Overview()
	for _, overview := range overview {
		fmt.Printf("`%s` image has been *%s* for %.1f hours.\n",
			overview.Tag,
			overview.Status,
			now.Sub(overview.Created).Hours())
	}
}

// pickCommand supports `deploy <hash>` otherwise defaults to `ls`
func pickCommand(osArgs []string) (string, string) {
	args := osArgs[1:]

	command := "ls"
	tag := ""

	if len(args) == 2 && args[0] == "deploy" {
		command = args[0]
		tag = args[1]
	}

	return command, tag
}
