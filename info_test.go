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
	recordType := "AAAA"
	dnsClient := &testDNSClient{
		getRecordsFunc: func(domain string) ([]dnsimple.Record, error) {
			return nil, fmt.Errorf("Unable to fetch DNS records")
		},
	}

	infoProvider := dnsimpleInfoProvider{
		client: dnsClient,
	}

	// act
	_, err := infoProvider.GetSubdomainRecord(domain, subdomain, recordType)

	// assert
	if err == nil {
		t.Fail()
		t.Errorf("GetSubdomainRecord(%q, %q, %q) should return an error if the DNS client responds with an error.", domain, subdomain, recordType)
	}

}

// GetSubdomainRecord should return an error if the DNS clients returns no records.
func Test_GetSubdomainRecord_DNSClientReturnsNoRecords_ErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	recordType := "AAAA"
	dnsClient := &testDNSClient{
		getRecordsFunc: func(domain string) ([]dnsimple.Record, error) {
			return nil, nil
		},
	}

	infoProvider := dnsimpleInfoProvider{
		client: dnsClient,
	}

	// act
	_, err := infoProvider.GetSubdomainRecord(domain, subdomain, recordType)

	// assert
	if err == nil {
		t.Fail()
		t.Errorf("GetSubdomainRecord(%q, %q, %q) should return an error if the DNS client does not return records.", domain, subdomain, recordType)
	}

}

// GetSubdomainRecord should return the first record that has a matching name.
func Test_GetSubdomainRecord_FirstRecordMatchingTheSubdomainIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	recordType := "AAAA"
	dnsClient := &testDNSClient{
		getRecordsFunc: func(domain string) ([]dnsimple.Record, error) {
			return []dnsimple.Record{
				dnsimple.Record{Name: "aaa", RecordType: "AAAA", Id: 1},
				dnsimple.Record{Name: "bbb", RecordType: "AAAA", Id: 2},
				dnsimple.Record{Name: "www", RecordType: "AAAA", Id: 3},
				dnsimple.Record{Name: "www", RecordType: "AAAA", Id: 4},
			}, nil
		},
	}

	infoProvider := dnsimpleInfoProvider{
		client: dnsClient,
	}

	// act
	resultRecord, _ := infoProvider.GetSubdomainRecord(domain, subdomain, recordType)

	// assert
	if resultRecord.Id != 3 {
		t.Fail()
		t.Errorf("GetSubdomainRecord(%q, %q, %q) should have returned the %q-record but returned %q instead.", domain, subdomain, recordType, subdomain, resultRecord.Name)
	}

}

// GetSubdomainRecord should return an error if no matching record is found.
func Test_GetSubdomainRecord_NoMatchingSubdomainRecordFound_ErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "nonexistingsubdomain"
	recordType := "AAAA"
	dnsClient := &testDNSClient{
		getRecordsFunc: func(domain string) ([]dnsimple.Record, error) {
			return []dnsimple.Record{
				dnsimple.Record{Name: "aaa", RecordType: "AAAA", Id: 1},
				dnsimple.Record{Name: "bbb", RecordType: "AAAA", Id: 2},
				dnsimple.Record{Name: "www", RecordType: "AAAA", Id: 3},
			}, nil
		},
	}

	infoProvider := dnsimpleInfoProvider{
		client: dnsClient,
	}

	// act
	_, err := infoProvider.GetSubdomainRecord(domain, subdomain, recordType)

	// assert
	if err == nil {
		t.Fail()
		t.Errorf("GetSubdomainRecord(%q, %q, %q) should return an error if no matching DNS record was found.", domain, subdomain, recordType)
	}

}

// GetSubdomainRecord should return an error if no record is found that matches the given record type.
func Test_GetSubdomainRecord_NoMatchingRecordTypeFound_ErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	recordType := "AAAA"
	dnsClient := &testDNSClient{
		getRecordsFunc: func(domain string) ([]dnsimple.Record, error) {
			return []dnsimple.Record{
				dnsimple.Record{Name: "aaa", RecordType: "AAAA", Id: 1},
				dnsimple.Record{Name: "bbb", RecordType: "AAAA", Id: 2},
				dnsimple.Record{Name: "www", RecordType: "A", Id: 3},
			}, nil
		},
	}

	infoProvider := dnsimpleInfoProvider{
		client: dnsClient,
	}

	// act
	_, err := infoProvider.GetSubdomainRecord(domain, subdomain, recordType)

	// assert
	if err == nil {
		t.Fail()
		t.Errorf("GetSubdomainRecord(%q, %q, %q) should return an error if no matching DNS record was found.", domain, subdomain, recordType)
	}

}
