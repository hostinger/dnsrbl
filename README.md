# HBL SaaS 

![Generic badge](https://img.shields.io/badge/V-1.0.0-green.svg) 

Hostinger Block List microservice based on Golang Echo framework. Hostinger Block List serves as a single source of truth for all blocked IP addresses in Hostinger organization.

# Table of contents
- [CLI](#CLI)
- [Contributing](#contributing)
- [License](#license)

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

# Contributing
Pull requests are welcome. For major changes, issue describing the change needs to be opened before.

# License
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
