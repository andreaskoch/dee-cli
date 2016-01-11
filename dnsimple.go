package main

import (
	"flag"
	"fmt"
	"github.com/pearkes/dnsimple"
	"net"
	"os"
	"strings"
)

var tempEmail = "user@example.com"
var tempToken = "TracsiflOgympacKoFieC"

const ACTIONUPDATE = "update"
const ACTIONLOGIN = "update"

var action string

var (
	updateSubdomainArguments = flag.NewFlagSet("update-subdomain", flag.ContinueOnError)
	domain                   = updateSubdomainArguments.String("domain", "", "Domain (e.g. example.com")
	subdomain                = updateSubdomainArguments.String("subdomain", "", "Subdomain (e.g. wwww)")
	ipAddress                = updateSubdomainArguments.String("ip", "", "IP address (e.g. 127.0.0.1. ::1")
)

func init() {

	executableName := os.Args[0]

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s updates DNS records via DNSimple.\n", executableName)
		fmt.Fprintf(os.Stderr, "\n")

		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s <action> [arguments ...]\n", executableName)
		fmt.Fprintf(os.Stderr, "\n")

		fmt.Fprintf(os.Stderr, "Actions:\n")
		fmt.Fprintf(os.Stderr, "%10s  %s\n", ACTIONUPDATE, "Update the DNS record for a given sub domain")
		fmt.Fprintf(os.Stderr, "\n")

		fmt.Fprintf(os.Stderr, "Action: %s\n", ACTIONUPDATE)
		updateSubdomainArguments.PrintDefaults()
	}

}

func main() {

	// get action
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	action = strings.TrimSpace(strings.ToLower(os.Args[1]))
	switch action {
	case ACTIONUPDATE:
		{
			update()
		}

	default:
		{
			fmt.Fprintf(os.Stderr, "Unknown action")
			os.Exit(1)
		}
	}
}

// update executes the update domain update action
func update() {
	dnsimpleClient, err := dnsimple.NewClient(tempEmail, tempToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create DNSimple client. Error: %s", err.Error())
		os.Exit(1)
	}

	// parse the arguments
	updateSubdomainArguments.Parse(os.Args[2:])

	// domain
	if *domain == "" {
		fmt.Fprintf(os.Stderr, "No domain supplied.")
		os.Exit(1)
	}

	// subdomain
	if *subdomain == "" {
		fmt.Fprintf(os.Stderr, "No subdomain supplied.")
		os.Exit(1)
	}

	// take ip from stdin
	if *ipAddress == "" {
		ipAddressFromStdin := ""
		fmt.Fscan(os.Stdin, &ipAddressFromStdin)
		ipAddress = &ipAddressFromStdin
	}

	if *ipAddress == "" {
		fmt.Fprintf(os.Stderr, "No IP address supplied.")
		os.Exit(1)
	}

	ip := net.ParseIP(*ipAddress)
	updateError := updateSubdomain(dnsimpleClient, *domain, *subdomain, ip)
	if updateError != nil {
		fmt.Fprintf(os.Stderr, "%s", updateError.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Updated: %s.%s → %s", *subdomain, *domain, ip.String())
}

// updateSubdomain updates the IP address of the given domain/subdomain
func updateSubdomain(client *dnsimple.Client, domain, subdomain string, ip net.IP) error {

	// get the subdomain record
	subdomainRecord, err := getSubdomainRecord(client, domain, subdomain)
	if err != nil {
		return err
	}

	// check if an update is necessary
	if subdomainRecord.Content == ip.String() {
		return fmt.Errorf("No update required. IP address did not change (%s).", subdomainRecord.Content)
	}

	// update the record
	changeRecord := &dnsimple.ChangeRecord{
		Name:  subdomainRecord.Name,
		Value: ip.String(),
		Type:  subdomainRecord.RecordType,
		Ttl:   fmt.Sprintf("%#v", subdomainRecord.Ttl),
	}

	_, updateError := client.UpdateRecord(domain, fmt.Sprintf("%v", subdomainRecord.Id), changeRecord)
	if updateError != nil {
		return updateError
	}

	return nil
}

// getSubdomainRecord return the subdomain record that matches the given name.
// If no matching subdomain was found or an error occurred while fetching the
// available records an error will be returned.
func getSubdomainRecord(client *dnsimple.Client, domain, subdomain string) (record dnsimple.Record, err error) {
	records, err := client.GetRecords(domain)
	if err != nil {
		return dnsimple.Record{}, err
	}

	for _, record := range records {
		if record.Name != subdomain {
			continue
		}

		return record, nil
	}

	return dnsimple.Record{}, fmt.Errorf("Domain %s.%s not found", subdomain, domain)
}
