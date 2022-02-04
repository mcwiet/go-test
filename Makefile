#################################################################################
# GLOBALS                                                                       #
#################################################################################

# Load environemt variables from .env file and export them (so they can be used by executed processes)
-include .env
export

APP_NAME_API = api
BUILD_DIR = ./dist
CDK_DIR = ./cdk.out
EVENTS_DIR = ./test/_request
GOOS ?= linux
GOARCH ?= amd64
TRUE_CONDITIONS = true TRUE 1

#################################################################################
# COMMANDS                                                                      #
#################################################################################

## Build everything
build: build-api build-infra
	@ echo "✅ Done building everything"

## Build the API application
build-api:
	@ echo "⏳ Start building API..."
	@ GOARCH=${GOARCH} GOOS=${GOOS} go build -o ${BUILD_DIR}/${APP_NAME_API} ./cmd/api
	@ echo "✅ Done building API"

## Build the infrastructure
build-infra:
	@ echo "⏳ Start building infrastructure..."
# CDK errors out if GOOS is set to linux
	@ GOOS="" && cdk synth 
	@ echo "✅ Done building infrastructure"

## Clean all build output
clean:
	@ echo "⏳ Start cleaning..."
	@ rm -rf ${BUILD_DIR}
	@ rm -rf ${CDK_DIR}
	@ echo "✅ Done cleaning"

## Build, package, and update the API application Lambda code (expects infrastructure to have been deployed)
deploy-api: build-api
	@ echo "⏳ Start updating API Lambda code..."
	@ rm -f ${BUILD_DIR}/bootstrap ${BUILD_DIR}/bootstrap.zip
	@ cp ${BUILD_DIR}/${APP_NAME_API} ${BUILD_DIR}/bootstrap
	@ zip -jr ${BUILD_DIR}/bootstrap.zip ${BUILD_DIR}/bootstrap
	@ aws lambda update-function-code --function-name go-api-lambda --zip-file fileb://${BUILD_DIR}/bootstrap.zip
	@ rm -f ${BUILD_DIR}/bootstrap ${BUILD_DIR}/bootstrap.zip
	@ echo "✅ Done updating API Lambda code"

## Deploy the infrastructure
deploy-infra:
	@ echo "⏳ Start deploying infrastructure..."
# CDK errors out if GOOS is set to linux
	@ GOOS="" && cdk deploy 
	@ echo "✅ Done deploying infrastructure"

## Install dependencies
install:
	@ echo "⏳ Start installing dependencies..."
	@ go mod download
	@ echo "✅ Done installing dependencies"

## Invoke the API; set API_REQUEST=[name of request] (e.g. use 'person' for ./test/_request/person.json)
invoke-api: build-infra
	@ echo "⏳ Invoking API with event '${EVENTS_DIR}/${API_REQUEST}.json'..."
	@ sam local invoke go-api-lambda -e ${EVENTS_DIR}/${API_REQUEST}.json  -t ${CDK_DIR}/go-api.template.json
	@ echo "\n✅ Done invoking API"

## Builds all code and runs all tests
release: build test-unit

## Run unit tests on library code (i.e. pkg/ directory)
test-unit: 
	@ rm -rf .coverage
ifeq (${SAVE_TEST_COVERAGE},$(filter ${SAVE_TEST_COVERAGE},${TRUE_CONDITIONS}))
	@ mkdir .coverage
	@ go test ./pkg/... -coverprofile ".coverage/pkg.out" 
else
	@ go test ./pkg/... -cover 
endif

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