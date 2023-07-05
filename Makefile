
# Image URL to use all building/pushing image targets
IMG ?= ghcr.io/bank-vaults/vault-operator:latest

# Tool chain
GOLANGCI_VERSION = 1.53.3
LICENSEI_VERSION = 0.8.0
KIND_VERSION = 0.20.0
KURUN_VERSION = 0.7.0
CODE_GENERATOR_VERSION = 0.27.1
CONTROLLER_GEN_VERSION = 0.12.0
HELM_DOCS_VERSION = 1.11.0

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.27.1

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	$(GOLANGCI_LINT) run --fix

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: license-check
license-check: ## Run license check
	licensei check
	licensei header

.PHONY: lint-go
lint-go:
	$(GOLANGCI_LINT) run $(if ${CI},--out-format github-actions,)

# TODO: add helm
.PHONY: lint-helm
lint-helm:
	helm lint deploy/charts/vault-operator

# TODO: add hadolint
.PHONY: lint-docker
lint-docker:
	hadolint Dockerfile

# TODO: add yamllint
.PHONY: lint-yaml
lint-yaml:
	yamllint $(if ${CI},-f github,) --no-warnings .

.PHONY: lint
lint: lint-go lint-helm lint-docker lint-yaml
lint: ## Run linters

.PHONY: test
test: generate fmt vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" \
		go test -race -v ./... -coverprofile cover.out

.PHONY: test-acceptance
test-acceptance: generate fmt vet envtest ## Run acceptance tests
	go test -race -v -timeout 900s -tags kubeall ./test

.PHONY: check
check: test lint ## Run checks (tests and linters)

##@ Autogeneration

.PHONY: generate-manifests
generate-manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." \
		output:rbac:artifacts:config=deploy/rbac \
		output:crd:artifacts:config=deploy/crd/bases
	cp deploy/crd/bases/bank-vaults.dev_vaults.yaml deploy/charts/vault-operator/crds/crd.yaml

.PHONY: generate-code
generate-code: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/custom-boilerplate.go.txt" paths="./..."
	./hack/update-codegen.sh v${CODE_GENERATOR_VERSION}

.PHONY: generate-helm-docs
generate-helm-docs:
	$(HELM_DOCS) -s file -c deploy/charts/ -t README.md.gotmpl

.PHONY: generate
generate: generate-manifests generate-code generate-helm-docs
generate: ## Run generation jobs

##@ Build

.PHONY: build
build: generate-manifests generate-code fmt vet ## Build manager binary.
	@mkdir -p build
	go build -race -o build/manager ./cmd/manager

.PHONY: run
run: generate-manifests generate-code fmt vet deploy ## Run a controller from your host.
	OPERATOR_NAME=vault-dev go run cmd/manager/main.go -verbose

# If you wish built the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64 ). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build: test ## Build docker image with the manager.
	$(CONTAINER_TOOL) build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	$(CONTAINER_TOOL) push ${IMG}

# PLATFORMS defines the target platforms for  the manager image be build to provide support to multiple
# architectures. (i.e. make docker-buildx IMG=myregistry/mypoperator:0.0.1). To use this option you need to:
# - able to use docker buildx . More info: https://docs.docker.com/build/buildx/
# - have enable BuildKit, More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image for your registry (i.e. if you do not inform a valid value via IMG=<myregistry/image:<tag>> then the export will fail)
# To properly provided solutions that supports more than one platform you should use this option.
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: docker-buildx
docker-buildx: test ## Build and push docker image for the manager for cross-platform support
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- $(CONTAINER_TOOL) buildx create --name project-v3-builder
	$(CONTAINER_TOOL) buildx use project-v3-builder
	- $(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --tag ${IMG} -f Dockerfile.cross .
	- $(CONTAINER_TOOL) buildx rm project-v3-builder
	rm Dockerfile.cross

.PHONY: helm-chart
helm-chart: ## Build Helm chart
	@mkdir -p build
	helm package -d build/ deploy/charts/vault-operator

.PHONY: build-artifacts
build-artifacts: docker-build helm-chart
build-artifacts: ## Build artifacts

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: up
up: ## Start development environment
	kind create cluster

.PHONY: stop
stop: ## Stop development environment
	# TODO: consider using k3d instead
	kind delete cluster

.PHONY: down
down: ## Destroy development environment
	kind delete cluster

.PHONY: clean
clean: undeploy ## Clean operator resources from a Kubernetes cluster

.PHONY: install
install: generate-manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build deploy/crd | $(KUBECTL) apply -f -

.PHONY: uninstall
uninstall: generate-manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build deploy/crd | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy
deploy: generate-manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd deploy/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build deploy/default | $(KUBECTL) apply -f -

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build deploy/default | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUBECTL ?= kubectl
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest

## Tool Versions
KUSTOMIZE_VERSION ?= v5.0.1
CONTROLLER_TOOLS_VERSION ?= v0.12.0

.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary. If wrong version is installed, it will be removed before downloading.
$(KUSTOMIZE): $(LOCALBIN)
	@if test -x $(LOCALBIN)/kustomize && ! $(LOCALBIN)/kustomize version | grep -q $(KUSTOMIZE_VERSION); then \
		echo "$(LOCALBIN)/kustomize version is not expected $(KUSTOMIZE_VERSION). Removing it before installing."; \
		rm -rf $(LOCALBIN)/kustomize; \
	fi
	test -s $(LOCALBIN)/kustomize || GOBIN=$(LOCALBIN) GO111MODULE=on go install sigs.k8s.io/kustomize/kustomize/v5@$(KUSTOMIZE_VERSION)

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary. If wrong version is installed, it will be overwritten.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen && $(LOCALBIN)/controller-gen --version | grep -q $(CONTROLLER_TOOLS_VERSION) || \
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

deps: $(GOLANGCI_LINT) $(LICENSEI) $(KIND) $(KURUN) $(CONTROLLER_GEN) $(KUSTOMIZE) $(ENVTEST) $(HELM_DOCS)
deps: ## Install dependencies

GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint
$(GOLANGCI_LINT): $(LOCALBIN)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- v${GOLANGCI_VERSION}

LICENSEI ?= $(LOCALBIN)/licensei
$(LICENSEI): $(LOCALBIN)
	curl -sfL https://raw.githubusercontent.com/goph/licensei/master/install.sh | bash -s -- v${LICENSEI_VERSION}

KIND ?= $(LOCALBIN)/kind
$(KIND): $(LOCALBIN)
	curl -Lo $(LOCALBIN)/kind https://kind.sigs.k8s.io/dl/v${KIND_VERSION}/kind-$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m | sed -e "s/aarch64/arm64/; s/x86_64/amd64/")
	@chmod +x $(LOCALBIN)/kind

KURUN ?= $(LOCALBIN)/kurun
$(KURUN): $(LOCALBIN)
	curl -Lo $(LOCALBIN)/kurun https://github.com/banzaicloud/kurun/releases/download/${KURUN_VERSION}/kurun-$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m | sed -e "s/aarch64/arm64/; s/x86_64/amd64/")
	@chmod +x $(LOCALBIN)/kurun

HELM_DOCS ?= $(LOCALBIN)/helm-docs
$(HELM_DOCS): $(LOCALBIN)
	curl -L https://github.com/norwoodj/helm-docs/releases/download/v${HELM_DOCS_VERSION}/helm-docs_${HELM_DOCS_VERSION}_$(shell uname)_x86_64.tar.gz | tar -zOxf - helm-docs > ./bin/helm-docs
	@chmod +x $(LOCALBIN)/helm-docs
