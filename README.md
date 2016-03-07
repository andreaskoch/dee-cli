# Dee CLI

Update DNS records from the command-line

`dee` is single self-contained command-line utility for updating subdomain records that are managed by DNSimple written go that works on Linux, Mac OS and Windows alike.

[![Build Status](https://travis-ci.org/andreaskoch/dee-cli.svg?branch=master)](https://travis-ci.org/andreaskoch/dee-cli)

## Usage

```bash
dee <action> [arguments ...]
```

Get help:

```bash
dee --help
```

**Actions**:

- `login` to the DNSimple API
- `logout`
- `create` an address record for a given domain
- `list` all available domain, subdomain and DNS records
- `update` a given address record by name
- `delete` a given address record by name
- `createorupdate` a given address record

### Action: `login`

Save DNSimple API credentials to disc.

**Arguments**:

- `-email`: The e-mail address of your DNSimple account
- `-apitoken`: The DNSimple API token

**Example**:

```bash
dee login -email apiuser@example.com -apitoken TracsiflOgympacKoFieC
```

The credentials are saved to: `~/.dee/credentials.json`

### Action: `logout`

Remove any stored DNSimple API credentials from disc.

```bash
dee logout
```

### Action: `list`

List all available domains or subdomains.

**Arguments**

- `-domain`: A domain name (optional)
- `-subdomain`: A subdomain name (optional)

**Examples**

List all available domains:

```bash
dee list
```

List all subdomain of a given domain:

```bash
dee list -domain example.com
```

List all DNS records for a given subdomain:

```bash
dee list -domain example.com -subdomain www
```

### Action: `create`

Create an address record.

**Arguments**:

- `-domain`: A domain name (required)
- `-subdomain`: The subdomain name (required)
- `-ip`: An IPv4 or IPv6 address (required)
- `-ttl`: The time to live (TTL) for the DNS record in seconds (default: 600)

**Examples**:

Create an `AAAA` record for the subdomain `www`:

```bash
dee create -domain example.com -subdomain www -ip 2001:0db8:0000:0042:0000:8a2e:0370:7334
```

Create an `AAAA` record for the subdomain `www` with TTL of 1 minute:

```bash
dee create -domain example.com -subdomain www -ip 2001:0db8:0000:0042:0000:8a2e:0370:7334 -ttl 60
```

Create an `A` record for the subdomain `www`:

```bash
dee create -domain example.com -subdomain www -ip 10.2.1.3
```

The `-ip` parameter can also be passed via Stdin:

```bash
echo "2001:0db8:0000:0042:0000:8a2e:0370:7334" | dee create -domain example.com -subdomain www -ttl 3600
```

### Action: `delete`

Deletes an address record.

**Arguments**:

- `-domain`: A domain name (required)
- `-subdomain`: The subdomain name (required)
- `-type`: The address record type (required, e.g. "AAAA", "A")

**Examples**:

Delete the IPv6 address record for www.example.com:

```bash
dee delete -domain example.com -subdomain www -type AAAA
```

Delete the IPv4 address record for www.example.com:

```bash
dee delete -domain example.com -subdomain www -type A
```

### Action: `update`

Update the DNS record for a given sub domain

**Arguments**:

- `-domain`: A domain name (e.g. `example.com`)
- `-subdomain`: A subdomain name (e.g. `www`)
- `-ip`: An IPv4 or IPv6 address

**Examples**:

Set the `AAAA` record of `www.example.com` to the given IP address:

```bash
dee update -domain example.com -subdomain www -ip 2001:0db8:0000:0042:0000:8a2e:0370:7334
```

The `-ip` address parameter can also be passed in via Stdin:

```bash
echo "2001:0db8:0000:0042:0000:8a2e:0370:7334" | dee update -domain example.com -subdomain www
```

### Action: `createorupdate`

The create-or-update action can be used if you are not sure if the address record you are trying to update does already exist.
If the record exists it will be just updated; if it does not yet exist it will be created with the given IP address.

**Arguments**:

- `-domain`: A domain name (required)
- `-subdomain`: The subdomain name (required)
- `-ip`: An IPv4 or IPv6 address (required)
- `-ttl`: The time to live (TTL) for the DNS record in seconds (default: 600)

## Dependencies

dee uses the [github.com/andreaskoch/dee-ns](https://github.com/andreaskoch/dee-ns) library for creating, reading, updating and delting DNSimple DNS records.

## Installation & Build

Compile the dee binary for your current platform:

```bash
git clone git@github.com:andreaskoch/dee-cli.git && cd dee-cli

go run make.go -install
```

Compile the dee binaries for Linux (64bit, ARM, ARM5, ARM6, ARM7), Mac OS (64bit) and Windows (64bit):

```bash
go run make.go -crosscompile
```

Or you can just use go with `GO15VENDOREXPERIMENT` enabled:

```bash
export $GO15VENDOREXPERIMENT=1
go install
```

## Contribute

If you find a bug or if you want to add or improve some feature please create an issue or send me a pull requests.
All contributions are welcome.

If you are planning to sign up for [DNSimple](https://dnsimple.com) and if you like this tool I would be happy if you use this link: [https://dnsimple.com/r/381546095cf6a2](https://dnsimple.com/r/381546095cf6a2)

One month of free service for the both of us :dancers:
