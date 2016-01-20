// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/pearkes/dnsimple"
)

type dnsInfoProvider interface {
	GetSubdomainRecord(domain, subdomain string) (record dnsimple.Record, err error)
}

// newDNSimpleInfoProvider creates a new instance of the DNSimple info-provider.
func newDNSimpleInfoProvider(client *dnsimple.Client) *dnsimpleInfoProvider {
	return &dnsimpleInfoProvider{
		client: client,
	}
}

// dnsimpleInfoProvider returns DNS records from the DNSimple API.
type dnsimpleInfoProvider struct {
	client dnsClient
}

// GetSubdomainRecord return the subdomain record that matches the given name.
// If no matching subdomain was found or an error occurred while fetching the
// available records an error will be returned.
func (infoProvider *dnsimpleInfoProvider) GetSubdomainRecord(domain, subdomain string) (record dnsimple.Record, err error) {
	records, err := infoProvider.client.GetRecords(domain)
	if err != nil {
		return dnsimple.Record{}, err
	}

	for _, record := range records {
		if record.Name != subdomain {
			continue
		}

		return record, nil
	}

	return dnsimple.Record{}, fmt.Errorf("Domain %s.%s not found", subdomain, domain)
}
