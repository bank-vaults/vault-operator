# Vault Operator

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/bank-vaults/vault-operator/ci.yaml?style=flat-square)](https://github.com/bank-vaults/vault-operator/actions/workflows/ci.yaml)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/bank-vaults/vault-operator/badge?style=flat-square)](https://api.securityscorecards.dev/projects/github.com/bank-vaults/vault-operator)

**Kubernetes operator for [Hashicorp Vault](https://www.vaultproject.io/).**

## Documentation

The official documentation for the operator is available at [https://bank-vaults.dev](https://bank-vaults.dev/docs/operator/).

## Development

**For an optimal developer experience, it is recommended to install [Nix](https://nixos.org/download.html) and [direnv](https://direnv.net/docs/installation.html).**

_Alternatively, install [Go](https://go.dev/dl/) on your computer then run `make deps` to install the rest of the dependencies._

Make sure Docker is installed with Compose and Buildx.

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
