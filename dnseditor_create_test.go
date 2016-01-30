// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/pearkes/dnsimple"
	"net"
	"testing"
)

// testDNSCreator creates DNS records.
type testDNSCreator struct {
	createSubdomainFunc func(domain, subdomain string, timeToLive int, ip net.IP) error
}

func (editor *testDNSCreator) CreateSubdomain(domain, subdomain string, timeToLive int, ip net.IP) error {
	return editor.createSubdomainFunc(domain, subdomain, timeToLive, ip)
}

// If any of the given parameters is invalid CreateSubdomain should respond with an error.
func Test_CreateSubdomain_ParametersInvalid_ErrorIsReturned(t *testing.T) {
	// arrange
	inputs := []struct {
		domain    string
		subdomain string
		ttl       int
		ip        net.IP
	}{
		{"example.com", "", 600, net.ParseIP("::1")},
		{"www", "", 600, net.ParseIP("::1")},
		{"", "", 600, net.ParseIP("::1")},
		{" ", " ", 600, net.ParseIP("::1")},
		{"example.com", "www", 600, nil},
	}
	editor := dnsEditor{}

	for _, input := range inputs {

		// act
		err := editor.CreateSubdomain(input.domain, input.subdomain, input.ttl, input.ip)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("CreateSubdomain(%q, %q, %q, %q) should return an error.", input.domain, input.subdomain, input.ttl, input.ip)
		}
	}
}

// CreateSubdomain should return an error if the given subdomain does not exist.
func Test_CreateSubdomain_ValidParameters_SubdomainNotFound_ErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	ttl := 600
	ip := net.ParseIP("::1")

	infoProvider := &testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (record dnsimple.Record, err error) {
			return dnsimple.Record{}, fmt.Errorf("")
		},
	}

	infoProviderFactory := testInfoProviderFactory{infoProvider}

	editor := dnsEditor{
		infoProviderFactory: infoProviderFactory,
	}

	// act
	err := editor.CreateSubdomain(domain, subdomain, ttl, ip)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("CreateSubdomain(%q, %q, %q, %q) should return an error if the subdomain does not exist.", domain, subdomain, ttl, ip)
	}
}

func Test_CreateSubdomain_ValidParameters_SubdomainExists_DNSRecordCreationFails_ErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	ttl := 600
	ip := net.ParseIP("::1")

	dnsClient := &testDNSClient{
		createRecordFunc: func(domain string, opts *dnsimple.ChangeRecord) (string, error) {
			return "", fmt.Errorf("Failed to create record")
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
	err := editor.CreateSubdomain(domain, subdomain, ttl, ip)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("CreateSubdomain(%q, %q, %q, %q) should return an error of the record creation failed at the DNS client.", domain, subdomain, ip, ttl)
	}
}

func Test_CreateSubdomain_ValidParameters_SubdomainExists_DNSRecordCreationSucceeds_NoErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	ttl := 3600
	ip := net.ParseIP("::1")

	dnsClient := &testDNSClient{
		createRecordFunc: func(domain string, opts *dnsimple.ChangeRecord) (string, error) {
			return "", nil
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
	err := editor.CreateSubdomain(domain, subdomain, ttl, ip)

	// assert
	if err != nil {
		t.Fail()
		t.Logf("CreateSubdomain(%q, %q, %q, %q) should not return an error if the DNS record creation succeeded.", domain, subdomain, ttl, ip)
	}
}
