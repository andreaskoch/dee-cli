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

// If any of the given parameters is invalid UpdateSubdomain should respond with an error.
func Test_UpdateSubdomain_ParametersInvalid_ErrorIsReturned(t *testing.T) {
	// arrange
	inputs := []struct {
		domain    string
		subdomain string
		ip        net.IP
	}{
		{"example.com", "", net.ParseIP("::1")},
		{"www", "", net.ParseIP("::1")},
		{"", "", net.ParseIP("::1")},
		{" ", " ", net.ParseIP("::1")},
		{"example.com", "www", nil},
	}
	updater := dnsimpleUpdater{}

	for _, input := range inputs {

		// act
		err := updater.UpdateSubdomain(input.domain, input.subdomain, input.ip)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("UpdateSubdomain(%q, %q, %q) should return an error.", input.domain, input.subdomain, input.ip)
		}
	}
}

// UpdateSubdomain should return an error if the given subdomain does not exist.
func Test_UpdateSubdomain_ValidParameters_SubdomainNotFound_ErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	ip := net.ParseIP("::1")

	updater := dnsimpleUpdater{
		infoProvider: &testDNSInfoProvider{
			getSubdomainRecordFunc: func(domain, subdomain string) (record dnsimple.Record, err error) {
				return dnsimple.Record{}, fmt.Errorf("")
			},
		},
	}

	// act
	err := updater.UpdateSubdomain(domain, subdomain, ip)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("UpdateSubdomain(%q, %q, %q) should return an error if the subdomain does not exist.", domain, subdomain, ip)
	}
}

func Test_UpdateSubdomain_ValidParameters_SubdomainExists_DNSRecordUpdateFails_ErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	ip := net.ParseIP("::1")

	dnsClient := &testDNSClient{
		updateRecordFunc: func(domain string, id string, opts *dnsimple.ChangeRecord) (string, error) {
			return "", fmt.Errorf("Record update failed")
		},
	}

	infoProvider := &testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain string) (record dnsimple.Record, err error) {
			return dnsimple.Record{}, nil
		},
	}

	updater := dnsimpleUpdater{
		client:       dnsClient,
		infoProvider: infoProvider,
	}

	// act
	err := updater.UpdateSubdomain(domain, subdomain, ip)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("UpdateSubdomain(%q, %q, %q) should return an error of the record update failed at the DNS client.", domain, subdomain, ip)
	}
}

func Test_UpdateSubdomain_ValidParameters_SubdomainExists_DNSRecordUpdateSucceeds_NoErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	ip := net.ParseIP("::1")

	dnsClient := &testDNSClient{
		updateRecordFunc: func(domain string, id string, opts *dnsimple.ChangeRecord) (string, error) {
			return "", nil
		},
	}

	infoProvider := &testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain string) (record dnsimple.Record, err error) {
			return dnsimple.Record{}, nil
		},
	}

	updater := dnsimpleUpdater{
		client:       dnsClient,
		infoProvider: infoProvider,
	}

	// act
	err := updater.UpdateSubdomain(domain, subdomain, ip)

	// assert
	if err != nil {
		t.Fail()
		t.Logf("UpdateSubdomain(%q, %q, %q) should not return an error if the DNS record update succeeds.", domain, subdomain, ip)
	}
}

// If the update will not change the IP the update is aborted and an error is returned.
func Test_UpdateSubdomain_ValidParameters_SubdomainExists_ExistingIPIsTheSame_ErrorIsReturned(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	ip := net.ParseIP("::1")

	dnsClient := &testDNSClient{
		updateRecordFunc: func(domain string, id string, opts *dnsimple.ChangeRecord) (string, error) {
			return "", nil
		},
	}

	existingRecord := dnsimple.Record{
		Name:       "example.com",
		Content:    "::1",
		RecordType: "AAAA",
		Ttl:        600,
	}

	infoProvider := &testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain string) (record dnsimple.Record, err error) {
			return existingRecord, nil
		},
	}

	updater := dnsimpleUpdater{
		client:       dnsClient,
		infoProvider: infoProvider,
	}

	// act
	err := updater.UpdateSubdomain(domain, subdomain, ip)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("UpdateSubdomain(%q, %q, %q) should return an error because the IP of the existing record is the same as in the update.", domain, subdomain, ip)
	}
}

func Test_UpdateSubdomain_ValidParameters_SubdomainExists_OnlyTheIPIsChangedOnTheDNSRecord(t *testing.T) {
	// arrange
	domain := "example.com"
	subdomain := "www"
	ip := net.ParseIP("::2")

	existingRecord := dnsimple.Record{
		Name:       "example.com",
		Content:    "::1",
		RecordType: "AAAA",
		Ttl:        600,
	}

	dnsClient := &testDNSClient{
		updateRecordFunc: func(domain string, id string, opts *dnsimple.ChangeRecord) (string, error) {

			// assert
			if opts.Name != existingRecord.Name {
				t.Fail()
				t.Logf("The DNS name should not change during an update (Old: %q, New: %q)", existingRecord.Name, opts.Name)
			}

			if opts.Type != existingRecord.RecordType {
				t.Fail()
				t.Logf("The DNS record type should not change during an update (Old: %q, New: %q)", existingRecord.RecordType, opts.Type)
			}

			if opts.Ttl != fmt.Sprintf("%d", existingRecord.Ttl) {
				t.Fail()
				t.Logf("The DNS record TTL should not change during an update (Old: %q, New: %q)", existingRecord.Ttl, opts.Ttl)
			}

			if opts.Value != ip.String() {
				t.Fail()
				t.Logf("The DNS record value should have changed to %q", ip.String())
			}

			return "", nil
		},
	}

	infoProvider := &testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain string) (record dnsimple.Record, err error) {
			return existingRecord, nil
		},
	}

	updater := dnsimpleUpdater{
		client:       dnsClient,
		infoProvider: infoProvider,
	}

	// act
	updater.UpdateSubdomain(domain, subdomain, ip)
}

// testDNSInfoProvider is a DNS info-provider used for testing.
type testDNSInfoProvider struct {
	getSubdomainRecordFunc func(domain, subdomain string) (record dnsimple.Record, err error)
}

func (infoProvider *testDNSInfoProvider) GetSubdomainRecord(domain, subdomain string) (record dnsimple.Record, err error) {
	return infoProvider.getSubdomainRecordFunc(domain, subdomain)
}
