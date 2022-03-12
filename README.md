# Go Test Project

## About

This is a test project for getting familiar with Go. On top of learning simple things like language syntax and directory structure, the intent is to focus on working with GraphQL and AWS CDK with Go. Clean architecture principles are followed, separating code into layers and frequently using dependency injection.

The application is a simple API for interacting with `pet` objects with create, read, update, and delete operations. The infrastructure is powered by AWS services such as AppSync, Lambda, DynamoDB, and Cognito. GitHub Actions power the CI/CD pipeline which use a `staging` environment for pull requests and a `production` environment for code released to the main branch. SonarCloud scans the code for potential bugs, security issues, unmaintainable code, and test coverage reporting.

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
1. Run `make build-env-file` to create the `.env` environment file (pulls values from the freshly deployed infrastructure)
1. Update the `.env` file with credentials for a test user (for activities like automated integration testing)
1. Run `make create-test-user` to add the test user to the Cognito User Pool
1. Run `make promote-test-user` to make the test user an admin, capable of executing all necessary API commands
1. Run `make test-unit` to run unit tests locally
1. Run `make test-integration` to run integration tests against the deployed environment in AWS
1. Run `make invoke-api-sam API_REQUEST=pet` to invoke the API Lambda locally in Docker, using requests stored in `test/_request/`

## Developer Notes

### API Schema

The API abides by a few general principles:

- All operations utilize a single `input` object for receiving data (e.g. `CreatePetInput`)
- All mutations return a mutation-specific `payload` object (e.g. `CreatePetPayload`)
- All queries return 'basic model' objects (e.g. `Pet`) or `connection` objects (e.g. `PetConnection`)
- Avoid generic `update` mutations; prefer mutations which change small amounts of data (e.g. `updatePetOwner`)

A few established GraphQL schemas, like GitHub's, were referenced for patterns as well.

### Design Comments

- Each API request has 3 layers: controller (parse the HTTP request and call services), services (run business logic) and data (interact with the data store)
- Dependency injection is used frequently to make unit testing easier and abide by clean architecture (enable use of stubs and mocks)
- Initial unit tests are simple and generally test a "working path" and an "error path"
- Authorization is mix of RBAC and ABAC - a user may be authorized to perform an action based on a role/group (e.g. admin) or based on attributes (e.g. requestor is the owner of the target pet)
  - **Pros**: easy mapping of security requirements to code, unit testable, no separate storage of permissions (e.g. data table which lists each pet permission a user has)
  - **Cons**: certain permission changes could require code change rather than changing data at runtime, risk of unintentionally removing a user's access to a resource which they could previously access
  - Took approach of custom code rather than library like Casbin to keep code simple (no need to learn domain-specific policy languages) and more easily implement certain features (e.g. check if one of user's groups is the 'admin' group)
- Dependencies between CDK stacks are implemented as SSM parameters (rather than stack outputs / exports); this leads to reduced coupling and allows stacks to be deleted without first deleting their dependent stacks ([see here for more context](https://tusharsharma.dev/posts/aws-cfn-with-ssm-parameters))
- To delete an environment, delete the CloudFormation stacks and manually delete any resources where the status was marked as **DELETE_SKIPPED** (such as a DynamoDB Table or Cognito User Pool)

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
1. [GitHub GraphQL Schema](https://docs.github.com/en/graphql/overview/public-schema)
