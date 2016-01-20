# dnsimple-cli

Update DNS records from the command-line via DNSimple

`dnsimple-cli` is a cross-platform command-line utility for updating subdomain records that are managed by DNSimple.

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

The IP address can also be passed in via Stdin:

```bash
echo "2001:0db8:0000:0042:0000:8a2e:0370:7334" | dnsimple-cli update -domain example.com -subdomain www
```

## Dependencies

dnsimple-cli uses the following third-party libraries:

- The [github.com/pearkes/dnsimple](https://github.com/pearkes/dnsimple) API library for communicating with the DNSimple API
- [github.com/mitchellh/go-homedir](https://github.com/mitchellh/go-homedir) for determining the user's home directory path
- [github.com/spf13/afero](https://github.com/spf13/afero) as a filesystem abstraction for testing

## Roadmap

- Actions
  - `logout`: Reset credentials
  - `create`: Create a subdomain record
  - `delete`: Delete a given subdomain record
  - `list`: List all subdomain records

## Contribute

If you find a bug or if you want to add or improve some feature please create an issue or send me a pull requests.
All contributions are welcome.
