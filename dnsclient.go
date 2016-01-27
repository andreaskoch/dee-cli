// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/pearkes/dnsimple"
)

// dnsClientFactory provides the ability to create DNS clients.
type dnsClientFactory interface {
	// CreateClient create a new dnsClient client instance.
	CreateClient() (dnsClient, error)
}

// dnsimpleClientFactory creates DNSimple clients.
type dnsimpleClientFactory struct {
	credentialStore credentialStore
}

// CreateClient create a new DNSimple client instance.
func (clientFactory dnsimpleClientFactory) CreateClient() (dnsClient, error) {
	// get the credentials
	credentials, credentialError := clientFactory.credentialStore.GetCredentials()
	if credentialError != nil {
		return nil, fmt.Errorf("%s", credentialError.Error())
	}

	// create a DNSimple client
	dnsimpleClient, dnsimpleClientError := dnsimple.NewClient(credentials.Email, credentials.Token)
	if dnsimpleClientError != nil {
		return nil, fmt.Errorf("Unable to create DNSimple client. Error: %s", dnsimpleClientError.Error())
	}

	return dnsimpleClient, nil
}

// dnsClient provides functions for updating DNS records.
type dnsClient interface {
	UpdateRecord(domain string, id string, opts *dnsimple.ChangeRecord) (string, error)
	GetRecords(domain string) ([]dnsimple.Record, error)
	GetDomains() ([]dnsimple.Domain, error)
}

type testDNSClient struct {
	updateRecordFunc func(domain string, id string, opts *dnsimple.ChangeRecord) (string, error)
	getRecordsFunc   func(domain string) ([]dnsimple.Record, error)
	getDomainsFunc   func() ([]dnsimple.Domain, error)
}

func (dnsClient *testDNSClient) UpdateRecord(domain string, id string, opts *dnsimple.ChangeRecord) (string, error) {
	return dnsClient.updateRecordFunc(domain, id, opts)
}

func (dnsClient *testDNSClient) GetRecords(domain string) ([]dnsimple.Record, error) {
	return dnsClient.getRecordsFunc(domain)
}

func (dnsClient *testDNSClient) GetDomains() ([]dnsimple.Domain, error) {
	return dnsClient.getDomainsFunc()
}
