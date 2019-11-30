# az-blob-hashdeep [![Build Status](https://travis-ci.org/evenh/az-blob-hashdeep.svg?branch=master)](https://travis-ci.org/evenh/az-blob-hashdeep) [![Go Report Card](https://goreportcard.com/badge/github.com/evenh/az-blob-hashdeep)](https://goreportcard.com/report/github.com/evenh/az-blob-hashdeep)

A simple tool for generating [hashdeep](https://github.com/jessek/hashdeep) compatible output for an Azure Blob Storage container. Useful for verifying migration of data to/from Azure.

This implementation requires that the `Content-MD5` is set for all the blobs in a container. It will **not** verify the hashes, it will only be used for the final output. 

## How to use?

Get precompiled binaries from the [releases page](https://github.com/evenh/az-blob-hashdeep/releases) or use the [Docker image](https://hub.docker.com/repository/docker/evenh/az-blob-hashdeep): `evenh/az-blob-hashdeep`.

```bash
export AZURE_ACCOUNT_NAME=myaccount
export AZURE_ACCOUNT_KEY=secretKey
export AZURE_CONTAINER=migrationcontainer

./az-blob-hashdeep generate --account-name=$AZURE_ACCOUNT_NAME \
                            --account-key=$AZURE_ACCOUNT_KEY \
                            --container=$AZURE_CONTAINER \
                            --output ~/$AZURE_ACCOUNT_NAME-$AZURE_CONTAINER.hashdeep
```

This will result in a output like:

```
%%%% HASHDEEP-1.0
%%%% size,md5,filename
## Invoked from: /Users/evenh/dev/evenh/az-blob-hashdeep
## $ ./az-blob-hashdeep generate --account-name=myaccount --account-key=secretKey --container=migrationcontainer --output /Users/evenh/myaccount-migrationcontainer.hashdeep
##
1026764,ddb5d9fb991f62be9c55383aefa8e8e3,00/00/000008af-2e78-4b21-9a0e-a44ee77d4606
97428,4fdb49a5de56a1b11c9c37264a1bb927,00/00/00006c79-1c38-45f8-a3b8-ebb299fc67a1
[more entries…]
```

### Use prefix
There is also a potential useful flag: `--prefix (-p)` which will prepend a prefix to file paths. Example:
```
Path: foo/bar/file.txt
Prefix: old-fs-01
Outputted path: old-fs-01/foo/bar/file.txt
```

## Contributing

Contributions in form of issues and pull requests are most welcome.

## License

This project is licensed under the Apache 2.0 License. See [LICENSE](./LICENSE) for more information.
