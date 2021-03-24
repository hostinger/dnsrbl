# HBL SaaS 

![Generic badge](https://img.shields.io/badge/V-1.0.0-green.svg) 

Hostinger Block List microservice based on Golang Echo framework. Hostinger Block List serves as a single source of truth for all blocked IP addresses in Hostinger organization.

# Table of contents
- [Docs](#Docs)
- [API](#API)
- [CLI](#CLI)
- [SDK](#SDK)
- [Development](#Development)
- [Contributing](#contributing)
- [License](#license)

# Docs
For API documentation we use Swagger. You can access Swagger documentation upon launch at http://127.0.0.1:8080/swagger/index.html

Upon API changes, a separate PR should be opened to update documentation. Documentation is automatically created with https://github.com/swaggo/swag
```
swag init -g cmd/hbl.go
```

# API
For API we use Golang Echo framework (https://echo.labstack.com/).

### Headers
Every request to API must have `Content-Type: application/json` set or the request will be ignored.

### Authentication
Every request, which has `KeyAuthMiddleware` function enabled expects `X-API-Key` request header with the API token. Requests without this or invalid token will return an `Unauthorized` response.

# CLI
There is a CLI application available, which helps interact with HBL API right from the terminal.

### Usage
```
./hblctl --help
Application which helps interact with Hostinger Block List API service.

Usage:
  hblctl [command]

Available Commands:
  allow
  block
  delete
  list
  sync

Flags:
      --config string           config file (default is $HOME/.hblctl.yaml)
      --hbl-api-host string     Host for connecting to the HBL API. (HBL_API_HOST)
      --hbl-api-key string      Key for connecting to the HBL API. (HBL_API_KEY)
      --hbl-api-port string     Port for connecting to the HBL API. (HBL_API_PORT)
      --hbl-api-scheme string   Scheme for connecting to the HBL API. (HBL_API_SCHEME)
  -h, --help                    help for hblctl

Use "hblctl [command] --help" for more information about a command.
```
### Block
```bash
./hblctl block <ip> <author> <comment> --hbl-api-host <api-host> --hbl-api-port <api-port> --hbl-api-scheme <api-scheme> --hbl-api-key <api-key>
```

### Allow
```bash
./hblctl allow <ip> <author> <comment> --hbl-api-host <api-host> --hbl-api-port <api-port> --hbl-api-scheme <api-scheme> --hbl-api-key <api-key>
```

### Delete
```bash
./hblctl delete <ip> --hbl-api-host <api-host> --hbl-api-port <api-port> --hbl-api-scheme <api-scheme> --hbl-api-key <api-key>
```

### List
```bash
./hblctl list [<ip>] --hbl-api-host <api-host> --hbl-api-port <api-port> --hbl-api-scheme <api-scheme> --hbl-api-key <api-key>
```

### Sync
```bash
./hblctl sync [<ip>] --hbl-api-host <api-host> --hbl-api-port <api-port> --hbl-api-scheme <api-scheme> --hbl-api-key <api-key>
```

# SDK
There is an official Golang SDK package available, which will help interact with HBL API through code.

Example usage:
```golang
import "github.com/hostinger/hbl/sdk"

func main() {
	c := sdk.NewClient("key", "url")
	if err := c.Allow(context.Background(), "127.0.0.1", "Author", "Comment"); err != nil {
		return err
	}
}
```
For more details on available functions see [SDK](https://github.com/hostinger/hbl/tree/master/sdk)

# Development
For local development we use Docker. Dockerfile expects an `.env` file to be created with credentials at the root directory. You can find this `.env` file inside Vault.

# Contributing
Pull requests are welcome. For major changes, issue describing the change needs to be opened before.

# License
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
