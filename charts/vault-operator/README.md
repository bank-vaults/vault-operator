# Vault Operator

![kube version: >=1.22.0-0](https://img.shields.io/badge/kube%20version->=1.22.0--0-informational?style=flat-square)

Kubernetes operator for [Hashicorp Vault](https://www.vaultproject.io/).

**Homepage:** <https://bank-vaults.dev/>

## TL;DR;

```bash
helm install --generate-name --wait oci://ghcr.io/bank-vaults/helm-charts/vault-operator
```

## Values

| Parameter                     | Description                                                     | Default                              |
|-------------------------------|-----------------------------------------------------------------|--------------------------------------|
| `image.pullPolicy`            | Container pull policy                                           | `IfNotPresent`                       |
| `image.repository`            | Container image to use                                          | `ghcr.io/bank-vaults/vault-operator` |
| `image.tag`                   | Container image tag to deploy operator in                       | `.Chart.AppVersion`                  |
| `image.imagePullSecrets`      | Image pull secrets for private repositories                     | `[]`                                 |
| `image.bankVaultsRepository`  | Container image to use for Bank-Vaults (deprecated)             |                                      |
| `image.bankVaultsTag`         | Container image tag to use for Bank-Vaults (deprecated)         |                                      |
| `bankVaults.image.repository` | Bank-Vaults image repository                                    | `ghcr.io/banzaicloud/bank-vaults`    |
| `bankVaults.image.tag`        | Bank-Vaults image tag (pinned to supported Bank-Vaults version) | `1.19.0`                             |
| `replicaCount`                | k8s replicas                                                    | `1`                                  |
| `resources.requests.cpu`      | Container requested CPU                                         | `100m`                               |
| `resources.requests.memory`   | Container requested memory                                      | `128Mi`                              |
| `resources.limits.cpu`        | Container CPU limit                                             | `100m`                               |
| `resources.limits.memory`     | Container memory limit                                          | `256Mi`                              |
| `crdAnnotations`              | Annotations for the Vault CRD                                   | `{}`                                 |
| `securityContext`             | Container security context for vault-operator deployment        | `{}`                                 |
| `podSecurityContext`          | Pod security context for vault-operator deployment              | `{}`                                 |
| `psp.enabled`                 | Deploy PSP resources                                            | `false`                              |
| `psp.vaultSA`                 | Used service account for vault                                  | `vault`                              |

## Credits

Thanks to Cosmin Cojocar for the original Vault Operator Helm chart!
