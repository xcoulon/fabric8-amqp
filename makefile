MINISHIFT_PROJECT = fabric8
VENDOR_DIR=vendor
LDFLAGS := -w
BUILD_DIR = bin
ARTEMIS_VERSION=2.6.3
MINISHIFT_REGISTRY := $(shell minishift openshift registry)
IMAGE_TAG ?= $(shell git rev-parse HEAD)
GITUNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(GITUNTRACKEDCHANGES),)
IMAGE_TAG := $(IMAGE_TAG)-dirty
endif
BUILD_TIME=`date -u '+%Y-%m-%dT%H:%M:%SZ'`
# Pass in build time variables to main
# LDFLAGS=-ldflags "-X ${PACKAGE_NAME}/controller.Commit=${IMAGE_TAG} -X ${PACKAGE_NAME}/controller.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL := help

# Check that given variables are set and all have non-empty values,
# die with an error otherwise.
#
# Params:
#   1. Variable name(s) to test.
#   2. (optional) Error message to print.
check_defined = \
    $(strip $(foreach 1,$1, \
        $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
    $(if $(value $1),, \
      $(error Undefined $1$(if $2, ($2))))

$(BUILD_DIR): 
	mkdir $(BUILD_DIR)


tools.timestamp:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/lint/golint
	@touch tools.timestamp

deps: tools.timestamp $(VENDOR_DIR) ## Runs dep to vendor project dependencies

$(VENDOR_DIR):
	@echo "checking dependencies..."
	$(GOPATH)/bin/dep ensure -v 

.PHONY: minishift-login
## login to oc minishift
minishift-login:
	@echo "Login to minishift..."
	@oc login --insecure-skip-tls-verify=true -u developer -p developer

.PHONY: minishift-registry-login
## login to the registry in Minishift (to push images)
minishift-registry-login:
	@echo "Login to minishift registry..."
	@eval $$(minishift docker-env) && docker login -u developer -p $(shell oc whoami -t) $(shell minishift openshift registry)

.PHONY: clean-artifacts
clean-artifacts:
	@rm -rf $(BUILD_DIR) && mkdir $(BUILD_DIR)

.PHONY: deploy-activemq-artemis
deploy-activemq-artemis: ## builds and deploy the ActiveMQ Artemis service on Minishift
	eval $$(minishift docker-env) && \
	docker build --build-arg ARTEMIS_VERSION=$(ARTEMIS_VERSION) -t fabric8/activemq-artemis:$(ARTEMIS_VERSION) -f apache-artemis/Dockerfile . && \
	eval $$(minishift docker-env) && \
	docker login -u developer -p $(shell oc whoami -t) $(shell minishift openshift registry) && \
	docker tag fabric8/activemq-artemis:$(ARTEMIS_VERSION)  $(MINISHIFT_REGISTRY)/$(MINISHIFT_PROJECT)/activemq-artemis:$(ARTEMIS_VERSION) && \
	docker push $(MINISHIFT_REGISTRY)/$(MINISHIFT_PROJECT)/activemq-artemis:$(ARTEMIS_VERSION) && \
	oc apply -f openshift/activemq-artemis-deploy.yaml

.PHONY: build-publisher
build-publisher: clean-artifacts ## builds the publisher Docker image
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o $(BUILD_DIR)/publisher publisher/*.go
	eval $$(minishift docker-env) && docker build -t fabric8-publisher:latest \
	  --build-arg BIN_DIR=$(BUILD_DIR) --build-arg BIN_NAME=publisher \
	  -f publisher/Dockerfile .

.PHONY: deploy-publisher
deploy-publisher: build-publisher ## builds and deploy the publisher service on Minishift
	eval $$(minishift docker-env) && \
	docker login -u developer -p $(shell oc whoami -t) $(shell minishift openshift registry) && \
	docker tag fabric8-publisher:latest $(MINISHIFT_REGISTRY)/$(MINISHIFT_PROJECT)/publisher:latest && \
	docker push $(MINISHIFT_REGISTRY)/$(MINISHIFT_PROJECT)/publisher:latest && \
	oc apply -f openshift/publisher.yaml

.PHONY: build-subscriber
build-subscriber: clean-artifacts ## builds the subscriber Docker image
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o $(BUILD_DIR)/subscriber subscriber/*.go
	eval $$(minishift docker-env) && docker build -t fabric8-subscriber:latest \
	  --build-arg BIN_DIR=$(BUILD_DIR) --build-arg BIN_NAME=subscriber \
	  -f subscriber/Dockerfile .
	


.PHONY: deploy-subscribers
deploy-subscribers: build-subscriber ## builds and deploy the subscriber service 1 on Minishift
	eval $$(minishift docker-env) && \
	docker login -u developer -p $(shell oc whoami -t) $(shell minishift openshift registry) && \
	docker tag fabric8-subscriber:latest $(MINISHIFT_REGISTRY)/$(MINISHIFT_PROJECT)/subscriber:latest && \
	docker push $(MINISHIFT_REGISTRY)/$(MINISHIFT_PROJECT)/subscriber:latest 
	

.PHONY: clean-minishift
clean-minishift: minishift-login ## removes the fabric8 project on Minishift
	oc project fabric8 && oc delete project fabric8




.PHONY: help
help: ## Prints this help
	@grep -E '^[^.]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'








