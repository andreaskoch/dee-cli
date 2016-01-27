// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/pearkes/dnsimple"
	"testing"
)

type testInfoProviderFactory struct {
	infoProvider dnsInfoProvider
}

func (infoFactory testInfoProviderFactory) CreateInfoProvider() dnsInfoProvider {
	return infoFactory.infoProvider
}

// testDNSInfoProvider is a DNS info-provider used for testing.
type testDNSInfoProvider struct {
	getDomainNamesFunc         func() ([]string, error)
	getAllDNSRecordsFunc       func(domain string) ([]dnsimple.Record, error)
	getSubdomainRecordFunc     func(domain, subdomain, recordType string) (dnsimple.Record, error)
	getSubdomainDNSRecordsFunc func(domain, subdomain string) ([]dnsimple.Record, error)
}

func (infoProvider testDNSInfoProvider) GetDomainNames() ([]string, error) {
	return infoProvider.getDomainNamesFunc()
}

func (infoProvider testDNSInfoProvider) GetAllDNSRecords(domain string) ([]dnsimple.Record, error) {
	return infoProvider.getAllDNSRecordsFunc(domain)
}

func (infoProvider testDNSInfoProvider) GetSubdomainRecord(domain, subdomain, recordType string) (record dnsimple.Record, err error) {
	return infoProvider.getSubdomainRecordFunc(domain, subdomain, recordType)
}

func (infoProvider testDNSInfoProvider) GetSubdomainDNSRecords(domain, subdomain string) ([]dnsimple.Record, error) {
	return infoProvider.getSubdomainDNSRecordsFunc(domain, subdomain)
}

func Test_NoDNSClientSupplied_GetSubdomainRecord_ErrorIsReturned(t *testing.T) {
	// arrange
	infoProvider := dnsimpleInfoProvider{}

	// act
	_, err := infoProvider.GetSubdomainRecord("example.com", "www", "AAAA")

	// assert
	if err == nil {
		t.Fail()
		t.Errorf("GetSubdomainRecord() should return an error if no DNS client was supplied to the info provider.")
	}

}

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

	clientFactory := testDNSClientFactory{dnsClient}
	infoProvider := dnsimpleInfoProvider{clientFactory}

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

	clientFactory := testDNSClientFactory{dnsClient}
	infoProvider := dnsimpleInfoProvider{clientFactory}

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

	clientFactory := testDNSClientFactory{dnsClient}
	infoProvider := dnsimpleInfoProvider{clientFactory}

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

	clientFactory := testDNSClientFactory{dnsClient}
	infoProvider := dnsimpleInfoProvider{clientFactory}

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

	clientFactory := testDNSClientFactory{dnsClient}
	infoProvider := dnsimpleInfoProvider{clientFactory}

	// act
	_, err := infoProvider.GetSubdomainRecord(domain, subdomain, recordType)

	// assert
	if err == nil {
		t.Fail()
		t.Errorf("GetSubdomainRecord(%q, %q, %q) should return an error if no matching DNS record was found.", domain, subdomain, recordType)
	}

}

// GetDomainNames should return an error if the info provider has no DNS client.
func Test_NoDNSClientGiven_GetDomainNames_ErrorIsReturned(t *testing.T) {
	// arrange
	clientFactory := testDNSClientFactory{}
	infoProvider := dnsimpleInfoProvider{clientFactory}

	// act
	_, err := infoProvider.GetDomainNames()

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetDomainNames() should return an error if there is no DNS client present.")
	}
}

// GetDomainNames should return an error if the DNS client returns one.
func Test_DNSClientReturnsAnError_GetDomainNames_ErrorReturned(t *testing.T) {
	// arrange
	dnsClient := &testDNSClient{
		getDomainsFunc: func() ([]dnsimple.Domain, error) {
			return []dnsimple.Domain{}, fmt.Errorf("Unable to fetch domains")
		},
	}

	clientFactory := testDNSClientFactory{dnsClient}
	infoProvider := dnsimpleInfoProvider{clientFactory}

	// act
	_, err := infoProvider.GetDomainNames()

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetDomainNames() should have returned an error because the DNS client returned one as well.")
	}
}

// GetDomainNames should return an empty list of the DNS client returns no domains.
func Test_DNSClientReturnsNoDomain_GetDomainNames_EmptyListIsReturned(t *testing.T) {
	// arrange
	dnsClient := &testDNSClient{
		getDomainsFunc: func() ([]dnsimple.Domain, error) {
			return []dnsimple.Domain{}, nil
		},
	}

	clientFactory := testDNSClientFactory{dnsClient}
	infoProvider := dnsimpleInfoProvider{clientFactory}

	// act
	names, _ := infoProvider.GetDomainNames()

	// assert
	if len(names) > 0 {
		t.Fail()
		t.Logf("GetDomainNames() not return any domain names because the DNS client returned none.")
	}
}

// GetDomainNames should returns all names of the domains returned by the DNS client.
func Test_DNSClientReturnsDomains_GetDomainNames_DomainNamesAreReturned(t *testing.T) {
	// arrange
	dnsClient := &testDNSClient{
		getDomainsFunc: func() ([]dnsimple.Domain, error) {
			return []dnsimple.Domain{
				dnsimple.Domain{Name: "example.com"},
				dnsimple.Domain{Name: "example.de"},
			}, nil
		},
	}

	clientFactory := testDNSClientFactory{dnsClient}
	infoProvider := dnsimpleInfoProvider{clientFactory}

	// act
	names, _ := infoProvider.GetDomainNames()

	// assert
	if len(names) != 2 {
		t.Fail()
		t.Logf("GetDomainNames() should have returned two domain.")
	}
}
