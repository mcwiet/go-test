# Go Test Project

## About

This is a test project for getting famliar with Go. On top of learning simple things like language syntax and directory structure, the intent is to focus on working with GraphQL and AWS CDK with Go. Clean architecture principles are followed as well, separating code into layers and fequently using dependency injection.

Code is setup to allow easy

## Getting Started

### Required Prerequisites

1. [Go](https://go.dev/doc/install): Language SDK
1. [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html): Preferably executing with a role that has admin privileges to the AWS account
1. [AWS CDK](https://docs.aws.amazon.com/cdk/v2/guide/cli.html): AWS infrastructure management

### Suggested Prerequisites

1. Ability to run makefile commands
1. [AWS SAM](https://aws.amazon.com/serverless/sam/) + [Docker](https://www.docker.com/products/docker-desktop): Locally invoking / testing Lambda functions
1. [jq](https://stedolan.github.io/jq/): Parsing JSON responses (used in Makefile commands)

This document will assume you have the ability to run commands from the Makefile.

### Installation

1. Clone repo
1. Install project dependencies - `make install`

## Usage

See instructions below for running the application locally or in an AWS account. For a full list of supported Makefile commands, use either `make` or `make help`.

### Local

The API can be run locally using SAM and Docker. Requests are expected in the format of a standard AppSync request - reusable examples can be placed in the `test/_data/_request` folder.

1. Build everything - `make build`
1. Invoke the API - `make invoke-api API_REQUEST=person`

Note that the invoke command does not automatically rebuild the API package. When rapidly testing changes with the same request, it can be helpful to use a command such as `make build-api invoke-api API_REQUEST=person`. Note that SAM uses the infrastructure build output as well so be mindful of rebuilding the infrastructure package as well.

### AWS Cloud

1. Build everything - `make build`
1. Deploy the infrastructure - `make deploy-infra`

Once the infrastructure has been deployed, changes to code can be quickly uploaded directly to the Lambda resolver using `make deploy-api` (saves time by not updating the entire CloudFormation infrastructure stack).

## References

1. [Directory structure recommendations](https://github.com/golang-standards/project-layout)
1. [GraphQL knowledge base](https://graphql.org/learn/)
1. [GraphQL specification](https://spec.graphql.org/)
