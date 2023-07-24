# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

# Target image name
IMG ?= ghcr.io/bank-vaults/vault-operator:dev

# Default test data
TEST_K8S_VERSION ?= 1.27.1
TEST_VAULT_VERSION ?= 1.14.0
TEST_BANK_VAULTS_VERSION ?= 1.20.0
TEST_BANK_VAULTS_IMAGE ?= ghcr.io/bank-vaults/bank-vaults:$(TEST_BANK_VAULTS_VERSION)
TEST_OPERATOR_VERSION ?= $(lastword $(subst :, ,$(IMG)))
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

.PHONY: fmt
fmt: ## Run go fmt against code
	$(GOLANGCI_LINT) run --fix

.PHONY: lint-go
lint-go: # Run golang lint check
	$(GOLANGCI_LINT) run $(if ${CI},--out-format github-actions,)

.PHONY: lint-helm
lint-helm: # Run helm lint check
	$(HELM) lint deploy/charts/vault-operator

.PHONY: lint-docker
lint-docker: # Run Dockerfile lint check
	$(HADOLINT) Dockerfile

.PHONY: lint-yaml
lint-yaml: # Run yaml lint check
	$(YAMLLINT) $(if ${CI},-f github,) --no-warnings .

.PHONY: lint
lint: lint-go lint-helm lint-docker lint-yaml
lint: ## Run lint checks

.PHONY: test
test: ## Run tests
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(TEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" \
		go test -race -v ./... -coverprofile cover.out

.PHONY: test-acceptance
test-acceptance: ## Run acceptance tests. If running on a local kind cluster, run "make import-test" before this
	VAULT_VERSION=$(TEST_VAULT_VERSION) BANK_VAULTS_VERSION=$(TEST_BANK_VAULTS_VERSION) OPERATOR_VERSION=$(TEST_OPERATOR_VERSION) \
		go test -race -v -timeout 900s -tags kubeall ./test

.PHONY: license-check
license-check: ## Run license check
	$(LICENSEI) check
	$(LICENSEI) header

.PHONY: check
check: lint test ## Run lint checks and tests

##@ Development

.PHONY: run
run: deploy ## Run manager from your host
	OPERATOR_NAME=vault-dev BANK_VAULTS_IMAGE=$(TEST_BANK_VAULTS_IMAGE)  go run cmd/main.go -verbose

.PHONY: up
up: ## Start kind development environment
	$(KIND) create cluster --name $(TEST_KIND_CLUSTER)

.PHONY: down
down: ## Destroy kind development environment
	$(KIND) delete cluster --name $(TEST_KIND_CLUSTER)

.PHONY: import-image
import-image: docker-build ## Import manager image to kind image repository
	$(KIND) load docker-image ${IMG} --name $(TEST_KIND_CLUSTER)

.PHONY: import-test
import-test: import-image ## Import images required for tests to kind image repository
	docker pull ghcr.io/bank-vaults/bank-vaults:$(TEST_BANK_VAULTS_VERSION)
	docker pull hashicorp/vault:$(TEST_VAULT_VERSION)

	$(KIND) load docker-image ghcr.io/bank-vaults/bank-vaults:$(TEST_BANK_VAULTS_VERSION) --name $(TEST_KIND_CLUSTER)
	$(KIND) load docker-image hashicorp/vault:$(TEST_VAULT_VERSION) --name $(TEST_KIND_CLUSTER)

##@ Build

.PHONY: build
build: ## Build manager binary
	@mkdir -p build
	go build -race -o build/manager ./cmd

# If you wish built the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64 ). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build: ## Build docker image
	docker build -t ${IMG} .

# PLATFORMS defines the target platforms for  the manager image be build to provide support to multiple
# architectures. (i.e. make docker-buildx IMG=myregistry/mypoperator:0.0.1). To use this option you need to:
# - able to use docker buildx . More info: https://docs.docker.com/build/buildx/
# - have enable BuildKit, More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image for your registry (i.e. if you do not inform a valid value via IMG=<myregistry/image:<tag>> then the export will fail)
# To properly provided solutions that supports more than one platform you should use this option.
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: docker-buildx
docker-buildx: ## Build docker image for cross-platform support
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- docker buildx create --name project-v3-builder
	docker buildx use project-v3-builder
	- docker buildx build --push --platform=$(PLATFORMS) --tag ${IMG} -f Dockerfile.cross .
	- docker buildx rm project-v3-builder
	rm Dockerfile.cross

.PHONY: helm-chart
helm-chart: ## Build helm chart
	@mkdir -p build
	$(HELM) package -d build/ deploy/charts/vault-operator

.PHONY: artifacts
artifacts: docker-build helm-chart
artifacts: ## Build docker image and helm chart


##@ Autogeneration

.PHONY: gen-manifests
gen-manifests: ## Generate webhook, RBAC, and CRD resources
	$(CONTROLLER_GEN) rbac:roleName=vault crd:maxDescLen=0 webhook paths="./..." \
		output:rbac:dir=deploy/rbac \
		output:crd:dir=deploy/crd/bases \
		output:webhook:dir=deploy/webhook
	cp deploy/crd/bases/vault.banzaicloud.com_vaults.yaml deploy/charts/vault-operator/crds/crd.yaml

.PHONY: gen-code
gen-code: ## Generate deepcopy, client, lister, and informer objects
	$(CONTROLLER_GEN) object:headerFile="hack/custom-boilerplate.go.txt" paths="./..."
	./hack/update-codegen.sh v${CODE_GENERATOR_VERSION}

.PHONY: gen-helm-docs
gen-helm-docs: ## Generate Helm chart documentation
	$(HELM_DOCS) -s file -c deploy/charts/ -t README.md.gotmpl

.PHONY: generate
generate: gen-manifests gen-code gen-helm-docs
generate: ## Generate manifests, code, and docs resources

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: gen-manifests ## Install CRDs into the K8s cluster
	$(KUSTOMIZE) build deploy/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: gen-manifests ## Uninstall CRDs from the K8s cluster. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build deploy/crd | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy
deploy: gen-manifests ## Deploy manager resources to the K8s cluster
	cd deploy/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build deploy/default | kubectl apply -f -

.PHONY: undeploy
clean: ## Clean manager resources from the K8s cluster. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build deploy/default | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

##@ Dependencies

# Dependency tool chain
GOLANGCI_VERSION = 1.53.3
LICENSEI_VERSION = 0.8.0
KIND_VERSION = 0.20.0
KURUN_VERSION = 0.7.0
CODE_GENERATOR_VERSION = 0.27.1
HELM_DOCS_VERSION = 1.11.0
KUSTOMIZE_VERSION = 5.1.0
CONTROLLER_TOOLS_VERSION = 0.12.0

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

KUSTOMIZE ?= $(or $(shell which kustomize),$(LOCALBIN)/kustomize)
$(KUSTOMIZE): $(LOCALBIN)
	@if test -x $(LOCALBIN)/kustomize && ! $(LOCALBIN)/kustomize version | grep -q v$(KUSTOMIZE_VERSION); then \
		echo "$(LOCALBIN)/kustomize version is not expected $(KUSTOMIZE_VERSION). Removing it before installing."; \
		rm -rf $(LOCALBIN)/kustomize; \
	fi
	test -s $(LOCALBIN)/kustomize || GOBIN=$(LOCALBIN) GO111MODULE=on go install sigs.k8s.io/kustomize/kustomize/v5@v$(KUSTOMIZE_VERSION)

CONTROLLER_GEN ?= $(or $(shell which controller-gen),$(LOCALBIN)/controller-gen)
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen && $(LOCALBIN)/controller-gen --version | grep -q v$(CONTROLLER_TOOLS_VERSION) || \
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@v$(CONTROLLER_TOOLS_VERSION)

ENVTEST ?= $(or $(shell which setup-envtest),$(LOCALBIN)/setup-envtest)
$(ENVTEST): $(LOCALBIN)
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

GOLANGCI_LINT ?= $(or $(shell which golangci-lint),$(LOCALBIN)/golangci-lint)
$(GOLANGCI_LINT): $(LOCALBIN)
	test -s $(LOCALBIN)/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- v${GOLANGCI_VERSION}

LICENSEI ?= $(or $(shell which licensei),$(LOCALBIN)/licensei)
$(LICENSEI): $(LOCALBIN)
	test -s $(LOCALBIN)/licensei || curl -sfL https://raw.githubusercontent.com/goph/licensei/master/install.sh | bash -s -- v${LICENSEI_VERSION}

HELM ?= $(or $(shell which helm),$(LOCALBIN)/helm)
$(HELM): $(LOCALBIN)
	test -s $(LOCALBIN)/helm || curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | USE_SUDO=false HELM_INSTALL_DIR=$(LOCALBIN) bash

KIND ?= $(or $(shell which kind),$(LOCALBIN)/kind)
$(KIND): $(LOCALBIN)
	@if [ ! -s "$(LOCALBIN)/kind" ]; then \
		curl -Lo $(LOCALBIN)/kind https://kind.sigs.k8s.io/dl/v${KIND_VERSION}/kind-$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m | sed -e "s/aarch64/arm64/; s/x86_64/amd64/"); \
		chmod +x $(LOCALBIN)/kind; \
	fi

HELM_DOCS ?= $(or $(shell which helm-docs),$(LOCALBIN)/helm-docs)
$(HELM_DOCS): $(LOCALBIN)
	@if [ ! -s "$(LOCALBIN)/helm-docs" ]; then \
		curl -L https://github.com/norwoodj/helm-docs/releases/download/v${HELM_DOCS_VERSION}/helm-docs_${HELM_DOCS_VERSION}_$(shell uname)_x86_64.tar.gz | tar -zOxf - helm-docs > ./bin/helm-docs; \
		chmod +x $(LOCALBIN)/helm-docs; \
	fi

KURUN ?= $(or $(shell which kurun),$(LOCALBIN)/kurun)
$(KURUN): $(LOCALBIN)
	@if [ ! -s "$(LOCALBIN)/kurun" ]; then \
		curl -Lo  $(LOCALBIN)/kurun https://github.com/banzaicloud/kurun/releases/download/${KURUN_VERSION}/kurun-$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m | sed -e "s/aarch64/arm64/; s/x86_64/amd64/"); \
		chmod +x  $(LOCALBIN)/kurun; \
	fi

# TODO: add support for hadolint and yamllint dependencies
HADOLINT ?= hadolint
YAMLLINT ?= yamllint

.PHONY: deps
deps: $(HELM) $(CONTROLLER_GEN) $(KUSTOMIZE) $(KIND)
deps: $(HELM_DOCS) $(ENVTEST) $(GOLANGCI_LINT) $(LICENSEI) $(KURUN)
deps: ## Download and install dependencies
