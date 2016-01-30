// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/pearkes/dnsimple"
	"net"
)

// The dnsRecordCreator interface offers functions for creating domain records.
type dnsRecordCreator interface {
	CreateSubdomain(domain, subDomainName string, timeToLive int, ip net.IP) error
}

// dnsimpleCreator updates DNSimple domain records.
type dnsimpleCreator struct {
	clientFactory       dnsClientFactory
	infoProviderFactory dnsInfoProviderFactory
}

// CreateSubdomain creates an address record for the given domain
func (creator *dnsimpleCreator) CreateSubdomain(domain, subdomain string, timeToLive int, ip net.IP) error {

	// validate parameters
	if isValidDomain(domain) == false {
		return fmt.Errorf("The domain name is invalid: %q", domain)
	}

	if isValidSubdomain(subdomain) == false {
		return fmt.Errorf("The domain name is invalid: %q", subdomain)
	}

	if ip == nil {
		return fmt.Errorf("No ip supplied")
	}

	client, clientError := creator.getClient()
	if clientError != nil {
		return fmt.Errorf("No DNS client available")
	}

	infoClient, infoClientError := creator.getInfoProvider()
	if infoClientError != nil {
		return fmt.Errorf("No DNS info provider available")
	}

	// check if the record already exists
	recordType := getDNSRecordTypeByIP(ip)
	if subdomainRecord, _ := infoClient.GetSubdomainRecord(domain, subdomain, recordType); subdomainRecord.Id != 0 {
		return fmt.Errorf("No address record of type %q found for %q", recordType, subdomain)
	}

	// create record
	changeRecord := &dnsimple.ChangeRecord{
		Name:  subdomain,
		Value: ip.String(),
		Type:  recordType,
		Ttl:   fmt.Sprintf("%s", timeToLive),
	}

	_, createError := client.CreateRecord(domain, changeRecord)
	if createError != nil {
		return createError
	}

	return nil
}

// getClient returns a DNS client instance or an error if the creation of the client failed.
func (creator *dnsimpleCreator) getClient() (dnsClient, error) {
	if creator.clientFactory == nil {
		return nil, fmt.Errorf("No DNS client factory available.")
	}

	client, err := creator.clientFactory.CreateClient()
	if err != nil {
		return nil, fmt.Errorf("Unable to create DNS client. %s", err.Error())
	}

	return client, nil
}

// getInfoProvider returns a DNS info provider instance or an error if the creation of the provider failed.
func (creator *dnsimpleCreator) getInfoProvider() (dnsInfoProvider, error) {
	if creator.infoProviderFactory == nil {
		return nil, fmt.Errorf("No DNS info provider factory available.")
	}

	infoProvider := creator.infoProviderFactory.CreateInfoProvider()
	return infoProvider, nil
}
