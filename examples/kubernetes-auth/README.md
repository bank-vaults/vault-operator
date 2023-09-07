# Kubernetes auth

This example demonstrates how to authenticate against Vault using a Kubernetes Service Account token.
It also contains an example for writing policies targeting the accessor (identity behind the token).

## Check access

Run the following command to open a shell in the `default` namespace:

```shell
kubectl run vault-shell --rm -i --tty --env "VAULT_ADDR=http://vault:8200" --image hashicorp/vault:1.14.1 -- sh
```

Exchange the Service Account token for a Vault token:

```shell
export VAULT_TOKEN=$(vault write -field=token auth/kubernetes/login role=default jwt=$(cat /run/secrets/kubernetes.io/serviceaccount/token))
```

## Simple policy

Write some data into Vault:

```shell
vault kv put shared/foo bar=baz
```

Read the secret back:

```shell
vault kv get shared/foo
```

## Accessor-based policy

The second policy limits access to `namespaces/NAMESPACE` (namespace in this case is `default`, because that's where the shell pod is running).

Try writing some data into Vault:

```shell
vault kv put namespaces/not-default/foo bar=baz
```

You should see a permission denied error:

```
Error writing data to namespaces/data/vault/foo: Error making API request.

URL: PUT http://vault:8200/v1/namespaces/data/vault/foo
Code: 403. Errors:

* 1 error occurred:
        * permission denied
```

Try writing some data into the right path:

```shell
vault kv put namespaces/default/foo bar=baz
```

Read the data back:

```shell
vault kv get namespaces/default/foo
```

Try reading from another path:

```shell
vault kv get namespaces/not-default/foo
```

It should fail again:

```
Error reading namespaces/data/not-default/foo: Error making API request.

URL: GET http://vault:8200/v1/namespaces/data/not-default/foo
Code: 403. Errors:

* 1 error occurred:
        * permission denied
```
