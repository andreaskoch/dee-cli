// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

func Test_updateAction_Name_UpdateIsReturned(t *testing.T) {

	// arrange
	updateAction := updateAction{}

	// act
	result := updateAction.Name()

	// assert
	if result != "update" {
		t.Fail()
		t.Logf("updateAction.Name() should have returned %q but returned %q instead.", "update", result)
	}

}

func Test_updateAction_Description_ResultIsNotEmpty(t *testing.T) {

	// arrange
	updateAction := updateAction{}

	// act
	result := updateAction.Description()

	// assert
	if isEmpty(result) {
		t.Fail()
		t.Logf("updateAction.Description() not be empty.")
	}

}

func Test_updateAction_Usage_ResultIsNotEmpty(t *testing.T) {

	// arrange
	updateAction := updateAction{}

	// act
	result := updateAction.Usage()

	// assert
	if isEmpty(result) {
		t.Fail()
		t.Logf("updateAction.Usage() not be empty.")
	}

}

// updateAction.Execute should return an error if the arguments are invalid.
func Test_updateAction_InvalidArguments_ErrorIsReturned(t *testing.T) {
	// arrange
	validArgumentsSet := [][]string{
		{
			"-domainname",
			"example.com",
			"-subdomain",
			"www",
			"-ip",
			"127.0.0.1",
		},
		{
			"-domainname",
			"example.com",
			"-subdomain",
			"www",
			"-ip",
			"",
		},
		{
			"-domainname",
			"example.com",
			"-subdomain",
			"www",
			"-ip",
			"dasjdksalkdjsakljdklsa",
		},
		{
			"-domainname",
			"example.com",
			"-subdomain",
			"www",
			"-ip",
			"4312908321098392108309",
		},
		{
			"-domain",
			"example.com",
			"-sub-domain",
			"a.b.c",
			"-ip",
			"127.0.0.1",
		},
		{
			"-domain",
			"example.com",
			"-subdomain",
			"www",
			"-ipaddress",
			"2001:0db8:0000:0042:0000:8a2e:0370:7334",
		},
		{
			"-domain",
			"example.co.uk",
			"-subdomain",
			"www",
			"-IP",
			"2001:0db8:0000:0042:0000:8a2e:0370:7334",
		},
		{
			"--DOMAIN",
			"example.com",
			"--subdomain",
			"www",
			"--ip",
			"127.0.0.1",
		},
		{
			"--domain = example.com",
			"--subdomain=www",
			"--ip=127.0.0.1",
		},
	}

	dnsUpdater := &testDNSUpdater{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			return nil
		},
	}

	updateAction := updateAction{dnsUpdater, nil}

	for _, arguments := range validArgumentsSet {

		// act
		_, err := updateAction.Execute(arguments)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("updateAction.Execute(%q) should return an error", arguments)
		}
	}
}

// updateAction.Execute should return an error if the arguments are valid but empty.
func Test_updateAction_ValidArgumentsButEmpty_ErrorIsReturned(t *testing.T) {
	// arrange
	validArgumentsSet := [][]string{
		{
			"-domain",
			"",
			"-subdomain",
			"www",
			"-ip",
			"127.0.0.1",
		},
		{
			"-domain",
			"example.com",
			"-subdomain",
			"www",
			"-ip",
			"",
		},
	}

	dnsUpdater := &testDNSUpdater{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			return nil
		},
	}

	updateAction := updateAction{dnsUpdater, nil}

	for _, arguments := range validArgumentsSet {

		// act
		_, err := updateAction.Execute(arguments)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("updateAction.Execute(%q) should return an error", arguments)
		}
	}
}

// updateAction.Execute should return an error if the given IP address is invalid.
func Test_updateAction_IPAddressIsInvalid_ErrorIsReturned(t *testing.T) {
	// arrange
	invalidIPs := []string{
		"127.0",
		"127.0.0.1.1",
		"255.255.255.255.255",
		"10000:21323:31231",
	}

	dnsUpdater := &testDNSUpdater{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			return nil
		},
	}

	updateAction := updateAction{dnsUpdater, nil}

	for _, invalidIP := range invalidIPs {

		// act
		arguments := []string{
			"-domain",
			"example.com",
			"-subdomain",
			"www",
			"-ip",
			invalidIP,
		}
		_, err := updateAction.Execute(arguments)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("updateAction.Execute(%q) should return an error because the given IP address (%q) is invalid", arguments, invalidIP)
		}
	}
}

// updateAction.Execute should not return an error if the arguments are valid and the subdomain updateAction.Execute succeeds.
func Test_updateAction_ValidArguments_NoErrorIsReturned(t *testing.T) {
	// arrange
	validArgumentsSet := [][]string{
		{
			"-domain",
			"example.com",
			"-subdomain",
			"www",
			"-ip",
			"127.0.0.1",
		},
		{
			"-domain",
			"example.com",
			"-subdomain",
			"a.b.c",
			"-ip",
			"127.0.0.1",
		},
		{
			"-domain",
			"example.com",
			"-subdomain",
			"www",
			"-ip",
			"2001:0db8:0000:0042:0000:8a2e:0370:7334",
		},
		{
			"-domain",
			"example.co.uk",
			"-subdomain",
			"www",
			"-ip",
			"2001:0db8:0000:0042:0000:8a2e:0370:7334",
		},
		{
			"--domain",
			"example.com",
			"--subdomain",
			"www",
			"--ip",
			"127.0.0.1",
		},
		{
			"--domain=example.com",
			"--subdomain=www",
			"--ip=127.0.0.1",
		},
	}

	dnsUpdater := &testDNSUpdater{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			return nil
		},
	}

	updateAction := updateAction{dnsUpdater, nil}

	for _, arguments := range validArgumentsSet {

		// act
		_, err := updateAction.Execute(arguments)

		// assert
		if err != nil {
			t.Fail()
			t.Logf("updateAction.Execute(dnsUpdater, %q) should not return an error: %q", arguments, err.Error())
		}
	}
}

// updateAction.Execute should return an error if the DNS Updater responds with one.
func Test_updateAction_ValidArguments_DNSUpdaterRespondsWithError_ErrorIsReturned(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-ip",
		"2001:0db8:0000:0042:0000:8a2e:0370:7334",
	}

	dnsUpdater := &testDNSUpdater{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			return fmt.Errorf("DNS Record Update failed")
		},
	}

	updateAction := updateAction{dnsUpdater, nil}

	// act
	_, err := updateAction.Execute(arguments)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("updateAction.Execute(dnsUpdater, %q) should return an error because the DNS updater returned one.", arguments)
	}
}

// updateAction.Execute should return a success message if the DNS updater succeeds.
func Test_updateAction_ValidArguments_DNSUpdaterSucceeds_SuccessMessageIsReturned(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-ip",
		"2001:0db8:0000:0042:0000:8a2e:0370:7334",
	}

	dnsUpdater := &testDNSUpdater{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			return nil
		},
	}

	updateAction := updateAction{dnsUpdater, nil}

	// act
	response, _ := updateAction.Execute(arguments)

	// assert
	if response == nil {
		t.Fail()
		t.Logf("updateAction.Execute(dnsUpdater, %q) should respond with a success message if the DNS updater succeeds.", arguments)
	}
}

// updateAction.Execute should return a success message that contains the subdomain, domain and ip.
func Test_updateAction_ValidArguments_DNSUpdaterSucceeds_SuccessMessageContainsTheSubdomainAndNewIP(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-ip",
		"2001:db8:0:42:0:8a2e:370:7334",
	}

	dnsUpdater := &testDNSUpdater{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			return nil
		},
	}

	updateAction := updateAction{dnsUpdater, nil}

	// act
	response, _ := updateAction.Execute(arguments)

	// assert
	containsIP := strings.Contains(response.Text(), "2001:db8:0:42:0:8a2e:370:7334")
	containsSubdomain := strings.Contains(response.Text(), "www")
	containsDomain := strings.Contains(response.Text(), "example.com")

	if !containsIP || !containsSubdomain || !containsDomain {
		t.Fail()
		t.Logf("updateAction.Execute(dnsUpdater, %q) should respond with a success message that contains the domain, subdomain and ip but responded with %q instead.", arguments, response.Text())
	}
}
