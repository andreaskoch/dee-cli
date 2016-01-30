// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/pearkes/dnsimple"
	"net"
)

// The updater interface offers functions for updating domain records.
type updater interface {
	UpdateSubdomain(domain, subDomainName string, ip net.IP) error
}

// dnsimpleUpdater updates DNSimple domain records.
type dnsimpleUpdater struct {
	clientFactory       dnsClientFactory
	infoProviderFactory dnsInfoProviderFactory
}

// UpdateSubdomain updates the IP address of the given domain/subdomain.
func (updater *dnsimpleUpdater) UpdateSubdomain(domain, subdomain string, ip net.IP) error {

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

	client, clientError := updater.getClient()
	if clientError != nil {
		return fmt.Errorf("No DNS client available")
	}

	infoClient, infoClientError := updater.getInfoProvider()
	if infoClientError != nil {
		return fmt.Errorf("No DNS info provider available")
	}

	// get the subdomain record
	recordType := getDNSRecordTypeByIP(ip)
	subdomainRecord, err := infoClient.GetSubdomainRecord(domain, subdomain, recordType)
	if err != nil {
		return fmt.Errorf("Unable to locate DNS record %q.", subdomain)
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
		Ttl:   fmt.Sprintf("%d", subdomainRecord.Ttl),
	}

	_, updateError := client.UpdateRecord(domain, fmt.Sprintf("%v", subdomainRecord.Id), changeRecord)
	if updateError != nil {
		return updateError
	}

	return nil
}

// getClient returns a DNS client instance or an error if the creation of the client failed.
func (updater *dnsimpleUpdater) getClient() (dnsClient, error) {
	if updater.clientFactory == nil {
		return nil, fmt.Errorf("No DNS client factory available.")
	}

	client, err := updater.clientFactory.CreateClient()
	if err != nil {
		return nil, fmt.Errorf("Unable to create DNS client. %s", err.Error())
	}

	return client, nil
}

// getInfoProvider returns a DNS info provider instance or an error if the creation of the provider failed.
func (updater *dnsimpleUpdater) getInfoProvider() (dnsInfoProvider, error) {
	if updater.infoProviderFactory == nil {
		return nil, fmt.Errorf("No DNS info provider factory available.")
	}

	infoProvider := updater.infoProviderFactory.CreateInfoProvider()
	return infoProvider, nil
}
