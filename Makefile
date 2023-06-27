# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

export PATH := $(abspath bin/):${PATH}

CONTAINER_IMAGE_REF = ghcr.io/bank-vaults/vault-operator:dev

# Dependency versions
GOLANGCI_VERSION = 1.53.3
LICENSEI_VERSION = 0.8.0
KIND_VERSION = 0.20.0
KURUN_VERSION = 0.7.0
CODE_GENERATOR_VERSION = 0.27.1
CONTROLLER_GEN_VERSION = 0.12.0
HELM_DOCS_VERSION = 1.11.0

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

.PHONY: build
build: ## Build binary
	@mkdir -p build
	go build -race -o build/manager ./cmd/manager

.PHONY: run
run: ## Run the operator locally talking to a Kubernetes cluster
	kubectl replace -f deploy/examples/default/crd.yaml || kubectl create -f deploy/examples/default/crd.yaml
	kubectl apply -f deploy/examples/default/rbac.yaml
	OPERATOR_NAME=vault-dev go run cmd/manager/main.go -verbose

.PHONY: clean
clean: ## Clean operator resources from a Kubernetes cluster
	kubectl delete -f deploy/examples/default/crd.yaml
	kubectl delete -f deploy/examples/default/rbac.yaml

.PHONY: artifacts
artifacts: container-image helm-chart
artifacts: ## Build artifacts

.PHONY: container-image
container-image: ## Build container image
	docker build -t ${CONTAINER_IMAGE_REF} .

.PHONY: helm-chart
helm-chart: ## Build Helm chart
	@mkdir -p build
	helm package -d build/ deploy/charts/vault-operator

.PHONY: check
check: test lint ## Run checks (tests and linters)

.PHONY: test
test: ## Run tests
	go test -race -v ./...

.PHONY: test-acceptance
test-acceptance: ## Run acceptance tests
	go test -race -v -timeout 900s -tags kubeall ./test

.PHONY: lint
lint: lint-go lint-helm lint-docker lint-yaml
lint: ## Run linters

.PHONY: lint-go
lint-go:
	golangci-lint run $(if ${CI},--out-format github-actions,)

.PHONY: lint-helm
lint-helm:
	helm lint deploy/charts/vault-operator

.PHONY: lint-docker
lint-docker:
	hadolint Dockerfile

.PHONY: lint-yaml
lint-yaml:
	yamllint $(if ${CI},-f github,) --no-warnings .

.PHONY: fmt
fmt: ## Format code
	golangci-lint run --fix

.PHONY: license-check
license-check: ## Run license check
	licensei check
	licensei header

.PHONY: generate
generate: generate-code generate-crds generate-helm-docs
generate: ## Run generation jobs

.PHONY: generate-code
generate-code: ## Regenerate clientset, deepcopy funcs, listers and informers
	./hack/scripts/update-codegen.sh v${CODE_GENERATOR_VERSION}

.PHONY: generate-crds
generate-crds: ## Regenerate CRDs in the Helm chart and examples
	controller-gen crd:maxDescLen=0 paths=./pkg/... output:crd:artifacts:config=./deploy/examples/default
	cp deploy/examples/default/vault.banzaicloud.com_vaults.yaml deploy/charts/vault-operator/crds/crd.yaml
	cp deploy/examples/default/vault.banzaicloud.com_vaults.yaml deploy/examples/default/crd.yaml
	rm deploy/examples/default/vault.banzaicloud.com_vaults.yaml

.PHONY: generate-helm-docs
generate-helm-docs:
	helm-docs -s file -c deploy/charts/ -t README.md.gotmpl

deps: bin/golangci-lint bin/licensei bin/kind bin/kurun bin/controller-gen bin/helm-docs
deps: ## Install dependencies

bin/golangci-lint:
	@mkdir -p bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- v${GOLANGCI_VERSION}

bin/licensei:
	@mkdir -p bin
	curl -sfL https://raw.githubusercontent.com/goph/licensei/master/install.sh | bash -s -- v${LICENSEI_VERSION}

bin/kind:
	@mkdir -p bin
	curl -Lo bin/kind https://kind.sigs.k8s.io/dl/v${KIND_VERSION}/kind-$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m | sed -e "s/aarch64/arm64/; s/x86_64/amd64/")
	@chmod +x bin/kind

bin/kurun:
	@mkdir -p bin
	curl -Lo bin/kurun https://github.com/banzaicloud/kurun/releases/download/${KURUN_VERSION}/kurun-$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m | sed -e "s/aarch64/arm64/; s/x86_64/amd64/")
	@chmod +x bin/kurun

bin/controller-gen:
	@mkdir -p bin
	set -ex ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	GOBIN=$(PWD)/bin/controller-gen-${CONTROLLER_GEN_VERSION}-bin/ go install sigs.k8s.io/controller-tools/cmd/controller-gen@v${CONTROLLER_GEN_VERSION} ;\
	mv $(PWD)/bin/controller-gen-${CONTROLLER_GEN_VERSION}-bin/controller-gen $(PWD)/bin/controller-gen ;\
	rm -rf $(PWD)/bin/controller-gen-${CONTROLLER_GEN_VERSION}-bin ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR

bin/helm-docs:
	@mkdir -p bin
	curl -L https://github.com/norwoodj/helm-docs/releases/download/v${HELM_DOCS_VERSION}/helm-docs_${HELM_DOCS_VERSION}_$(shell uname)_x86_64.tar.gz | tar -zOxf - helm-docs > ./bin/helm-docs
	@chmod +x bin/helm-docs

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'
