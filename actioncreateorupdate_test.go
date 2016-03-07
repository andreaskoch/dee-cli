// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/pearkes/dnsimple"
	"net"
	"strings"
	"testing"
)

func Test_createOrUpdateAction_Name_CorrectActionNameIsReturned(t *testing.T) {

	// arrange
	createOrUpdateAction := createOrUpdateAction{}

	// act
	result := createOrUpdateAction.Name()

	// assert
	if result != "createorupdate" {
		t.Fail()
		t.Logf("createOrUpdateAction.Name() should have returned %q but returned %q instead.", "createorupdate", result)
	}

}

func Test_createOrUpdateAction_Description_ResultIsNotEmpty(t *testing.T) {

	// arrange
	createOrUpdateAction := createOrUpdateAction{}

	// act
	result := createOrUpdateAction.Description()

	// assert
	if isEmpty(result) {
		t.Fail()
		t.Logf("createOrUpdateAction.Description() not be empty.")
	}

}

func Test_createOrUpdateAction_Usage_ResultIsNotEmpty(t *testing.T) {

	// arrange
	createOrUpdateAction := createOrUpdateAction{}

	// act
	result := createOrUpdateAction.Usage()

	// assert
	if isEmpty(result) {
		t.Fail()
		t.Logf("createOrUpdateAction.Usage() not be empty.")
	}

}

// createOrUpdateAction.Execute should return an error if the arguments are invalid.
func Test_createOrUpdateAction_InvalidArguments_ErrorIsReturned(t *testing.T) {
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

	dnsCreator := &testDNSEditor{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			return nil
		},
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			return nil
		},
	}

	editorFactory := testDNSEditorFactory{dnsCreator, nil}

	dnsInfoProvider := testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (dnsimple.Record, error) {
			return dnsimple.Record{}, nil
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	createOrUpdateAction := createOrUpdateAction{editorFactory, infoProviderFactory, nil}

	for _, arguments := range validArgumentsSet {

		// act
		_, err := createOrUpdateAction.Execute(arguments)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("createOrUpdateAction.Execute(%q) should return an error", arguments)
		}
	}
}

// createOrUpdateAction.Execute should return an error if the argument values are invalid.
func Test_createOrUpdateAction_InvalidArgumentValues_ErrorIsReturned(t *testing.T) {
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

	dnsCreator := &testDNSEditor{
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			return nil
		},
	}

	editorFactory := testDNSEditorFactory{dnsCreator, nil}

	dnsInfoProvider := testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (dnsimple.Record, error) {
			return dnsimple.Record{}, nil
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	createOrUpdateAction := createOrUpdateAction{editorFactory, infoProviderFactory, nil}

	for _, arguments := range validArgumentsSet {

		// act
		_, err := createOrUpdateAction.Execute(arguments)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("createOrUpdateAction.Execute(%q) should return an error", arguments)
		}
	}
}

// createOrUpdateAction.Execute should return an error if the given IP address is invalid.
func Test_createOrUpdateAction_IPAddressIsInvalid_ErrorIsReturned(t *testing.T) {
	// arrange
	invalidIPs := []string{
		"127.0",
		"127.0.0.1.1",
		"255.255.255.255.255",
		"10000:21323:31231",
	}

	dnsCreator := &testDNSEditor{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			return nil
		},
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			return nil
		},
	}

	editorFactory := testDNSEditorFactory{dnsCreator, nil}

	dnsInfoProvider := testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (dnsimple.Record, error) {
			return dnsimple.Record{}, nil
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	createOrUpdateAction := createOrUpdateAction{editorFactory, infoProviderFactory, nil}

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
		_, err := createOrUpdateAction.Execute(arguments)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("createOrUpdateAction.Execute(%q) should return an error because the given IP address (%q) is invalid", arguments, invalidIP)
		}
	}
}

// createOrUpdateAction.Execute should not return an error if the arguments are valid and the subdomain createOrUpdateAction.Execute succeeds.
func Test_createOrUpdateAction_ValidArguments_NoErrorIsReturned(t *testing.T) {
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

	dnsCreator := &testDNSEditor{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			return nil
		},
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			return nil
		},
	}

	editorFactory := testDNSEditorFactory{dnsCreator, nil}

	dnsInfoProvider := testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (dnsimple.Record, error) {
			return dnsimple.Record{}, nil
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	createOrUpdateAction := createOrUpdateAction{editorFactory, infoProviderFactory, nil}

	for _, arguments := range validArgumentsSet {

		// act
		_, err := createOrUpdateAction.Execute(arguments)

		// assert
		if err != nil {
			t.Fail()
			t.Logf("createOrUpdateAction.Execute(%q) should not return an error: %q", arguments, err.Error())
		}
	}
}

// createOrUpdateAction.Execute should return an error if no info provider is supplied
func Test_createOrUpdateAction_ValidArguments_NoInfoProvider_ErrorIsReturned(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-ip",
		"2001:0db8:0000:0042:0000:8a2e:0370:7334",
	}

	dnsCreator := &testDNSEditor{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			return nil
		},
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			return fmt.Errorf("DNS Record create failed")
		},
	}

	editorFactory := testDNSEditorFactory{dnsCreator, nil}

	createOrUpdateAction := createOrUpdateAction{editorFactory, nil, nil}

	// act
	_, err := createOrUpdateAction.Execute(arguments)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("createOrUpdateAction.Execute(%q) should return an error because no info provider was given.", arguments)
	}
}

// createOrUpdateAction.Execute should return an error if the DNS editor responds with one.
func Test_createOrUpdateAction_ValidArguments_CreateFails_ErrorIsReturned(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-ip",
		"2001:0db8:0000:0042:0000:8a2e:0370:7334",
	}

	dnsCreator := &testDNSEditor{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			return nil
		},
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			return fmt.Errorf("DNS Record create failed")
		},
	}

	editorFactory := testDNSEditorFactory{dnsCreator, nil}

	dnsInfoProvider := testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (dnsimple.Record, error) {
			return dnsimple.Record{}, fmt.Errorf("No found")
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	createOrUpdateAction := createOrUpdateAction{editorFactory, infoProviderFactory, nil}

	// act
	_, err := createOrUpdateAction.Execute(arguments)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("createOrUpdateAction.Execute(%q) should return an error because the DNS creator returned one.", arguments)
	}
}

// createOrUpdateAction.Execute should try an update if the record already exists.
func Test_createOrUpdateAction_ValidArguments_RecordAlreadyExists_UpdateIsExecuted(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-ip",
		"2001:0db8:0000:0042:0000:8a2e:0370:7334",
	}

	updateWasCalled := false
	dnsCreator := &testDNSEditor{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			updateWasCalled = true
			return nil
		},
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			t.Fail()
			return nil
		},
	}

	editorFactory := testDNSEditorFactory{dnsCreator, nil}

	dnsInfoProvider := testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (dnsimple.Record, error) {
			return dnsimple.Record{
				Name:       "www",
				Content:    "2001:0db8:0000:0042:0000:8a2e:0370:7334",
				RecordType: "AAAA",
			}, nil
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	createOrUpdateAction := createOrUpdateAction{editorFactory, infoProviderFactory, nil}

	// act
	createOrUpdateAction.Execute(arguments)

	// assert
	if updateWasCalled == false {
		t.Fail()
		t.Logf("createOrUpdateAction.Execute(%q) should have triggered the update action.", arguments)
	}
}

// createOrUpdateAction.Execute should return a success message if the DNS creator succeeds.
func Test_createOrUpdateAction_ValidArguments_CreateSucceds_SuccessMessageIsReturned(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-ip",
		"2001:0db8:0000:0042:0000:8a2e:0370:7334",
	}

	dnsCreator := &testDNSEditor{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			return nil
		},
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			return nil
		},
	}

	editorFactory := testDNSEditorFactory{dnsCreator, nil}

	dnsInfoProvider := testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (dnsimple.Record, error) {
			return dnsimple.Record{}, nil
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	createOrUpdateAction := createOrUpdateAction{editorFactory, infoProviderFactory, nil}

	// act
	response, _ := createOrUpdateAction.Execute(arguments)

	// assert
	if response == nil {
		t.Fail()
		t.Logf("createOrUpdateAction.Execute(%q) should respond with a success message if the DNS creator succeeds.", arguments)
	}
}

// createOrUpdateAction.Execute should return an error if the DNS editor factory returns an error.
func Test_createOrUpdateAction_ValidArguments_DNSEditorCreationFails_ErrorIsReturned(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-ip",
		"2001:db8:0:42:0:8a2e:370:7334",
	}

	editorFactory := testDNSEditorFactory{nil, fmt.Errorf("Unable to create DNS editor")}

	dnsInfoProvider := testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (dnsimple.Record, error) {
			return dnsimple.Record{}, fmt.Errorf("Record does not exist")
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	createOrUpdateAction := createOrUpdateAction{editorFactory, infoProviderFactory, nil}

	// act
	_, err := createOrUpdateAction.Execute(arguments)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("createOrUpdateAction.Execute(%q) should return an error if the DNS editor factory returned an error.", arguments)
	}
}

// createOrUpdateAction.Execute should return a success message that contains the subdomain, domain and ip.
func Test_createOrUpdateAction_ValidArguments_UpdateSucceds_SuccessMessageContainsTheSubdomainAndNewIP(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-ip",
		"2001:db8:0:42:0:8a2e:370:7334",
	}

	dnsCreator := &testDNSEditor{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			return nil
		},
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			return nil
		},
	}

	editorFactory := testDNSEditorFactory{dnsCreator, nil}

	dnsInfoProvider := testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (dnsimple.Record, error) {
			return dnsimple.Record{}, nil
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	createOrUpdateAction := createOrUpdateAction{editorFactory, infoProviderFactory, nil}

	// act
	response, _ := createOrUpdateAction.Execute(arguments)

	// assert
	containsIP := strings.Contains(response.Text(), "2001:db8:0:42:0:8a2e:370:7334")
	containsSubdomain := strings.Contains(response.Text(), "www")
	containsDomain := strings.Contains(response.Text(), "example.com")

	if !containsIP || !containsSubdomain || !containsDomain {
		t.Fail()
		t.Logf("createOrUpdateAction.Execute(%q) should respond with a success message that contains the domain, subdomain and ip but responded with %q instead.", arguments, response.Text())
	}
}

// createOrUpdateAction.Execute should return an error if the update fails
func Test_createOrUpdateAction_ValidArguments_UpdateFails_ErrorIsReturned(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-ip",
		"2001:db8:0:42:0:8a2e:370:7334",
	}

	dnsCreator := &testDNSEditor{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			return fmt.Errorf("Create failed")
		},
	}

	editorFactory := testDNSEditorFactory{dnsCreator, nil}

	dnsInfoProvider := testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (dnsimple.Record, error) {
			return dnsimple.Record{}, nil
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	createOrUpdateAction := createOrUpdateAction{editorFactory, infoProviderFactory, nil}

	// act
	_, err := createOrUpdateAction.Execute(arguments)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("createOrUpdateAction.Execute(%q) should return an error if the update failed.", arguments)
	}
}

// createOrUpdateAction.Execute should return a success message that contains the subdomain, domain and ip.
func Test_createOrUpdateAction_ValidArguments_CreateSucceds_SuccessMessageContainsTheSubdomainAndNewIP(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-ip",
		"2001:db8:0:42:0:8a2e:370:7334",
	}

	dnsCreator := &testDNSEditor{
		updateSubdomainFunc: func(domain, subdomain string, ip net.IP) error {
			return nil
		},
		createSubdomainFunc: func(domain, subdomain string, timeToLive int, ip net.IP) error {
			return nil
		},
	}

	editorFactory := testDNSEditorFactory{dnsCreator, nil}

	dnsInfoProvider := testDNSInfoProvider{
		getSubdomainRecordFunc: func(domain, subdomain, recordType string) (dnsimple.Record, error) {
			return dnsimple.Record{}, fmt.Errorf("Record does not exist")
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	createOrUpdateAction := createOrUpdateAction{editorFactory, infoProviderFactory, nil}

	// act
	response, _ := createOrUpdateAction.Execute(arguments)

	// assert
	containsIP := strings.Contains(response.Text(), "2001:db8:0:42:0:8a2e:370:7334")
	containsSubdomain := strings.Contains(response.Text(), "www")
	containsDomain := strings.Contains(response.Text(), "example.com")

	if !containsIP || !containsSubdomain || !containsDomain {
		t.Fail()
		t.Logf("createOrUpdateAction.Execute(%q) should respond with a success message that contains the domain, subdomain and ip but responded with %q instead.", arguments, response.Text())
	}
}
