// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/pearkes/dnsimple"
)

// dnsClient provides functions for updating DNS records.
type dnsClient interface {
	UpdateRecord(domain string, id string, opts *dnsimple.ChangeRecord) (string, error)
	GetRecords(domain string) ([]dnsimple.Record, error)
}

type testDNSClient struct {
	updateRecordFunc func(domain string, id string, opts *dnsimple.ChangeRecord) (string, error)
	getRecordsFunc   func(domain string) ([]dnsimple.Record, error)
}

func (dnsClient *testDNSClient) UpdateRecord(domain string, id string, opts *dnsimple.ChangeRecord) (string, error) {
	return dnsClient.updateRecordFunc(domain, id, opts)
}

func (dnsClient *testDNSClient) GetRecords(domain string) ([]dnsimple.Record, error) {
	return dnsClient.getRecordsFunc(domain)
}
