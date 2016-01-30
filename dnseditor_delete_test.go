// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/pearkes/dnsimple"
	"testing"
)

// testDNSDeleter deletes DNS records.
type testDNSDeleter struct {
	deleteSubdomainFunc func(domain, subdomain string, recordType string) error
}

func (editor *testDNSDeleter) DeleteSubdomain(domain, subdomain string, recordType string) error {
	return editor.deleteSubdomainFunc(domain, subdomain, recordType)
}

// If any of the given parameters is invalid DeleteSubdomain should respond with an error.
func Test_DeleteSubdomain_ParametersInvalid_ErrorIsReturned(t *testing.T) {
	// arrange
	inputs := []struct {
		domain     string
		subdomain  string
		recordType string
	}{
		{"example.com", "", "AAAA"},
		{"www", "", "AAAA"},
		{"", "", "AAAA"},
		{" ", " ", "AAAA"},
		{"example.com", "www", "AAAA"},
	}
	editor := dnsEditor{}

	for _, input := range inputs {

		// act
		err := editor.DeleteSubdomain(input.domain, input.subdomain, input.recordType)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("DeleteSubdomain(%q, %q, %q) should return an error.", input.domain, input.subdomain, input.recordType)
		}
	}
}

// DeleteSubdomain should return an error if the given subdomain does not exist.
func Test_DeleteSubdomain_ValidParameters_SubdomainNotFound_ErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	recordType := "AAAA"

	infoProvider := &testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (record dnsimple.Record, err error) {
			return dnsimple.Record{}, fmt.Errorf("Subdomain does not exist")
		},
	}

	infoProviderFactory := testInfoProviderFactory{infoProvider}

	editor := dnsEditor{
		infoProviderFactory: infoProviderFactory,
	}

	// act
	err := editor.DeleteSubdomain(domain, subdomain, recordType)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("DeleteSubdomain(%q, %q, %q) should return an error if the subdomain does not exist.", domain, subdomain, recordType)
	}
}

func Test_DeleteSubdomain_ValidParameters_SubdomainExists_DNSRecordDeleteFails_ErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	recordType := "AAAA"

	dnsClient := &testDNSClient{
		destroyRecordFunc: func(domain string, id string) error {
			return fmt.Errorf("Record deletion failed")
		},
	}

	infoProvider := &testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (record dnsimple.Record, err error) {
			return dnsimple.Record{}, nil
		},
	}

	dnsClientFactory := testDNSClientFactory{dnsClient}
	infoProviderFactory := testInfoProviderFactory{infoProvider}

	editor := dnsEditor{
		clientFactory:       dnsClientFactory,
		infoProviderFactory: infoProviderFactory,
	}

	// act
	err := editor.DeleteSubdomain(domain, subdomain, recordType)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("DeleteSubdomain(%q, %q, %q) should return an error of the record deletion failed at the DNS client.", domain, subdomain, recordType)
	}
}

func Test_DeleteSubdomain_ValidParameters_SubdomainExists_DNSRecordDeleteSucceeds_NoErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	recordType := "AAAA"

	dnsClient := &testDNSClient{
		destroyRecordFunc: func(domain string, id string) error {
			return nil
		},
	}

	infoProvider := &testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (record dnsimple.Record, err error) {
			return dnsimple.Record{}, nil
		},
	}

	dnsClientFactory := testDNSClientFactory{dnsClient}
	infoProviderFactory := testInfoProviderFactory{infoProvider}

	editor := dnsEditor{
		clientFactory:       dnsClientFactory,
		infoProviderFactory: infoProviderFactory,
	}

	// act
	err := editor.DeleteSubdomain(domain, subdomain, recordType)

	// assert
	if err != nil {
		t.Fail()
		t.Logf("DeleteSubdomain(%q, %q, %q) should not return an error if the DNS record deletion succeeds.", domain, subdomain, recordType)
	}
}
