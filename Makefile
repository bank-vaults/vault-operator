# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

export PATH := $(abspath bin/):${PATH}

# Target image name
CONTAINER_IMAGE_REF ?= ghcr.io/bank-vaults/vault-operator:dev

# Default test data
TEST_K8S_VERSION ?= 1.32.0
TEST_VAULT_VERSION ?= 1.14.8
TEST_BANK_VAULTS_VERSION ?= v1.32.0-softhsm
TEST_BANK_VAULTS_IMAGE ?= ghcr.io/bank-vaults/bank-vaults:$(TEST_BANK_VAULTS_VERSION)
TEST_OPERATOR_VERSION ?= $(lastword $(subst :, ,$(CONTAINER_IMAGE_REF)))
TEST_KIND_CLUSTER ?= vault-operator

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

##@ General

# Targets commented with ## will be visible in "make help" info.
# Comments marked with ##@ will be used as categories for a group of targets.

.PHONY: help
default: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Checks

.PHONY: check
check: test lint## Run tests and lint checks

.PHONY: test
test: ## Run tests
	KUBEBUILDER_ASSETS="$(shell $(SETUP_ENVTEST_BIN) use $(TEST_K8S_VERSION) --bin-dir ./bin -p path)" \
		go test -race -v ./... -coverprofile cover.out

.PHONY: test-acceptance
test-acceptance: ## Run acceptance tests
	VAULT_VERSION=$(TEST_VAULT_VERSION) BANK_VAULTS_VERSION=$(TEST_BANK_VAULTS_VERSION) OPERATOR_VERSION=$(TEST_OPERATOR_VERSION) \
		go test -race -v -timeout 900s -tags kubeall ./test

.PHONY: lint
lint: lint-go lint-helm lint-docker lint-yaml
lint: ## Run lint checks

.PHONY: lint-go
lint-go:
	$(GOLANGCI_LINT_BIN) run

.PHONY: lint-helm
lint-helm:
	$(HELM_BIN) lint deploy/charts/vault-operator

.PHONY: lint-docker
lint-docker:
	$(HADOLINT_BIN) Dockerfile

.PHONY: lint-yaml
lint-yaml:
	$(YAMLLINT_BIN) $(if ${CI},-f github,) --no-warnings .

.PHONY: fmt
fmt: ## Format code
	$(GOLANGCI_LINT_BIN) run --fix

.PHONY: license-check
license-check: ## Run license check
	$(LICENSEI_BIN) check
	$(LICENSEI_BIN) header

##@ Development

.PHONY: run
run: deploy ## Run manager from your host
	OPERATOR_NAME=vault-dev BANK_VAULTS_IMAGE=$(TEST_BANK_VAULTS_IMAGE)  go run cmd/main.go -verbose

.PHONY: up
up: ## Start kind development environment
	$(KIND_BIN) create cluster --name $(TEST_KIND_CLUSTER)

.PHONY: down
down: ## Destroy kind development environment
	$(KIND_BIN) delete cluster --name $(TEST_KIND_CLUSTER)

.PHONY: prepare-kind
prepare-kind: ## Prepare kind cluster for development/testing
	docker pull ghcr.io/bank-vaults/bank-vaults:$(TEST_BANK_VAULTS_VERSION)
	docker pull hashicorp/vault:$(TEST_VAULT_VERSION)

	$(KIND_BIN) load docker-image ${CONTAINER_IMAGE_REF} --name $(TEST_KIND_CLUSTER)
	$(KIND_BIN) load docker-image ghcr.io/bank-vaults/bank-vaults:$(TEST_BANK_VAULTS_VERSION) --name $(TEST_KIND_CLUSTER)
	$(KIND_BIN) load docker-image hashicorp/vault:$(TEST_VAULT_VERSION) --name $(TEST_KIND_CLUSTER)

##@ Build

.PHONY: build
build: ## Build binary
	@mkdir -p build
	go build -race -o build/manager ./cmd

.PHONY: artifacts
artifacts: docker-build helm-chart
artifacts: ## Build docker image and helm chart

# If you wish built the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64 ). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build: ## Build docker image
	docker build -t ${CONTAINER_IMAGE_REF} .

# PLATFORMS defines the target platforms for  the manager image be build to provide support to multiple
# architectures. (i.e. make docker-buildx CONTAINER_IMAGE_REF=myregistry/mypoperator:0.0.1). To use this option you need to:
# - able to use docker buildx . More info: https://docs.docker.com/build/buildx/
# - have enable BuildKit, More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image for your registry (i.e. if you do not inform a valid value via CONTAINER_IMAGE_REF=<myregistry/image:<tag>> then the export will fail)
# To properly provided solutions that supports more than one platform you should use this option.
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: docker-buildx
docker-buildx: ## Build docker image for cross-platform support
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- docker buildx create --name project-v3-builder
	docker buildx use project-v3-builder
	- docker buildx build --push --platform=$(PLATFORMS) --tag ${CONTAINER_IMAGE_REF} -f Dockerfile.cross .
	- docker buildx rm project-v3-builder
	rm Dockerfile.cross

.PHONY: helm-chart
helm-chart: ## Build helm chart
	@mkdir -p build
	$(HELM_BIN) package -d build/ deploy/charts/vault-operator

##@ Autogeneration

.PHONY: generate
generate: gen-manifests gen-code gen-helm-docs
generate: ## Generate manifests, code, and docs resources

.PHONY: gen-manifests
gen-manifests: ## Generate webhook, RBAC, and CRD resources
	$(CONTROLLER_GEN_BIN) rbac:roleName=vault crd:maxDescLen=0 webhook paths="./..." \
		output:rbac:dir=deploy/rbac \
		output:crd:dir=deploy/crd/bases \
		output:webhook:dir=deploy/webhook
	cp deploy/crd/bases/vault.banzaicloud.com_vaults.yaml deploy/charts/vault-operator/crds/crd.yaml

.PHONY: gen-code
gen-code: ## Generate deepcopy, client, lister, and informer objects
	$(CONTROLLER_GEN_BIN) object:headerFile="hack/custom-boilerplate.go.txt" paths="./..."
	./hack/update-codegen.sh v${CODE_GENERATOR_VERSION}

.PHONY: gen-helm-docs
gen-helm-docs: ## Generate Helm chart documentation
	$(HELM_DOCS_BIN) -s file -c deploy/charts/ -t README.md.gotmpl

##@ Deployment

.PHONY: install
install: gen-manifests ## Install CRDs into the K8s cluster
	$(KUSTOMIZE_BIN) build deploy/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: gen-manifests ## Uninstall CRDs from the K8s cluster
	$(KUSTOMIZE_BIN) build deploy/crd | kubectl delete -f -

.PHONY: deploy
deploy: gen-manifests ## Deploy resources to the K8s cluster
	cd deploy/manager && $(PWD)/$(KUSTOMIZE_BIN) edit set image controller=${CONTAINER_IMAGE_REF}
	$(KUSTOMIZE_BIN) build deploy/default | kubectl apply -f -

.PHONY: undeploy
clean: ## Clean resources from the K8s cluster
	$(KUSTOMIZE_BIN) build deploy/default | kubectl delete -f -

##@ Dependencies

.PHONY: deps
deps: bin/controller-gen bin/golangci-lint bin/helm bin/helm-docs bin/kind
deps: bin/kurun bin/kustomize bin/licensei bin/setup-envtest
deps: ## Install dependencies

# Dependency versions
GOLANGCI_LINT_VERSION = 2.0.2
LICENSEI_VERSION = 0.9.0
KIND_VERSION = 0.25.0
HELM_VERSION = 3.16.3
KURUN_VERSION = 0.7.0
CODE_GENERATOR_VERSION = 0.27.1
HELM_DOCS_VERSION = 1.14.2
KUSTOMIZE_VERSION = 5.5.0
CONTROLLER_TOOLS_VERSION = 0.16.5

# Dependency binaries
GOLANGCI_LINT_BIN := golangci-lint
LICENSEI_BIN := licensei
KIND_BIN := kind
HELM_BIN := helm
HELM_DOCS_BIN := helm-docs
SETUP_ENVTEST_BIN := setup-envtest
KUSTOMIZE_BIN := kustomize
CONTROLLER_GEN_BIN := controller-gen

# TODO: add support for hadolint and yamllint dependencies
HADOLINT_BIN := hadolint
YAMLLINT_BIN := yamllint

# If we have "bin" dir, use those binaries instead
ifneq ($(wildcard ./bin/.),)
	GOLANGCI_LINT_BIN := bin/$(GOLANGCI_LINT_BIN)
	LICENSEI_BIN := bin/$(LICENSEI_BIN)
	KIND_BIN := bin/$(KIND_BIN)
	HELM_BIN := bin/$(HELM_BIN)
	HELM_DOCS_BIN := bin/$(HELM_DOCS_BIN)
	SETUP_ENVTEST_BIN := bin/$(SETUP_ENVTEST_BIN)
	KUSTOMIZE_BIN := bin/$(KUSTOMIZE_BIN)
	CONTROLLER_GEN_BIN := bin/$(CONTROLLER_GEN_BIN)
endif

# full path to "bin" required for go install
DEPS_BIN_PATH ?= $(shell pwd)/bin

bin/kustomize:
	@mkdir -p bin
	GOBIN=$(DEPS_BIN_PATH) GO111MODULE=on go install sigs.k8s.io/kustomize/kustomize/v5@v$(KUSTOMIZE_VERSION)

bin/controller-gen:
	@mkdir -p bin
	GOBIN=$(DEPS_BIN_PATH) go install sigs.k8s.io/controller-tools/cmd/controller-gen@v$(CONTROLLER_TOOLS_VERSION)

bin/setup-envtest:
	@mkdir -p bin
	GOBIN=$(DEPS_BIN_PATH) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

bin/golangci-lint:
	@mkdir -p bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- v${GOLANGCI_LINT_VERSION}

bin/licensei:
	@mkdir -p bin
	curl -sfL https://raw.githubusercontent.com/goph/licensei/master/install.sh | bash -s -- v${LICENSEI_VERSION}

bin/kind:
	@mkdir -p bin
	curl -Lo bin/kind https://kind.sigs.k8s.io/dl/v${KIND_VERSION}/kind-$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m | sed -e "s/aarch64/arm64/; s/x86_64/amd64/")
	@chmod +x bin/kind

bin/helm:
	@mkdir -p bin
	curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | USE_SUDO=false HELM_INSTALL_DIR=bin DESIRED_VERSION=v${HELM_VERSION} bash
	@chmod +x bin/helm

bin/helm-docs:
	@mkdir -p bin
	curl -L https://github.com/norwoodj/helm-docs/releases/download/v${HELM_DOCS_VERSION}/helm-docs_${HELM_DOCS_VERSION}_$(shell uname)_x86_64.tar.gz | tar -zOxf - helm-docs > ./bin/helm-docs
	@chmod +x bin/helm-docs

bin/kurun:
	@mkdir -p bin
	curl -Lo  bin/kurun https://github.com/banzaicloud/kurun/releases/download/${KURUN_VERSION}/kurun-$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m | sed -e "s/aarch64/arm64/; s/x86_64/amd64/")
	@chmod +x  bin/kurun
