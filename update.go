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

// newDNSimpleUpdater creates a new DNSimple updater instance.
func newDNSimpleUpdater(client *dnsimple.Client, infoProvider *dnsimpleInfoProvider) updater {
	return &dnsimpleUpdater{
		client:       client,
		infoProvider: infoProvider,
	}
}

// dnsimpleUpdater updates DNSimple domain records.
type dnsimpleUpdater struct {
	client       dnsClient
	infoProvider dnsInfoProvider
}

// updateSubdomain updates the IP address of the given domain/subdomain
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

	// get the subdomain record
	recordType := getDNSRecordTypeByIP(ip)
	subdomainRecord, err := updater.infoProvider.GetSubdomainRecord(domain, subdomain, recordType)
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
		Ttl:   fmt.Sprintf("%d", subdomainRecord.Ttl),
	}

	_, updateError := updater.client.UpdateRecord(domain, fmt.Sprintf("%v", subdomainRecord.Id), changeRecord)
	if updateError != nil {
		return updateError
	}

	return nil
}
