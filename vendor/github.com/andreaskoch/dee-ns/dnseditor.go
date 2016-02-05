// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deens

import (
	"fmt"
	"github.com/pearkes/dnsimple"
	"net"
)

// The DNSRecordCreator interface offers functions for creating domain records.
type DNSRecordCreator interface {

	// CreateSubdomain creates a new subdomain address record.
	CreateSubdomain(domain, subDomainName string, timeToLive int, ip net.IP) error
}

// The DNSRecordUpdater interface offers functions for updating domain records.
type DNSRecordUpdater interface {

	// UpdateSubdomain sets ip address of the given subdomain.
	UpdateSubdomain(domain, subDomainName string, ip net.IP) error
}

// The DNSRecordDeleter interface offers functions for creating domain records.
type DNSRecordDeleter interface {

	// DeleteSubdomain removes subdomain address record of the given type.
	DeleteSubdomain(domain, subDomainName string, recordType string) error
}

// The DNSRecordEditor interface provides functions for editing DNS records.
type DNSRecordEditor interface {
	DNSRecordCreator
	DNSRecordUpdater
	DNSRecordDeleter
}

// NewDNSEditor creates an new DNSRecordEditor instance.
func NewDNSEditor(client DNSClient, infoProvider DNSInfoProvider) DNSRecordEditor {
	return &DNSEditor{client, infoProvider}
}

// DNSEditor updates DNSimple domain records.
type DNSEditor struct {
	client       DNSClient
	infoProvider DNSInfoProvider
}

// CreateSubdomain creates an address record for the given domain
func (editor *DNSEditor) CreateSubdomain(domain, subdomain string, timeToLive int, ip net.IP) error {

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

	// check if the record already exists
	recordType := getDNSRecordTypeByIP(ip)
	if _, err := editor.infoProvider.GetSubdomainRecord(domain, subdomain, recordType); err != nil {
		return fmt.Errorf("No address record of type %q found for %q", recordType, subdomain)
	}

	// create record
	changeRecord := &dnsimple.ChangeRecord{
		Name:  subdomain,
		Value: ip.String(),
		Type:  recordType,
		Ttl:   fmt.Sprintf("%s", timeToLive),
	}

	_, createError := editor.client.CreateRecord(domain, changeRecord)
	if createError != nil {
		return createError
	}

	return nil
}

// UpdateSubdomain updates the IP address of the given domain/subdomain.
func (editor *DNSEditor) UpdateSubdomain(domain, subdomain string, ip net.IP) error {

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

	// get the subdomain record
	recordType := getDNSRecordTypeByIP(ip)
	subdomainRecord, err := editor.infoProvider.GetSubdomainRecord(domain, subdomain, recordType)
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

	_, updateError := editor.client.UpdateRecord(domain, fmt.Sprintf("%v", subdomainRecord.Id), changeRecord)
	if updateError != nil {
		return updateError
	}

	return nil
}

// DeleteSubdomain deletes the address record of the given domain
func (editor *DNSEditor) DeleteSubdomain(domain, subdomain string, recordType string) error {

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

	// check if the record already exists
	subdomainRecord, subdomainError := editor.infoProvider.GetSubdomainRecord(domain, subdomain, recordType)
	if subdomainError != nil {
		return fmt.Errorf("No address record of type %q found for %q", recordType, subdomain)
	}

	deleteError := editor.client.DestroyRecord(domain, fmt.Sprintf("%d", subdomainRecord.Id))
	if deleteError != nil {
		return deleteError
	}

	return nil
}
