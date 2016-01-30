// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
	"testing"
)

func Test_deleteAction_Name_UpdateIsReturned(t *testing.T) {

	// arrange
	deleteAction := deleteAction{}

	// act
	result := deleteAction.Name()

	// assert
	if result != "delete" {
		t.Fail()
		t.Logf("deleteAction.Name() should have returned %q but returned %q instead.", "delete", result)
	}

}

func Test_deleteAction_Description_ResultIsNotEmpty(t *testing.T) {

	// arrange
	deleteAction := deleteAction{}

	// act
	result := deleteAction.Description()

	// assert
	if isEmpty(result) {
		t.Fail()
		t.Logf("deleteAction.Description() not be empty.")
	}

}

func Test_deleteAction_Usage_ResultIsNotEmpty(t *testing.T) {

	// arrange
	deleteAction := deleteAction{}

	// act
	result := deleteAction.Usage()

	// assert
	if isEmpty(result) {
		t.Fail()
		t.Logf("deleteAction.Usage() not be empty.")
	}

}

// deleteAction.Execute should return an error if the arguments are invalid.
func Test_deleteAction_InvalidArguments_ErrorIsReturned(t *testing.T) {
	// arrange
	validArgumentsSet := [][]string{
		{
			"-domainname",
			"example.com",
			"-subdomain",
			"www",
			"-type",
			"AAAA",
		},
		{
			"-domain",
			"example.com",
			"-subdomain-name",
			"www",
			"-type",
			"AAAA",
		},
		{
			"-domain",
			"example.com",
			"-subdomain",
			"",
			"-type",
			"AAAA",
		},
		{
			"-domain",
			"example.com",
			"-subdomain",
			"www",
			"-dnsrecordtype",
			"",
		},
	}

	dnsDeleter := &testDNSDeleter{
		deleteSubdomainFunc: func(domain, subdomain string, recordType string) error {
			return nil
		},
	}

	deleteAction := deleteAction{dnsDeleter}

	for _, arguments := range validArgumentsSet {

		// act
		_, err := deleteAction.Execute(arguments)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("deleteAction.Execute(%q) should return an error", arguments)
		}
	}
}

// deleteAction.Execute should return an error if the argument values are invalid.
func Test_deleteAction_InvalidArgumentValues_ErrorIsReturned(t *testing.T) {
	// arrange
	validArgumentsSet := [][]string{
		{
			"-domain",
			"",
			"-subdomain",
			"www",
			"-type",
			"AAAA",
		},
		{
			"-domain",
			"example.com",
			"-subdomain",
			"",
			"-type",
			"AAAA",
		},
		{
			"-domain",
			"example.com",
			"-subdomain",
			"www",
			"-type",
			"",
		},
	}

	dnsDeleter := &testDNSDeleter{
		deleteSubdomainFunc: func(domain, subdomain string, recordType string) error {
			return nil
		},
	}

	deleteAction := deleteAction{dnsDeleter}

	for _, arguments := range validArgumentsSet {

		// act
		_, err := deleteAction.Execute(arguments)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("deleteAction.Execute(%q) should return an error", arguments)
		}
	}
}

// deleteAction.Execute should not return an error if the arguments are valid and the subdomain deleteAction.Execute succeeds.
func Test_deleteAction_ValidArguments_NoErrorIsReturned(t *testing.T) {
	// arrange
	validArgumentsSet := [][]string{
		{
			"-domain",
			"example.com",
			"-subdomain",
			"www",
			"-type",
			"AAAA",
		},
		{
			"--domain=example.com",
			"--subdomain=www",
			"--type=AAAA",
		},
	}

	dnsDeleter := &testDNSDeleter{
		deleteSubdomainFunc: func(domain, subdomain string, recordType string) error {
			return nil
		},
	}

	deleteAction := deleteAction{dnsDeleter}

	for _, arguments := range validArgumentsSet {

		// act
		_, err := deleteAction.Execute(arguments)

		// assert
		if err != nil {
			t.Fail()
			t.Logf("deleteAction.Execute(dnsDeleter, %q) should not return an error: %q", arguments, err.Error())
		}
	}
}

// deleteAction.Execute should return an error if the DNS Updater responds with one.
func Test_deleteAction_ValidArguments_DNSCreatorRespondsWithError_ErrorIsReturned(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-type",
		"AAAA",
	}

	dnsDeleter := &testDNSDeleter{
		deleteSubdomainFunc: func(domain, subdomain string, recordType string) error {
			return fmt.Errorf("DNS Record delete failed")
		},
	}

	deleteAction := deleteAction{dnsDeleter}

	// act
	_, err := deleteAction.Execute(arguments)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("deleteAction.Execute(dnsDeleter, %q) should return an error because the DNS deleter returned one.", arguments)
	}
}

// deleteAction.Execute should return a success message if the DNS deleter succeeds.
func Test_deleteAction_ValidArguments_DNSCreatorSucceeds_SuccessMessageIsReturned(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-type",
		"AAAA",
	}

	dnsDeleter := &testDNSDeleter{
		deleteSubdomainFunc: func(domain, subdomain string, recordType string) error {
			return nil
		},
	}

	deleteAction := deleteAction{dnsDeleter}

	// act
	response, _ := deleteAction.Execute(arguments)

	// assert
	if response == nil {
		t.Fail()
		t.Logf("deleteAction.Execute(dnsDeleter, %q) should respond with a success message if the DNS deleter succeeds.", arguments)
	}
}

// deleteAction.Execute should return a success message that contains the subdomain, domain and record type.
func Test_deleteAction_ValidArguments_DNSCreatorSucceeds_SuccessMessageContainsTheSubdomainAndType(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
		"-type",
		"AAAA",
	}

	dnsDeleter := &testDNSDeleter{
		deleteSubdomainFunc: func(domain, subdomain string, recordType string) error {
			return nil
		},
	}

	deleteAction := deleteAction{dnsDeleter}

	// act
	response, _ := deleteAction.Execute(arguments)

	// assert
	containsSubdomain := strings.Contains(response.Text(), "www")
	containsDomain := strings.Contains(response.Text(), "example.com")
	containsRecordType := strings.Contains(response.Text(), "AAAA")

	if !containsSubdomain || !containsDomain || !containsRecordType {
		t.Fail()
		t.Logf("deleteAction.Execute(dnsDeleter, %q) should respond with a success message that contains the domain, subdomain and record type but responded with %q instead.", arguments, response.Text())
	}
}
