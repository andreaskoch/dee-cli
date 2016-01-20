// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/pearkes/dnsimple"
	"testing"
)

// GetSubdomainRecord should return an error if the DNS clients returns an error instead of DNS records.
func Test_GetSubdomainRecord_DNSClientReturnsError_ErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	dnsClient := &testDNSClient{
		getRecordsFunc: func(domain string) ([]dnsimple.Record, error) {
			return nil, fmt.Errorf("Unable to fetch DNS records")
		},
	}

	infoProvider := dnsimpleInfoProvider{
		client: dnsClient,
	}

	// act
	_, err := infoProvider.GetSubdomainRecord(domain, subdomain)

	// assert
	if err == nil {
		t.Fail()
		t.Errorf("GetSubdomainRecord(%q, %q) should return an error if the DNS client responds with an error.", domain, subdomain)
	}

}

// GetSubdomainRecord should return an error if the DNS clients returns no records.
func Test_GetSubdomainRecord_DNSClientReturnsNoRecords_ErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	dnsClient := &testDNSClient{
		getRecordsFunc: func(domain string) ([]dnsimple.Record, error) {
			return nil, nil
		},
	}

	infoProvider := dnsimpleInfoProvider{
		client: dnsClient,
	}

	// act
	_, err := infoProvider.GetSubdomainRecord(domain, subdomain)

	// assert
	if err == nil {
		t.Fail()
		t.Errorf("GetSubdomainRecord(%q, %q) should return an error if the DNS client does not return records.", domain, subdomain)
	}

}

// GetSubdomainRecord should return the first record that has a matching name.
func Test_GetSubdomainRecord_FirstRecordMatchingTheSubdomainIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	dnsClient := &testDNSClient{
		getRecordsFunc: func(domain string) ([]dnsimple.Record, error) {
			return []dnsimple.Record{
				dnsimple.Record{Name: "aaa", Id: 1},
				dnsimple.Record{Name: "bbb", Id: 2},
				dnsimple.Record{Name: "www", Id: 3},
				dnsimple.Record{Name: "www", Id: 4},
			}, nil
		},
	}

	infoProvider := dnsimpleInfoProvider{
		client: dnsClient,
	}

	// act
	resultRecord, _ := infoProvider.GetSubdomainRecord(domain, subdomain)

	// assert
	if resultRecord.Id != 3 {
		t.Fail()
		t.Errorf("GetSubdomainRecord(%q, %q) should have returned the %q-record but returned %q instead.", domain, subdomain, subdomain, resultRecord.Name)
	}

}

// GetSubdomainRecord should return an error if no matching record is found.
func Test_GetSubdomainRecord_NoMatchingRecordFound_ErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "nonexistingsubdomain"
	dnsClient := &testDNSClient{
		getRecordsFunc: func(domain string) ([]dnsimple.Record, error) {
			return []dnsimple.Record{
				dnsimple.Record{Name: "aaa", Id: 1},
				dnsimple.Record{Name: "bbb", Id: 2},
				dnsimple.Record{Name: "www", Id: 3},
			}, nil
		},
	}

	infoProvider := dnsimpleInfoProvider{
		client: dnsClient,
	}

	// act
	_, err := infoProvider.GetSubdomainRecord(domain, subdomain)

	// assert
	if err == nil {
		t.Fail()
		t.Errorf("GetSubdomainRecord(%q, %q) should return an error if no matching DNS record was found.", domain, subdomain)
	}

}
