# dnsimple-cli

A command-line utility for updating subdomain records that are managed by DNSimple

## Usage

Set the `AAAA` record of `www.example.com` to the given IP address:

```bash
dnsimple-cli -domain example.com -subdomain www -ip 2001:0db8:0000:0042:0000:8a2e:0370:7334
```

## Dependencies

- The [github.com/pearkes/dnsimple](https://github.com/pearkes/dnsimple) API library for communicating with the DNSimple API
