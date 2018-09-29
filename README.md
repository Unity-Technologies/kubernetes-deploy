# kubernetes-deploy

A Go library for deploying to Kubernetes and retrieving status of that deploy.

Useful for deploy bots.

Calls Kubernetes API to update the `image` of an existing deployment and checks that new images are running successful.
