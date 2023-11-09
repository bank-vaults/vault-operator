# Vault Operator

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/bank-vaults/vault-operator/ci.yaml?style=flat-square)](https://github.com/bank-vaults/vault-operator/actions/workflows/ci.yaml)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/bank-vaults/vault-operator/badge?style=flat-square)](https://api.securityscorecards.dev/projects/github.com/bank-vaults/vault-operator)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/8056/badge)](https://www.bestpractices.dev/projects/8056)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/vault-operator)](https://artifacthub.io/packages/search?repo=vault-operator)

**Kubernetes operator for [Hashicorp Vault](https://www.vaultproject.io/).**

## Documentation

The official documentation for the operator is available at [https://bank-vaults.dev](https://bank-vaults.dev/docs/operator/).

## Version compatibility matrix

Please see [VERSIONS.md](VERSIONS.md) for version compatibility.

## Development

**For an optimal developer experience, it is recommended to install [Nix](https://nixos.org/download.html) and [direnv](https://direnv.net/docs/installation.html).**

_Alternatively, install [Go](https://go.dev/dl/) on your computer then run `make deps` to install the rest of the dependencies._

Make sure Docker is installed with Compose and Buildx.

Fetch required tools:

```shell
make deps
```

Run project dependencies:

```shell
make up
```

Run the operator:

```shell
make run
```

Run the test suite:

```shell
make test
make test-acceptance
```

Run linters:

```shell
make lint # pass -j option to run them in parallel
```

Some linter violations can automatically be fixed:

```shell
make fmt
```

Build artifacts locally:

```shell
make artifacts
```

Once you are done, you can tear down project dependencies:

```shell
make down
```

## License

The project is licensed under the [Apache 2.0 License](LICENSE).
