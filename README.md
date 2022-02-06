# Go Test Project

## About

This is a test project for getting famliar with Go. On top of learning simple things like language syntax and directory structure, the intent is to focus on working with GraphQL and AWS CDK with Go. Clean architecture principles are followed, separating code into layers and fequently using dependency injection.

The application is a simple API for interacting with `person` objects with create, read, update, and delete operations. The infrastructure is powered by AWS services such as AppSync, Lambda, DynamoDB, and Cognito. GitHub Actions power the CI/CD pipeline which use a `staging` environment for pull requests and a `production` environment for code released to the main branch. SonarCloud scans the code for potential bugs, security issues, unmaintainable code, and test coverage reporting.

## Getting Started

### Prerequisites

1. [Go](https://go.dev/doc/install): Language SDK
1. [AWS CDK](https://docs.aws.amazon.com/cdk/v2/guide/cli.html): AWS infrastructure management (be sure to [bootstrap](https://docs.aws.amazon.com/cdk/v2/guide/bootstrapping.html) the account/region)
1. [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html): Preferably executing with a role that has admin privileges to the AWS account
1. [AWS SAM](https://aws.amazon.com/serverless/sam/) + [Docker](https://www.docker.com/products/docker-desktop): Locally invoking / testing Lambda functions
1. [jq](https://stedolan.github.io/jq/): Parsing JSON responses (used in Makefile commands)
1. [Make](https://www.gnu.org/software/make/): Running common, preset commands

This document will assume you have the ability to run commands from the Makefile.

## Usage

The instructions below provide steps for setting up an environment both locally and in AWS. For a full list of supported Makefile commands, use either run `make` or `make help` to list all commands and a brief description.

1. Run `make build-infra deploy-infra` to build and deploy a development environment to an AWS account
1. Run `make build-env` to create the `.env` environment file (pulls values from the freshly deployed infrastructure)
1. Update the `.env` file with credentials for a test user (for activities such as automated integration testing), then run `make create-test-user` to add that user to the Cognito User Pool
1. Run `make test-unit` to run unit tests locally
1. Run `make test-integration` to run integration tests against the deployed environment in AWS
1. Run `make invoke-api-sam API_REQUEST=person` to invoke the API Lambda locally in Docker, using requests stored in `test/_request/`

## Developer Notes

### Design Comments

- Each API request has 3 layers: controller (parse the HTTP request and call services), services (run business logic) and data (interact with the data store)
- Dependency injection is used as much as possible to make unit testing easier (enable use of stubs and mocks)
- Initial unit tests are simple and generally test a "working path" and an "error path"
- Dependencies between CDK stacks are implemented as SSM parameters (rather than stack outputs / exports); this leads to reduced coupling and allows stacks to be deleted without first deleting their dependenent stacks ([see here for more context](https://tusharsharma.dev/posts/aws-cfn-with-ssm-parameters))
- To delete an environemt, delete the CloudFormation stacks and manually delete any resources where the status was marked as **DELETE_SKIPPED** (such as a DynamoDB Table or Cognito User Pool)

### Adding a New API

These are the recommended high level steps to adding functionality to the API:

1. Update the schema (`api/schema.graphql`)
1. Update the CDK
1. Create a controller and update the API entrypoint (`cmd/api/main.go`)
1. Create the service
1. Create the data access object
1. Add tests (if not already done)

### References

1. [Directory structure recommendations](https://github.com/golang-standards/project-layout)
1. [GraphQL knowledge base](https://graphql.org/learn/)
1. [GraphQL specification](https://spec.graphql.org/)
