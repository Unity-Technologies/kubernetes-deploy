# kubernetes-deploy

A Go library for deploying to Kubernetes and retrieving status of that deploy.

Useful for deploy bots.

Included `main.go` sample app demonstrates its features.

# How does it work?

Calls Kubernetes API to update the `image` of an existing deployment and checks that new images are running successful.

Authorization for Kubernetes API via `Authorization: Bearer <token>` header. Offers a `BearerTokenRetriever` interface to provide flexibility -- either grab a long-lived token once from an environment variables, or refresh and retrieve short-lived tokens prior to each call.

NOTE that you need to have a running deployment of your app in your Kubernetes cluster first. Deploying just changes an existing deployment's `image` to a new Docker container.

# Getting Started

Copy the sample .env file and fill in your values.

    cp .env-sample .env

`main.go` supports two commands:

    go run main.go ls     # to list current pods in deployment

    go run main.go deploy 77d0ea51fdc30234918f2726d26479c66b7f777   # deploy container with this tag
