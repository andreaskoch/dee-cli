// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

// The dnsRecordDeleter interface offers functions for creating domain records.
type dnsRecordDeleter interface {
	DeleteSubdomain(domain, subDomainName string, recordType string) error
}

// dnsimpleDeleter updates DNSimple domain records.
type dnsimpleDeleter struct {
	clientFactory       dnsClientFactory
	infoProviderFactory dnsInfoProviderFactory
}

// DeleteSubdomain deletes the address record of the given domain
func (deleter *dnsimpleDeleter) DeleteSubdomain(domain, subdomain string, recordType string) error {

	// validate parameters
	if isValidDomain(domain) == false {
		return fmt.Errorf("The domain name is invalid: %q", domain)
	}

	if isValidSubdomain(subdomain) == false {
		return fmt.Errorf("The domain name is invalid: %q", subdomain)
	}

	if recordType != "AAAA" && recordType != "A" {
		return fmt.Errorf("The given record type is invalid: %q", subdomain)
	}

	client, clientError := deleter.getClient()
	if clientError != nil {
		return fmt.Errorf("No DNS client available")
	}

	infoClient, infoClientError := deleter.getInfoProvider()
	if infoClientError != nil {
		return fmt.Errorf("No DNS info provider available")
	}

	// check if the record already exists
	subdomainRecord, subdomainError := infoClient.GetSubdomainRecord(domain, subdomain, recordType)
	if subdomainError != nil {
		return fmt.Errorf("No address record of type %q found for %q", recordType, subdomain)
	}

	deleteError := client.DestroyRecord(domain, fmt.Sprintf("%d", subdomainRecord.Id))
	if deleteError != nil {
		return deleteError
	}

	return nil
}

// getClient returns a DNS client instance or an error if the creation of the client failed.
func (deleter *dnsimpleDeleter) getClient() (dnsClient, error) {
	if deleter.clientFactory == nil {
		return nil, fmt.Errorf("No DNS client factory available.")
	}

	client, err := deleter.clientFactory.CreateClient()
	if err != nil {
		return nil, fmt.Errorf("Unable to create DNS client. %s", err.Error())
	}

	return client, nil
}

// getInfoProvider returns a DNS info provider instance or an error if the creation of the provider failed.
func (deleter *dnsimpleDeleter) getInfoProvider() (dnsInfoProvider, error) {
	if deleter.infoProviderFactory == nil {
		return nil, fmt.Errorf("No DNS info provider factory available.")
	}

	infoProvider := deleter.infoProviderFactory.CreateInfoProvider()
	return infoProvider, nil
}
