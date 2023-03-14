MAIN_APP_NAME := main-app
MAIN_APP_VERSION := v0.0.1
MAIN_APP_IMG := $(MAIN_APP_NAME):$(MAIN_APP_VERSION)
PROCESSOR_APP_NAME := processor-app
PROCESSOR_APP_VERSION := v0.0.1
PROCESSOR_APP_IMG := $(PROCESSOR_APP_NAME):$(PROCESSOR_APP_VERSION)

default: help

buildAllImages: buildMainImg buildProcessorImg ## Build Main and Processor Images

buildMainImg: ## Build Main Image
	@docker build -t $(MAIN_APP_IMG) -f ./pkg/main/Dockerfile .

buildProcessorImg: ## Build Processor (secondary) Image
	@docker build -t $(PROCESSOR_APP_IMG) -f ./pkg/processor/Dockerfile .

run: ## Run App
	@docker-compose -f ./docker-compose.yml up -d

runRequiredServices: ## Run Required Services (Jaeger, Kafka, OpenTelemetry)
	@docker-compose -f ./deployments/docker-compose.yml up -d

clean: ## Clean docker
	@docker-compose -f ./docker-compose.yml up -d
	@docker-compose -f ./deployments/docker-compose.yml down

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[38;5;69m%-30s\033[38;5;38m %s\033[0m\n", $$1, $$2}'