# dnsimple-cli

Update DNS records from the command-line via DNSimple

`dnsimple-cli` is single self-contained command-line utility for updating subdomain records that are managed by DNSimple written go that works on Linux, Mac OS and Windows alike.

[![Build Status](https://travis-ci.org/andreaskoch/dnsimple-cli.svg?branch=master)](https://travis-ci.org/andreaskoch/dnsimple-cli)

## Usage

```bash
dnsimple-cli <action> [arguments ...]
```

Get help:

```bash
dnsimple-cli --help
```

### Action: `login`

Save DNSimple API credentials to disc.

**Arguments**:

- `-email`: The e-mail address of your DNSimple account
- `-apitoken`: The DNSimple API token

**Example**:

```bash
dnsimple-cli login -email apiuser@example.com -apitoken TracsiflOgympacKoFieC
```

The credentials are saved to: `~/.dnsimple-cli/credentials.json`

### Action: `logout`

Remove any stored DNSimple API credentials from disc.

```bash
dnsimple-cli logout
```

### Action: `list`

List all available domains or subdomains.

**Arguments**

- `-domain`: A domain name (optional)
- `-subdomain`: A subdomain name (optional)

**Example**

List all available domains:

```bash
dnsimple-cli list
```

List all subdomain of a given domain:

```bash
dnsimple-cli list -domain example.com
```

List all DNS records for a given subdomain:

```bash
dnsimple-cli list -domain example.com -subdomain www
```

### Action: `create`

Create an address DNS record.

**Arguments**:

- `-domain`: A domain name (required)
- `-subdomain`: The subdomain name (required)
- `-ip`: An IPv4 or IPv6 address (required)
- `-ttl`: The time to live (TTL) for the DNS record in seconds (default: 600)

**Example**:

Create an `AAAA` record for the subdomain `www`:

```bash
dnsimple-cli create -domain example.com -subdomain www -ip 2001:0db8:0000:0042:0000:8a2e:0370:7334
```

Create an `AAAA` record for the subdomain `www` with TTL of 1 minute:

```bash
dnsimple-cli create -domain example.com -subdomain www -ip 2001:0db8:0000:0042:0000:8a2e:0370:7334 -ttl 60
```

Create an `A` record for the subdomain `www`:

```bash
dnsimple-cli create -domain example.com -subdomain www -ip 10.2.1.3
```

The `-ip` parameter can also be passed via Stdin:

```bash
echo "2001:0db8:0000:0042:0000:8a2e:0370:7334" | dnsimple-cli create -domain example.com -subdomain www -ttl 3600
```

### Action: `update`

Update the DNS record for a given sub domain

**Arguments**:

- `-domain`: A domain name (e.g. `example.com`)
- `-subdomain`: A subdomain name (e.g. `www`)
- `-ip`: An IPv4 or IPv6 address

**Example**:

Set the `AAAA` record of `www.example.com` to the given IP address:

```bash
dnsimple-cli update -domain example.com -subdomain www -ip 2001:0db8:0000:0042:0000:8a2e:0370:7334
```

The `-ip` address parameter can also be passed in via Stdin:

```bash
echo "2001:0db8:0000:0042:0000:8a2e:0370:7334" | dnsimple-cli update -domain example.com -subdomain www
```

## Dependencies

dnsimple-cli uses the following third-party libraries:

- The [github.com/pearkes/dnsimple](https://github.com/pearkes/dnsimple) API library for communicating with the DNSimple API
- [github.com/mitchellh/go-homedir](https://github.com/mitchellh/go-homedir) for determining the user's home directory path
- [github.com/spf13/afero](https://github.com/spf13/afero) as a filesystem abstraction for testing

## Installation & Build

Compile the dnsimple-cli binary for your current platform:

```bash
git clone git@github.com:andreaskoch/dnsimple-cli.git && cd dnsimple-cli

go run make.go -install
```

Compile the dnsimple-cli binaries for Linux (64bit, ARM, ARM5, ARM6, ARM7), Mac OS (64bit) and Windows (64bit):

```bash
go run make.go -crosscompile
```

Or you can just use go with `GO15VENDOREXPERIMENT` enabled:

```bash
export $GO15VENDOREXPERIMENT=1
go install
```

## Roadmap

- Actions
  - `delete`: Delete a given subdomain record

## Contribute

If you find a bug or if you want to add or improve some feature please create an issue or send me a pull requests.
All contributions are welcome.
