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

// The dnsRecordUpdater interface offers functions for updating domain records.
type dnsRecordUpdater interface {
	UpdateSubdomain(domain, subDomainName string, ip net.IP) error
}

// The dnsRecordDeleter interface offers functions for creating domain records.
type dnsRecordDeleter interface {
	DeleteSubdomain(domain, subDomainName string, recordType string) error
}

// dnsEditor updates DNSimple domain records.
type dnsEditor struct {
	clientFactory       dnsClientFactory
	infoProviderFactory dnsInfoProviderFactory
}

// CreateSubdomain creates an address record for the given domain
func (editor *dnsEditor) CreateSubdomain(domain, subdomain string, timeToLive int, ip net.IP) error {

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

	client, clientError := editor.getClient()
	if clientError != nil {
		return fmt.Errorf("No DNS client available")
	}

	infoClient, infoClientError := editor.getInfoProvider()
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

// UpdateSubdomain updates the IP address of the given domain/subdomain.
func (editor *dnsEditor) UpdateSubdomain(domain, subdomain string, ip net.IP) error {

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

	client, clientError := editor.getClient()
	if clientError != nil {
		return fmt.Errorf("No DNS client available")
	}

	infoClient, infoClientError := editor.getInfoProvider()
	if infoClientError != nil {
		return fmt.Errorf("No DNS info provider available")
	}

	// get the subdomain record
	recordType := getDNSRecordTypeByIP(ip)
	subdomainRecord, err := infoClient.GetSubdomainRecord(domain, subdomain, recordType)
	if err != nil {
		return fmt.Errorf("No address record of type %q found for %q", recordType, subdomain)
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

// DeleteSubdomain deletes the address record of the given domain
func (editor *dnsEditor) DeleteSubdomain(domain, subdomain string, recordType string) error {

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

	client, clientError := editor.getClient()
	if clientError != nil {
		return fmt.Errorf("No DNS client available")
	}

	infoClient, infoClientError := editor.getInfoProvider()
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
func (editor *dnsEditor) getClient() (dnsClient, error) {
	if editor.clientFactory == nil {
		return nil, fmt.Errorf("No DNS client factory available")
	}

	client, err := editor.clientFactory.CreateClient()
	if err != nil {
		return nil, fmt.Errorf("Unable to create DNS client. %s", err.Error())
	}

	return client, nil
}

// getInfoProvider returns a DNS info provider instance or an error if the creation of the provider failed.
func (editor *dnsEditor) getInfoProvider() (dnsInfoProvider, error) {
	if editor.infoProviderFactory == nil {
		return nil, fmt.Errorf("No DNS info provider factory available")
	}

	infoProvider := editor.infoProviderFactory.CreateInfoProvider()
	return infoProvider, nil
}
