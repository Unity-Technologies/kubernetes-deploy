package main

import (
	"crypto/tls"
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

	tokenProvider := new(bearerTokenProvider)

	cluster := deploy.KubernetesClusterNamespace{
		Client:             client,
		Description:        os.Getenv("DESCRIPTION"),
		Endpoint:           os.Getenv("KUBERNETES_ENDPOINT"),
		Namespace:          os.Getenv("KUBERNETES_NAMESPACE"),
		DeploymentName:     os.Getenv("KUBERNETES_DEPLOYMENT_NAME"),
		ContainerName:      os.Getenv("KUBERNETES_DEPLOYMENT_CONTAINERNAME"),
		ContainerImage:     os.Getenv("KUBERNETES_DEPLOYMENT_IMAGE_PREFIX"),
		BearerTokenService: tokenProvider,
	}

	command, containerTag := pickCommand(os.Args)

	if command == "ls" {
		// Get a list of pods and their status
		podList, err := cluster.GetPodList()
		if err != nil {
			log.Printf("Unable to retrieve pod list due to %s", err.Error())
			return
		}
		log.Printf("%+v", podList)
	}

	if command == "deploy" {
		// Deploy container named KUBERNETES_DEPLOYMENT_IMAGE_PREFIX:tag
		err := cluster.Deploy(containerTag)
		if err != nil {
			log.Printf("Unable to deploy %q due to %s", containerTag, err.Error())
			return
		}
		log.Println("deployed!")
	}
}

// Sample BearerTokenRetriever

type bearerTokenProvider struct{}

func (*bearerTokenProvider) RetrieveToken() string {
	return os.Getenv("KUBERNETES_ENDPOINT_BEARER_TOKEN")
}

// Helpers

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
