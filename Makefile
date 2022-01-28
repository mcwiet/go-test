#################################################################################
# GLOBALS                                                                       #
#################################################################################

APP_NAME_API = "api"
BUILD_DIR = "./dist"

#################################################################################
# COMMANDS                                                                      #
#################################################################################

## Build everything
build: build-api
	@ echo "✅ Done building everything"

## Build the API application
build-api:
	@ echo "⏳ Start building API..."
	@ go build -o ${BUILD_DIR}/${APP_NAME_API} ./cmd/api
	@ echo "✅ Done building API"

## Clean all build output
clean:
	@ echo "⏳ Start cleaning..."
	@ rm -rf ${BUILD_DIR}
	@ echo "✅ Done cleaning"

## Start the API application
start-api:
	@ echo "🏎  Starting the API"
	@ ${BUILD_DIR}/${APP_NAME_API}

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