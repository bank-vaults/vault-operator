# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

export PATH := $(abspath bin/):${PATH}

# Target image name
CONTAINER_IMAGE_REF ?= ghcr.io/bank-vaults/vault-operator:dev

CRD_DIR ?= deploy/crd/bases
HELM_DIR ?= deploy/charts/vault-operator

# Default test data
TEST_K8S_VERSION ?= 1.35.0
TEST_VAULT_VERSION ?= 2.0.1
TEST_BANK_VAULTS_VERSION ?= v1.33.1-softhsm
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
check: test lint ## Run tests and lint checks

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

.PHONY: fix
fix: ## Auto-fix lint issues (runs golangci-lint --fix)
	$(GOLANGCI_LINT_BIN) run --fix

.PHONY: license-cache
license-cache: ## Populate license cache
	$(LICENSEI_BIN) cache

.PHONY: license-check
license-check: ## Run license check
	$(LICENSEI_BIN) check
	$(LICENSEI_BIN) header

##@ Development

.PHONY: run
run: install ## Run manager from your host
	OPERATOR_NAME=vault-dev BANK_VAULTS_IMAGE=$(TEST_BANK_VAULTS_IMAGE) go run cmd/main.go -verbose

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

# PLATFORMS lists the target platforms for the manager image. Override via:
#   make docker-buildx PLATFORMS=linux/amd64 CONTAINER_IMAGE_REF=myregistry/myop:0.0.1
# Requires `docker buildx` with BuildKit enabled and push permissions on the target registry.
PLATFORMS ?= linux/arm64,linux/amd64
.PHONY: docker-buildx
docker-buildx: ## Build docker image for cross-platform support
	docker buildx build --push --platform=$(PLATFORMS) --tag ${CONTAINER_IMAGE_REF} .

.PHONY: helm-chart
helm-chart: gen-helm-crds
helm-chart: ## Build helm chart
	@mkdir -p build
	$(HELM_BIN) package -d build/ deploy/charts/vault-operator

##@ Autogeneration

.PHONY: generate
generate: gen-manifests gen-code gen-helm-docs
generate: ## Generate manifests, code, and docs resources

.PHONY: gen-helm-crds
gen-helm-crds: ## Generate CRDs for Helm chart
	./hack/crds-generate.sh $(CRD_DIR) $(HELM_DIR)

.PHONY: gen-manifests
gen-manifests: gen-helm-crds
gen-manifests: ## Generate webhook, RBAC, and CRD resources
	$(CONTROLLER_GEN_BIN) rbac:roleName=vault crd:maxDescLen=0 webhook paths="./..." \
		output:rbac:dir=deploy/rbac \
		output:crd:dir=deploy/crd/bases \
		output:webhook:dir=deploy/webhook

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
undeploy: ## Remove operator resources from the K8s cluster (inverse of `deploy`)
	$(KUSTOMIZE_BIN) build deploy/default | kubectl delete -f -

##@ Dependencies

.PHONY: deps
deps: bin/controller-gen bin/golangci-lint bin/helm bin/helm-docs bin/kind
deps: bin/kustomize bin/licensei bin/setup-envtest
deps: ## Install dependencies

# Dependency versions
GOLANGCI_LINT_VERSION = 2.12.2
LICENSEI_VERSION = 0.9.0
KIND_VERSION = 0.31.0
HELM_VERSION = 4.2.0
CODE_GENERATOR_VERSION = 0.36.1
HELM_DOCS_VERSION = 1.14.2
KUSTOMIZE_VERSION = 5.8.1
CONTROLLER_TOOLS_VERSION = 0.21.0

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


.PHONY: deps-clean
deps-clean: ## Remove all installed dependency binaries (forces fresh install on next `make deps`)
	rm -rf bin/

bin/kustomize: bin/kustomize-v$(KUSTOMIZE_VERSION)
	@ln -sf $(notdir $<) $@

bin/kustomize-v$(KUSTOMIZE_VERSION):
	@mkdir -p bin
	GOBIN=$(DEPS_BIN_PATH) GO111MODULE=on go install sigs.k8s.io/kustomize/kustomize/v5@v$(KUSTOMIZE_VERSION)
	@mv bin/kustomize $@

bin/controller-gen: bin/controller-gen-v$(CONTROLLER_TOOLS_VERSION)
	@ln -sf $(notdir $<) $@

bin/controller-gen-v$(CONTROLLER_TOOLS_VERSION):
	@mkdir -p bin
	GOBIN=$(DEPS_BIN_PATH) go install sigs.k8s.io/controller-tools/cmd/controller-gen@v$(CONTROLLER_TOOLS_VERSION)
	@mv bin/controller-gen $@

bin/setup-envtest:
	@mkdir -p bin
	GOBIN=$(DEPS_BIN_PATH) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

bin/golangci-lint: bin/golangci-lint-v$(GOLANGCI_LINT_VERSION)
	@ln -sf $(notdir $<) $@

bin/golangci-lint-v$(GOLANGCI_LINT_VERSION):
	@mkdir -p bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/v${GOLANGCI_LINT_VERSION}/install.sh | bash -s -- -b ./bin v${GOLANGCI_LINT_VERSION}
	@mv bin/golangci-lint $@

bin/licensei: bin/licensei-v$(LICENSEI_VERSION)
	@ln -sf $(notdir $<) $@

bin/licensei-v$(LICENSEI_VERSION):
	@mkdir -p bin
	curl -sfL https://raw.githubusercontent.com/goph/licensei/master/install.sh | bash -s -- v${LICENSEI_VERSION}
	@mv bin/licensei $@

bin/kind: bin/kind-v$(KIND_VERSION)
	@ln -sf $(notdir $<) $@

bin/kind-v$(KIND_VERSION):
	@mkdir -p bin
	curl -fLo $@ https://kind.sigs.k8s.io/dl/v${KIND_VERSION}/kind-$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m | sed -e "s/aarch64/arm64/; s/x86_64/amd64/")
	@chmod +x $@

bin/helm: bin/helm-v$(HELM_VERSION)
	@ln -sf $(notdir $<) $@

bin/helm-v$(HELM_VERSION):
	@mkdir -p bin
	curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | USE_SUDO=false HELM_INSTALL_DIR=bin DESIRED_VERSION=v${HELM_VERSION} bash
	@mv bin/helm $@

bin/helm-docs: bin/helm-docs-v$(HELM_DOCS_VERSION)
	@ln -sf $(notdir $<) $@

bin/helm-docs-v$(HELM_DOCS_VERSION):
	@mkdir -p bin
	curl -fsSL https://github.com/norwoodj/helm-docs/releases/download/v${HELM_DOCS_VERSION}/helm-docs_${HELM_DOCS_VERSION}_$(shell uname)_$(shell uname -m | sed -e "s/aarch64/arm64/; s/x86_64/x86_64/").tar.gz | tar -zOxf - helm-docs > $@
	@chmod +x $@
