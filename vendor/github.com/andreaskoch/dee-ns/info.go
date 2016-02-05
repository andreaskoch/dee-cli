// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deens

import (
	"fmt"
	"github.com/pearkes/dnsimple"
)

// The DNSInfoProvider interface offer DNS info functions.
type DNSInfoProvider interface {

	// GetDomainNames returns a list of domain names.
	// Returns an error if the domain names cannot be fetched.
	GetDomainNames() ([]string, error)

	// GetDomainRecords returns all DNS records for the given domain.
	// Returns an error of the DNS records cannot be fetched or the
	// given domain was not found.
	GetDomainRecords(domain string) ([]dnsimple.Record, error)

	// GetSubdomainRecord returns the DNS record for the given domain, subdomain and record type.
	// Returns an error if no DNS record was found.
	GetSubdomainRecord(domain, subdomain, recordType string) (dnsimple.Record, error)

	// GetSubdomainRecords returns a list of all available DNS records for the
	// given domain and subdomain.
	GetSubdomainRecords(domain, subdomain string) ([]dnsimple.Record, error)
}

// NewDNSInfoProvider creates a new DNS info provider instance.
func NewDNSInfoProvider(client DNSClient) DNSInfoProvider {
	return &dnsimpleInfoProvider{client}
}

// dnsimpleInfoProvider returns DNS records from the DNSimple API.
type dnsimpleInfoProvider struct {
	client DNSClient
}

// GetDomainNames returns a list of all available domain names.
func (infoProvider *dnsimpleInfoProvider) GetDomainNames() ([]string, error) {

	domains, err := infoProvider.client.GetDomains()
	if err != nil {
		return nil, err
	}

	var domainNames []string
	for _, domain := range domains {
		domainNames = append(domainNames, domain.Name)
	}

	return domainNames, nil
}

// GetDomainRecords returns all DNS records for the given domain.
func (infoProvider *dnsimpleInfoProvider) GetDomainRecords(domain string) ([]dnsimple.Record, error) {

	return infoProvider.getDNSRecords(domain, func(record dnsimple.Record) bool {
		return true
	})

}

// GetSubdomainRecord return the subdomain record that matches the given name and record type.
// If no matching subdomain was found or an error occurred while fetching the available records
// an error will be returned.
func (infoProvider *dnsimpleInfoProvider) GetSubdomainRecord(domain, subdomain, recordType string) (dnsimple.Record, error) {

	// get all records that have matching subdomain name and record type
	records, err := infoProvider.getDNSRecords(domain, func(record dnsimple.Record) bool {
		return record.Name == subdomain && record.RecordType == recordType
	})

	// error while fetching DNS records
	if err != nil {
		return dnsimple.Record{}, err
	}

	// no records found
	if len(records) == 0 {
		return dnsimple.Record{}, fmt.Errorf("No record found for %s.%s", subdomain, domain)
	}

	// return the first record found
	return records[0], nil
}

// GetSubdomainRecords returns all DNS records for the given subdomain.
func (infoProvider *dnsimpleInfoProvider) GetSubdomainRecords(domain, subdomain string) ([]dnsimple.Record, error) {

	return infoProvider.getDNSRecords(domain, func(record dnsimple.Record) bool {
		return record.Name == subdomain
	})

}

// getDNSRecords returns all DNS records for the given domain that pass the given filter expression.
func (infoProvider *dnsimpleInfoProvider) getDNSRecords(domain string, includeInResult func(record dnsimple.Record) bool) ([]dnsimple.Record, error) {

	// get all DNS records for the given domain
	records, err := infoProvider.client.GetRecords(domain)
	if err != nil {
		return nil, err
	}

	var filteredRecords []dnsimple.Record
	for _, record := range records {
		if !includeInResult(record) {
			continue
		}

		filteredRecords = append(filteredRecords, record)
	}

	return filteredRecords, nil
}
