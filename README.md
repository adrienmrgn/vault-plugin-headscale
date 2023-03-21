# vault-plugin-headscale
Secret engine to create/remove users and generate preauthkey from a Headscale control plane from Hashicorp Vault

## Build

The plugin is a Go binary, built using the repo's Makefile. Binary is generated in `\bin`

```shell
make build
```

## Installation
Follow Vault documentation to enable this plugin on your Vault cluster.

For tests purpose, the plugin is built and added to a developpment Docker image (see [docker file](.docker/Dockerfile.vault))
A [docker-compose file](./docker-compose.yaml) is provided to run a `headscale` server and a `vault` server. The plugin is loaded as a development plugin in the test scenario.

The Makefile target `make compose` :
* build the plugin
* build the vault image with the plugin inside
* spin the containers up
* enable the Headscale plugin

## Usage

Once the plugin leaded by Vault and enable at `/headscale`, here's how to configure and use it.

### Generate Headscale access key
```shell
export HEADSCALE_API_KEY=$(docker exec headscale headscale apikey create -o yaml)
```
### Configure secret engine
```shell 
export VAULT_ADDR=http://127.0.0.1:8200
export VAULT_TOKEN=root
export HEADSCALE_API_URL="http://headscale:8080"
vault write headscale/config/access api_key="${HEADSCALE_API_KEY}" api_url="${HEADSCALE_API_URL}"
vault read headscale/config/access
```

### Create a user
```shell
vault write headscale/user/ name=foo
```

### Get a user
```shell
vault read headscale/user/foo 
```

### Generate a key
```shell
vault read headscale/creds/foo \
 ephemeral=true \
 reusable=true \
 tags=hello,world \
 expiration=2024-01-01T00:00:00Z
```
