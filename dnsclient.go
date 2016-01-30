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
	CreateRecord(domain string, opts *dnsimple.ChangeRecord) (string, error)
}
