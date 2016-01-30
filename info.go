// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/pearkes/dnsimple"
)

type dnsInfoProviderFactory interface {
	CreateInfoProvider() dnsInfoProvider
}

type dnsimpleInfoProviderFactory struct {
	clientFactory dnsClientFactory
}

func (infoFactory dnsimpleInfoProviderFactory) CreateInfoProvider() dnsInfoProvider {
	dnsimpleInfoProvider := &dnsimpleInfoProvider{infoFactory.clientFactory}
	return dnsimpleInfoProvider
}

type dnsInfoProvider interface {
	GetDomainNames() ([]string, error)
	GetDomainRecords(domain string) ([]dnsimple.Record, error)
	GetSubdomainRecord(domain, subdomain, recordType string) (dnsimple.Record, error)
	GetSubdomainRecords(domain, subdomain string) ([]dnsimple.Record, error)
}

// dnsimpleInfoProvider returns DNS records from the DNSimple API.
type dnsimpleInfoProvider struct {
	clientFactory dnsClientFactory
}

// GetDomainNames returns a list of all available domain names.
func (infoProvider *dnsimpleInfoProvider) GetDomainNames() ([]string, error) {

	client, clientError := infoProvider.getClient()
	if clientError != nil {
		return nil, fmt.Errorf("No DNS client available")
	}

	domains, err := client.GetDomains()
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

	client, clientError := infoProvider.getClient()
	if clientError != nil {
		return nil, fmt.Errorf("No DNS client available")
	}

	// get all DNS records for the given domain
	records, err := client.GetRecords(domain)
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

// getClient returns a DNS client instance or an error if the creation of the client failed.
func (infoProvider *dnsimpleInfoProvider) getClient() (dnsClient, error) {
	if infoProvider.clientFactory == nil {
		return nil, fmt.Errorf("No DNS client factory available")
	}

	client, err := infoProvider.clientFactory.CreateClient()
	if err != nil {
		return nil, fmt.Errorf("Unable to create DNS client. %s", err.Error())
	}

	return client, nil
}
