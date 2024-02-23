#!/bin/bash

set -xeo pipefail

# REQUIREMENTS:
# - kustomize
# - kubectl
# - helm3
# - https://github.com/subfuzion/envtpl
# - https://github.com/hankjacobs/cidr
# - jq
# - kind
# - vault
#

METALLB_VERSION=v0.13.10
VAULT_VERSION=1.14.8
BANK_VAULTS_VERSION=v1.31.0
VAULT_TOKEN=$(uuidgen)
export VAULT_TOKEN

if [ $# = 0 ]; then
    echo "The Bank-Vaults Multi-DC CLI"
    echo
    echo "Usage:"
    echo "  $0 [command]"
    echo
    echo "Available Commands:"
    echo "  install    Installs a Vault cluster to one or more Kubernetes clusters"
    echo "  uninstall  Uninstalls a Vault cluster from one or more Kubernetes clusters"
    echo "  status     Displays the status a cluster to one or more Kubernetes clusters"
    exit 0
fi

function waitfor {
    WAIT_MAX=0
    until "$@" &>/dev/null || [ $WAIT_MAX -eq 45 ]; do
        sleep 1
        ((WAIT_MAX = WAIT_MAX + 1))
    done
}

function metallb_setup {
    export METALLB_ADDRESS_RANGE=$1
    # kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/${METALLB_VERSION}/manifests/namespace.yaml
    kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/${METALLB_VERSION}/config/manifests/metallb-native.yaml
    kubectl wait --namespace metallb-system --for condition=Available=true deploy --selector=app=metallb --timeout=90s
    kubectl wait --namespace metallb-system --for=condition=ready pod --selector=app=metallb --timeout=90s
    envtpl deploy/dev/multi-dc/test/metallb-config.yaml | kubectl apply -f -
}

function cidr_range {
    local cidr=$1
    cidr "${cidr}" | tr -d ' '
}

function node_setup {
    local instance=$1
    local lb_subnet=$2

    kind create cluster --name "${instance}"
    metallb_setup "$(cidr_range "${lb_subnet}")"
}

function infra_setup {
    # get the kind Docker network subnet
    # SUBNET=$(docker network inspect kind --format '{{(index .IPAM.Config 0).Subnet}}')

    node_setup primary 172.18.1.255/25

    node_setup secondary 172.18.2.255/25

    node_setup tertiary 172.18.3.255/25

    docker run -d --rm --network kind -e VAULT_DEV_ROOT_TOKEN_ID="${VAULT_TOKEN}" --name central-vault hashicorp/vault:"${VAULT_VERSION}"
    CENTRAL_VAULT_ADDRESS=$(docker inspect central-vault --format '{{.NetworkSettings.Networks.kind.IPAddress}}')
    export CENTRAL_VAULT_ADDRESS
}

function install_instance {
    local INSTANCE=$1

    kind load image-archive docker.tar --name "${INSTANCE}"

    helm upgrade --install vault-operator ./deploy/charts/vault-operator --wait --set image.tag=${OPERATOR_VERSION} --set image.pullPolicy=Never --set bankVaults.image.tag=${BANK_VAULTS_VERSION}

    kubectl apply -k deploy/rbac/
    envtpl deploy/dev/multi-dc/test/cr-"${INSTANCE}".yaml | kubectl apply -f -

    echo "Waiting for for ${INSTANCE} vault instance..."
    waitfor kubectl get pod/vault-"${INSTANCE}"-0

    kubectl wait --for=condition=ready pod/vault-"${INSTANCE}"-0 --timeout=180s
}

COMMAND=$1

if [ "$COMMAND" = "install" ]; then

    infra_setup

    CENTRAL_VAULT_ADDRESS=$(docker inspect central-vault --format '{{.NetworkSettings.Networks.kind.IPAddress}}')
    export CENTRAL_VAULT_ADDRESS

    ## Primary
    kubectl config use-context kind-primary

    install_instance primary

    RAFT_LEADER_ADDRESS=$(kubectl get service vault-primary -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    export RAFT_LEADER_ADDRESS

    kubectl get secrets vault-primary-tls -o json | jq 'del(.metadata.ownerReferences)' | jq 'del(.metadata.resourceVersion)' | jq 'del(.metadata.uid)' >vault-primary-tls.json

    ## Secondary
    kubectl config use-context kind-secondary

    kubectl apply -f vault-primary-tls.json

    install_instance secondary

    ## Tertiary
    kubectl config use-context kind-tertiary

    kubectl apply -f vault-primary-tls.json

    install_instance tertiary

    ## Cleanup

    rm vault-primary-tls.json

    echo -e "\nMulti-DC Vault cluster setup completed."

elif [ "$COMMAND" = "status" ]; then

    VAULT_SKIP_VERIFY="true"
    export VAULT_SKIP_VERIFY

    VAULT_ADDR=https://$(implement_me):8200
    export VAULT_ADDR

    vault operator raft list-peers -format json | jq

elif [ "$COMMAND" = "uninstall" ]; then

    kind delete cluster --name primary
    kind delete cluster --name secondary
    kind delete cluster --name tertiary
    docker rm -f central-vault
else

    echo "unknown command: $COMMAND"

fi
