# Vault Operator examples

These examples demonstrate different Vault Operator features.

## Examples

- [Base](base/): Minimal example for launching Vault
- [File storage](file/): Persistent volumes and file storage
- [Startup secrets](startup-secrets/): Initialize Vault with secrets
- [Kubernetes auth](kubernetes-auth/): Authenticate against Vault using Kubernetes Service Accounts

## Prerequisites

- Ability to setup a Kubernetes cluster (eg. using [KinD](https://kind.sigs.k8s.io/))
- kubectl
- kustomize
- [Helm](https://helm.sh/)
- [vault CLI](https://developer.hashicorp.com/vault/downloads)
- kubectl [view-secret plugin](https://github.com/elsesiy/kubectl-view-secret) _(optional)_

It is recommended that you check out this repository and run examples from there as some examples require additional steps (but you may be able to run some of them directly):

```shell
git clone git@github.com:bank-vaults/vault-operator.git
cd vault-operator/examples
```

## Set up a Kubernetes cluster

Vault Operator should run on recent Kubernetes versions.
We generally use [KinD](https://kind.sigs.k8s.io) for demos, but [k3d](https://k3d.io) is also a popular option.

You can launch a local cluster using KinD by running the following command:

```shell
kind create cluster
```

## Install Vault Operator

Install the latest version of the Vault Operator:

```shell
helm upgrade --install --wait --namespace vault-system --create-namespace vault-operator oci://ghcr.io/bank-vaults/helm-charts/vault-operator
```

## Install Vault

Choose one of the examples in this folder, follow instructions (if any) in the README and install the example:

```shell
EXAMPLE=base

kustomize build $EXAMPLE | kubectl apply -f -
```

Or if you haven't checked out the repository:

```shell
EXAMPLE=base

kustomize build github.com/bank-vaults/vault-operator/examples/$EXAMPLE | kubectl apply -f -
```

> [!IMPORTANT]
> Examples are generally mutually exclusive, so make sure you don't install two of them on the same cluster at the same time.
> If you want to experiment with multiple examples, consider using [vCluster](https://www.vcluster.com/).

## Checking Vault

After installation, you may want to check the installation and interact with Vault.

First, wait for Vault to become ready:

```shell
kubectl wait pods vault-0 --for condition=Ready --timeout=120s
```

> [!NOTE]
> You may need to pass a namespace parameter to the above and following commands: `--namespace YOUR_NAMESPACE`

Set the Vault token from the Kubernetes secret:

```shell
export VAULT_TOKEN=$(kubectl get secrets vault-unseal-keys -o jsonpath={.data.vault-root} | base64 --decode)
```

Tell the CLI where Vault is listening:

```shell
export VAULT_ADDR=http://127.0.0.1:8200
```

> [!NOTE]
> If you are running an example with TLS configured, the address should be `https://127.0.0.1:8200`.

Port forward to the Vault service:

```shell
kubectl port-forward service/vault 8200 1>/dev/null &
```

Check Vault status:

```shell
vault status
```

Open the UI (and login with the root token):

```shell
open $VAULT_ADDR
```

The same commands again:

```shell
kubectl wait pods vault-0 --for condition=Ready --timeout=120s
export VAULT_TOKEN=$(kubectl get secrets vault-unseal-keys -o jsonpath={.data.vault-root} | base64 --decode)
export VAULT_ADDR=http://127.0.0.1:8200
kubectl port-forward service/vault 8200 1>/dev/null &
vault status
```

TODO: Write a script with the above commands?

## Cleanup

Kill background jobs:

```shell
kill %1
```

Delete the installed example:

```shell
kustomize build $EXAMPLE | kubectl delete -f -
```

Delete unseal keys:

```shell
kubectl delete secret vault-unseal-keys
```

Delete the Operator:

```shell
helm -n vault-system delete vault-operator
```

Delete the cluster:

```shell
kind delete cluster
```
