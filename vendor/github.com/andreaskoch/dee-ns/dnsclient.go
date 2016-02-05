// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deens

import (
	"fmt"
	"github.com/pearkes/dnsimple"
)

// NewDNSClient creates a new DNS client instance for the given credentials.
func NewDNSClient(credentials APICredentials) (DNSClient, error) {
	dnsimpleClient, dnsimpleClientError := dnsimple.NewClient(credentials.Email, credentials.Token)
	if dnsimpleClientError != nil {
		return nil, fmt.Errorf("Unable to create DNSimple client. Error: %s", dnsimpleClientError.Error())
	}

	return dnsimpleClient, nil
}

// DNSClient provides functions for updating DNS records.
type DNSClient interface {
	// UpdateRecord update the DNS record with the given id.
	UpdateRecord(domain string, id string, opts *dnsimple.ChangeRecord) (string, error)

	// GetRecords returns all DNS records for the given domain.
	GetRecords(domain string) ([]dnsimple.Record, error)

	// GetDomains returns a list of domain.
	GetDomains() ([]dnsimple.Domain, error)

	// CreateRecord creates a new DNS record for the given domain.
	CreateRecord(domain string, opts *dnsimple.ChangeRecord) (string, error)

	// DestroyRecord deletes the DNS record with the given id.
	DestroyRecord(domain string, id string) error
}
