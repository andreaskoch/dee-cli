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

func Test_createAction_Name_UpdateIsReturned(t *testing.T) {

	// arrange
	createAction := createAction{}

	// act
	result := createAction.Name()

	// assert
	if result != "create" {
		t.Fail()
		t.Logf("createAction.Name() should have returned %q but returned %q instead.", "create", result)
	}

}

func Test_createAction_Description_ResultIsNotEmpty(t *testing.T) {

	// arrange
	createAction := createAction{}

	// act
	result := createAction.Description()

	// assert
	if isEmpty(result) {
		t.Fail()
		t.Logf("createAction.Description() not be empty.")
	}

}

func Test_createAction_Usage_ResultIsNotEmpty(t *testing.T) {

	// arrange
	createAction := createAction{}

	// act
	result := createAction.Usage()

	// assert
	if isEmpty(result) {
		t.Fail()
		t.Logf("createAction.Usage() not be empty.")
	}

}

// createAction.Execute should return an error if the arguments are invalid.
func Test_createAction_InvalidArguments_ErrorIsReturned(t *testing.T) {
	// arrange
	validArgumentsSet := [][]string{
		{
			"-domainname",
			"example.com",
			"-subdomain",
			"www",
			"-ttl",
			"3600",
			"-ip",
			"127.0.0.1",
		},
		{
			"-domain",
			"example.com",
			"-subdomain",
			"www",
			"-ttl",
			"3600",
			"-ip",
			"",
		},
		{
			"-domain",
			"example.com",
			"-subdomain",
			"www",
			"-timetolive",
			"700",
			"-ip",
			"dsadklsajdklsakldj",
		},
		{
			"-domain",
			"example.com",
			"-subdomain",
			"www",
			"-ip",
			"3123213123213213213213",
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

	dnsCreator := &testDNSCreator{
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			return nil
		},
	}

	createAction := createAction{dnsCreator, nil}

	for _, arguments := range validArgumentsSet {

		// act
		_, err := createAction.Execute(arguments)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("createAction.Execute(%q) should return an error", arguments)
		}
	}
}

// createAction.Execute should return an error if the argument values are invalid.
func Test_createAction_InvalidArgumentValues_ErrorIsReturned(t *testing.T) {
	// arrange
	validArgumentsSet := [][]string{
		{
			"-domain",
			"",
			"-subdomain",
			"www",
			"-ttl",
			"600",
			"-ip",
			"127.0.0.1",
		},
		{
			"-domain",
			"example.com",
			"-subdomain",
			"",
			"-ttl",
			"600",
			"-ip",
			"127.0.0.1",
		},
		{
			"-domain",
			"example.com",
			"-subdomain",
			"www",
			"-ttl",
			"-600",
			"-ip",
			"127.0.0.1",
		},
		{
			"-domain",
			"example.com",
			"-subdomain",
			"www",
			"-ttl",
			"One Hour",
			"-ip",
			"127.0.0.1",
		},
		{
			"-domain",
			"example.com",
			"-subdomain",
			"www",
			"-ttl",
			"600",
			"-ip",
			"",
		},
	}

	dnsCreator := &testDNSCreator{
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			return nil
		},
	}

	createAction := createAction{dnsCreator, nil}

	for _, arguments := range validArgumentsSet {

		// act
		_, err := createAction.Execute(arguments)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("createAction.Execute(%q) should return an error", arguments)
		}
	}
}

// createAction.Execute should return an error if the given IP address is invalid.
func Test_createAction_IPAddressIsInvalid_ErrorIsReturned(t *testing.T) {
	// arrange
	invalidIPs := []string{
		"127.0",
		"127.0.0.1.1",
		"255.255.255.255.255",
		"10000:21323:31231",
	}

	dnsCreator := &testDNSCreator{
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			return nil
		},
	}

	createAction := createAction{dnsCreator, nil}

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
		_, err := createAction.Execute(arguments)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("createAction.Execute(%q) should return an error because the given IP address (%q) is invalid", arguments, invalidIP)
		}
	}
}

// createAction.Execute should not return an error if the arguments are valid and the subdomain createAction.Execute succeeds.
func Test_createAction_ValidArguments_NoErrorIsReturned(t *testing.T) {
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

	dnsCreator := &testDNSCreator{
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			return nil
		},
	}

	createAction := createAction{dnsCreator, nil}

	for _, arguments := range validArgumentsSet {

		// act
		_, err := createAction.Execute(arguments)

		// assert
		if err != nil {
			t.Fail()
			t.Logf("createAction.Execute(dnsCreator, %q) should not return an error: %q", arguments, err.Error())
		}
	}
}

// createAction.Execute should return an error if the DNS creator responds with one.
func Test_createAction_ValidArguments_DNSCreatorRespondsWithError_ErrorIsReturned(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-ip",
		"2001:0db8:0000:0042:0000:8a2e:0370:7334",
	}

	dnsCreator := &testDNSCreator{
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			return fmt.Errorf("DNS Record create failed")
		},
	}

	createAction := createAction{dnsCreator, nil}

	// act
	_, err := createAction.Execute(arguments)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("createAction.Execute(dnsCreator, %q) should return an error because the DNS creator returned one.", arguments)
	}
}

// createAction.Execute should return a success message if the DNS creator succeeds.
func Test_createAction_ValidArguments_DNSCreatorSucceeds_SuccessMessageIsReturned(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-ip",
		"2001:0db8:0000:0042:0000:8a2e:0370:7334",
	}

	dnsCreator := &testDNSCreator{
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			return nil
		},
	}

	createAction := createAction{dnsCreator, nil}

	// act
	response, _ := createAction.Execute(arguments)

	// assert
	if response == nil {
		t.Fail()
		t.Logf("createAction.Execute(dnsCreator, %q) should respond with a success message if the DNS creator succeeds.", arguments)
	}
}

// createAction.Execute should return a success message that contains the subdomain, domain and ip.
func Test_createAction_ValidArguments_DNSCreatorSucceeds_SuccessMessageContainsTheSubdomainAndNewIP(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-ip",
		"2001:db8:0:42:0:8a2e:370:7334",
	}

	dnsCreator := &testDNSCreator{
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			return nil
		},
	}

	createAction := createAction{dnsCreator, nil}

	// act
	response, _ := createAction.Execute(arguments)

	// assert
	containsIP := strings.Contains(response.Text(), "2001:db8:0:42:0:8a2e:370:7334")
	containsSubdomain := strings.Contains(response.Text(), "www")
	containsDomain := strings.Contains(response.Text(), "example.com")

	if !containsIP || !containsSubdomain || !containsDomain {
		t.Fail()
		t.Logf("createAction.Execute(dnsCreator, %q) should respond with a success message that contains the domain, subdomain and ip but responded with %q instead.", arguments, response.Text())
	}
}
