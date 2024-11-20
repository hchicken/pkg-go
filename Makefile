GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

.PHONY: release

# generate release
release:
	@read -p "Enter package name: " package; \
	if [ -z "$$package" ]; then \
		echo "Error: Package name cannot be empty"; \
		exit 1; \
	fi; \
	read -p "Enter version number (e.g., 1.0.0): " version; \
	if [ -z "$$version" ]; then \
		echo "Error: Version number cannot be empty"; \
		exit 1; \
	fi; \
	if ! echo "$$version" | grep -qE '^[0-9]+\.[0-9]+\.[0-9]+$$'; then \
		echo "Error: Version must follow semantic versioning (e.g., 1.0.0)"; \
		exit 1; \
	fi; \
	echo "Creating release for $$package version $$version"; \
	read -p "Continue? [y/N] " confirm; \
	if [ "$$confirm" != "y" ] && [ "$$confirm" != "Y" ]; then \
		echo "Release cancelled"; \
		exit 1; \
	fi; \
	sh release.sh "$$package" "$$version"

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help