# Dee NS

A go-library for updating DNS records

`github.com/andreaskoch/dee-ns` is a golang library for updating subdomain records that are managed by DNSimple.

[![Build Status](https://travis-ci.org/andreaskoch/dee-ns.svg?branch=master)](https://travis-ci.org/andreaskoch/dee-ns)

## Usage

Create a new DNS client instance and fetch all domain names:

```go
import (
	"fmt"
	"os"
	"github.com/andreaskoch/dee-ns"
)

func main() {
	// assemble the API credentials
	credentials := deens.APICredentials{"john.doe@example.com", "ApItOken"}

	// create a new DNS client
	dnsClient, clientError := deens.NewDNSClient(credentials)
	if clientError != nil {
		fmt.Fprintf(os.Stderr, "Unable to create DNS client: %s", clientError.Error())
		os.Exit(1)
	}

	// fetch all available domains
	domains, domainsError := dnsClient.GetDomains()
	if domainsError != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch domains: %s", domainsError.Error())
		os.Exit(1)
	}

	// print domain names to stdout
	for _, domain := range domains {
		fmt.Fprintf(os.Stdout, "%s\n", domain.Name)
	}
}

```

Create a new DNS info provider:

```go
import (
	"fmt"
	"os"
	"github.com/andreaskoch/dee-ns"
)

func main() {

	// create a DNS client
	credentials := APICredentials{"john.doe@example.com", "ApItOken"}
	dnsClient, clientError := deens.NewDNSClient(credentials)
	if clientError != nil {
		fmt.Fprintf(os.Stderr, "Unable to create DNS client: %s", clientError.Error())
		os.Exit(1)
	}

	// create a new DNS info provider instance
	infoProvider := deens.NewDNSInfoProvider(dnsClient)

	// get all domain names
	domainNames, domainNamesError := infoProvider.GetDomainNames()
	if domainNamesError != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch domain names: %s", domainNamesError.Error())
		os.Exit(1)
	}

	// print a list all domain names
	for _, domainName := range domainNames {
		fmt.Fprintf(os.Stdout, "%s\n", domainName)
	}
}
```

Create a new DNS editor instance:

```go
import (
	"fmt"
	"net"
	"os"
)

func main() {

	// create a DNS client
	credentials := APICredentials{"john.doe@example.com", "ApItOken"}
	dnsClient, clientError := NewDNSClient(credentials)
	if clientError != nil {
		fmt.Fprintf(os.Stderr, "Unable to create DNS client: %s", clientError.Error())
		os.Exit(1)
	}

	// create a new DNS info provider instance
	dnsInfoProvider := NewDNSInfoProvider(dnsClient)

	// create a new DNS editor instance
	dnsEditor := NewDNSEditor(dnsClient, dnsInfoProvider)

	// create an DNS A record for www.example.com pointing to 127.0.0.1
	domain := "example.com"
	subDomainName := "www"
	timeToLive := 600
	ip := net.ParseIP("127.0.0.1")

	createSubdomainError := dnsEditor.CreateSubdomain(domain, subDomainName, timeToLive, ip)
	if createSubdomainError != nil {
		fmt.Fprintf(os.Stderr, "Failed to create subdomain %s.%s: %s", subDomainName, domain, clientError.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Created subdomain %s.%s", subDomainName, domain)
}
```

## Dependencies

dee-ns uses the [github.com/pearkes/dnsimple](https://github.com/pearkes/dnsimple) library for communicating with the DNSimple API.

## Contribute

If you find a bug or if you want to add or improve some feature please create an issue or send me a pull requests.
All contributions are welcome.
