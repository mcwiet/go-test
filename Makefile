#################################################################################
# GLOBALS                                                                       #
#################################################################################

ENV_FILE = ./.env

# Load environemt variables from env file and export them (so they can be used by executed processes)
-include ${ENV_FILE}
export

# Constants
APP_NAME_API = api
AWS_LAMBDA_GOOS = linux
AWS_LAMBDA_GOARCH = amd64
BUILD_DIR = ./dist
CDK_DIR = ./cdk.out
CMD_API_URL = aws ssm get-parameter --name /go/${ENV}/appsync-url | jq '.Parameter.Value'
CMD_USER_POOL_ID = aws ssm get-parameter --name /go/${ENV}/user-pool-id | jq '.Parameter.Value'
CMD_USER_POOL_APP_CLIENT_ID = aws ssm get-parameter --name /go/${ENV}/user-pool-api-client-id | jq '.Parameter.Value'
EVENTS_DIR = ./test/_request
TRUE_CONDITIONS = true TRUE 1

# Conditional constants
ENV ?= development

#################################################################################
# COMMANDS                                                                      #
#################################################################################

## Build everything
build: build-api build-infra
	@ echo "✅ Done building everything"

## Build the API application
build-api:
	@ echo "⏳ Start building API..."
	@ go build -o ${BUILD_DIR}/${APP_NAME_API} ./cmd/api
	@ echo "✅ Done building API"

## Build the infrastructure
build-infra:
	@ echo "⏳ Start building ${ENV} infrastructure..."
	@ cdk synth 
	@ echo "✅ Done building ${ENV} infrastructure"

## Build an env file containing values for infrastructure-dependent environment variables
build-env-file:
	@ $(eval USER_POOL_ID=$(shell ${CMD_USER_POOL_ID}))
	@ $(eval USER_POOL_APP_CLIENT_ID=$(shell ${CMD_USER_POOL_APP_CLIENT_ID}))
	@ $(eval API_URL=$(shell ${CMD_API_URL}))
	@ echo "🚨 FYI: Deleting existing ${ENV_FILE} file"
	@ rm -f ${ENV_FILE}
	@ echo "⏳ Start building ${ENV_FILE} file for ${ENV}..."
	@ echo "# AUTOGENERATED VALUES" >> ${ENV_FILE}
	@ echo "ENV=${ENV}" >> ${ENV_FILE}
	@ echo "AWS_ACCOUNT=${AWS_ACCOUNT}" >> ${ENV_FILE}
	@ echo "AWS_REGION=${AWS_REGION}" >> ${ENV_FILE}
	@ echo "USER_POOL_ID=${USER_POOL_ID}" >> ${ENV_FILE}
	@ echo "USER_POOL_APP_CLIENT_ID=${USER_POOL_APP_CLIENT_ID}" >> ${ENV_FILE}
	@ echo "API_URL=${API_URL}" >> ${ENV_FILE}
	@ echo "" >> ${ENV_FILE}
	@ echo "# MANUAL VALUES" >> ${ENV_FILE}
	@ echo "TEST_USER_EMAIL=${TEST_USER_EMAIL}" >> ${ENV_FILE}
	@ echo "TEST_USER_PASSWORD=${TEST_USER_PASSWORD}" >> ${ENV_FILE}
	@ echo "✅ Done building ${ENV_FILE} file for ${ENV}"

## Create a user in the user pool for the current environment
create-test-user:
	@ $(eval USER_POOL_ID=$(shell ${CMD_USER_POOL_ID}))
ifndef TEST_USER_EMAIL
	@ echo "🚨 MANUAL ACTION: Set value for TEST_USER_EMAIL"
else
ifndef TEST_USER_PASSWORD
	@ echo "🚨 MANUAL ACTION: Set value for TEST_USER_PASSWORD"
else
	@ echo "⏳ Start creating ${ENV} user '${TEST_USER_EMAIL}'..."
	@ aws cognito-idp admin-create-user --user-pool-id ${USER_POOL_ID} --username ${TEST_USER_EMAIL} --message-action SUPPRESS
	@ aws cognito-idp admin-set-user-password --user-pool-id ${USER_POOL_ID} --username ${TEST_USER_EMAIL} --password ${TEST_USER_PASSWORD} --permanent
	@ echo "Updated attributes - password set"
	@ aws cognito-idp admin-update-user-attributes --user-pool-id ${USER_POOL_ID} --username ${TEST_USER_EMAIL} --user-attributes Name=email_verified,Value=true
	@ echo "Updated attributes - email verified"
	@ echo "✅ Done creating ${ENV} user '${TEST_USER_EMAIL}'..."
	@ echo "🚨 MANUAL ACTION: Update ${ENV_FILE} file with user credentials (if needed)"
endif
endif

## Clean all build output
clean:
	@ echo "⏳ Start cleaning..."
	@ rm -rf ${BUILD_DIR}
	@ rm -rf ${CDK_DIR}
	@ echo "✅ Done cleaning"

## Delete a user in the user pool for the current environment
delete-test-user:
	@ $(eval USER_POOL_ID=$(shell ${CMD_USER_POOL_ID}))
ifndef TEST_USER_EMAIL
	@ echo "🚨 MANUAL ACTION: Set value for TEST_USER_EMAIL"
else
	@ echo "⏳ Start deleting ${ENV} user '${TEST_USER_EMAIL}'..."
	@ aws cognito-idp admin-delete-user --user-pool-id ${USER_POOL_ID} --username ${TEST_USER_EMAIL}
	@ echo "✅ Done deleting ${ENV} user '${TEST_USER_EMAIL}'..."
endif

## Deploy the infrastructure
deploy-infra:
	@ echo "⏳ Start deploying ${ENV} infrastructure..."
	@ cdk deploy --all
	@ echo "✅ Done deploying ${ENV} infrastructure"

## Install dependencies
install:
	@ echo "⏳ Start installing dependencies..."
	@ go mod download
	@ echo "✅ Done installing dependencies"

## Invoke the API; set API_REQUEST=[name of request] (e.g. use 'person' for ./test/_request/person.json)
invoke-api: build-infra
	@ echo "⏳ Invoking API with event '${EVENTS_DIR}/${API_REQUEST}.json'..."
	@ sam local invoke go-${ENV}-api-lambda -e ${EVENTS_DIR}/${API_REQUEST}.json -t ${CDK_DIR}/go-${ENV}-api.template.json
	@ echo "\n✅ Done invoking API"

## Run integration tests
test-integration:
	@ echo "⏳ Start running ${ENV} integration tests..."
	@ go test ./test/integration/... -v
	@ echo "✅ Done running ${ENV} integration tests"

## Run unit tests on library code (i.e. pkg/ directory)
test-unit: 
	@ echo "⏳ Start running unit tests..."
	@ rm -rf .coverage
ifeq (${SAVE_TEST_COVERAGE},$(filter ${SAVE_TEST_COVERAGE},${TRUE_CONDITIONS}))
	@ mkdir .coverage
	@ go test ./pkg/... -coverprofile ".coverage/pkg.out" 
else
	@ go test ./pkg/... -cover 
endif
	@ echo "✅ Done running unit tests"

## Build, package, and update the API application Lambda code (expects infrastructure to have been deployed)
update-api:
	@ echo "⏳ Start updating API Lambda code..."
	@ GOARCH=${AWS_LAMBDA_GOARCH} GOOS=${AWS_LAMBDA_GOOS} go build -o ${BUILD_DIR}/${APP_NAME_API} ./cmd/api
	@ rm -f ${BUILD_DIR}/bootstrap ${BUILD_DIR}/bootstrap.zip
	@ cp ${BUILD_DIR}/${APP_NAME_API} ${BUILD_DIR}/bootstrap
	@ zip -jr ${BUILD_DIR}/bootstrap.zip ${BUILD_DIR}/bootstrap
	@ aws lambda update-function-code --function-name go-${ENV}-api-lambda --zip-file fileb://${BUILD_DIR}/bootstrap.zip
	@ rm -f ${BUILD_DIR}/bootstrap ${BUILD_DIR}/bootstrap.zip
	@ echo "✅ Done updating API Lambda code"

#################################################################################
# RESERVED                                                                      #
#################################################################################

.DEFAULT_GOAL := help
.PHONY: help
help:
	@echo "$$(tput bold)Available rules:$$(tput sgr0)"
	@echo
	@sed -n -e "/^## / { \
		h; \
		s/.*//; \
		:doc" \
		-e "H; \
		n; \
		s/^## //; \
		t doc" \
		-e "s/:.*//; \
		G; \
		s/\\n## /---/; \
		s/\\n/ /g; \
		p; \
	}" ${MAKEFILE_LIST} \
	| LC_ALL='C' sort --ignore-case \
	| awk -F '---' \
		-v ncol=$$(tput cols) \
		-v indent=19 \
		-v col_on="$$(tput setaf 6)" \
		-v col_off="$$(tput sgr0)" \
	'{ \
		printf "%s%*s%s ", col_on, -indent, $$1, col_off; \
		n = split($$2, words, " "); \
		line_length = ncol - indent; \
		for (i = 1; i <= n; i++) { \
			line_length -= length(words[i]) + 1; \
			if (line_length <= 0) { \
				line_length = ncol - indent - length(words[i]) - 1; \
				printf "\n%*s ", -indent, " "; \
			} \
			printf "%s ", words[i]; \
		} \
		printf "\n"; \
	}' \
	| more $(shell test $(shell uname) = Darwin && echo '--no-init --raw-control-chars')