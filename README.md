# tondns

[![Go Report Card](https://goreportcard.com/badge/github.com/slashfast/tondns)](https://goreportcard.com/report/github.com/slashfast/tondns)
[![License](https://img.shields.io/github/license/slashfast/tondns)](LICENSE) [![Static Badge](https://img.shields.io/badge/Based_on_TON-ffffff?logo=ton)](https://ton.org)
 ![Docker Image Size (tag)](https://img.shields.io/docker/image-size/slashfast/tondns/latest)



CLI tool for checking and auto-renewing TON DNS domains.

## Features

- 🔍 Domain Status Check
  - Expiration date tracking
  - Ownership verification
  - Domain resolution status
- 🔄 Smart Auto-renewal
  - Configurable renewal thresholds
  - Automatic renewal payments
  - Bulk domain processing
- ⚡ Performance
  - Concurrent domain checking
  - Efficient API usage
  - Support for proxy connections

## Installation

### From Source

```sh
git clone https://github.com/slashfast/tondns.git
cd tondns
go build
# or install to $GOPATH/bin or $HOME/go/bin
go install
```

### Using Docker

```sh
docker pull slashfast/tondns
# or build locally
docker build -t tondns .
```

## Usage

The tool has two main modes of operation:

### Full Mode (Seed Required)
Check all domains in your wallet using a seed:
```sh
tondns check --pretty --seed "your seed phrase" 
```


### Lite Mode (No Wallet Required)
Check the status of specified domains without a wallet connection:

```sh
# Check domain status
tondns check durov.ton gems.ton --lite --pretty
```

## Docker Usage

```sh
docker run tondns check --seed "your seed phrase"
```

## Examples

```sh
# Pass the seed from an environment variable and check all domains
SEED="your seed phrase" tondns check --pretty

# Renew domains if they expire soon
tondns check --renew --pretty --seed $SEED

# Force renew domains (regardless of expiration)
tondns check --renew --force --pretty --seed $SEED

# Check the status of specific domains
tondns check durov.ton gems.ton --lite --pretty

# Check domains using a SOCKS5 proxy (affects only API requests to Getgems)
tondns check --pretty --seed "your seed phrase" --proxy socks5://127.0.0.1:1080
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request


## Data Source
tondns uses the [Getgems](https://getgems.io/) API to retrieve up-to-date information about TON DNS domains, including current ownership status and expiration dates.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
