// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/pearkes/dnsimple"
	"strings"
	"testing"
)

// testDNSInfoProvider is a DNS info-provider used for testing.
type testDNSInfoProvider struct {
	getDomainNamesFunc      func() ([]string, error)
	getDomainRecordsFunc    func(domain string) ([]dnsimple.Record, error)
	getSubdomainRecordFunc  func(domain, subdomain, recordType string) (dnsimple.Record, error)
	getSubdomainRecordsFunc func(domain, subdomain string) ([]dnsimple.Record, error)
}

func (infoProvider testDNSInfoProvider) GetDomainNames() ([]string, error) {
	return infoProvider.getDomainNamesFunc()
}

func (infoProvider testDNSInfoProvider) GetDomainRecords(domain string) ([]dnsimple.Record, error) {
	return infoProvider.getDomainRecordsFunc(domain)
}

func (infoProvider testDNSInfoProvider) GetSubdomainRecord(domain, subdomain, recordType string) (record dnsimple.Record, err error) {
	return infoProvider.getSubdomainRecordFunc(domain, subdomain, recordType)
}

func (infoProvider testDNSInfoProvider) GetSubdomainRecords(domain, subdomain string) ([]dnsimple.Record, error) {
	return infoProvider.getSubdomainRecordsFunc(domain, subdomain)
}

// The Name function should return "list"
func Test_listAction_Name_ResultIsLogin(t *testing.T) {
	// arrange
	listAction := listAction{}

	// act
	result := listAction.Name()

	// assert
	if result != "list" {
		t.Fail()
		t.Logf("listAction.Name() should return %q but returned %q instead", "list", result)
	}
}

// The Description function return something
func Test_listAction_Description_ResultIsNotEmpty(t *testing.T) {
	// arrange
	listAction := listAction{}

	// act
	result := listAction.Description()

	// assert
	if isEmpty(result) {
		t.Fail()
		t.Logf("listAction.Description() should return something")
	}
}

// The Usage function return something
func Test_listAction_Usage_ResultIsNotEmpty(t *testing.T) {
	// arrange
	listAction := listAction{}

	// act
	result := listAction.Usage()

	// assert
	if isEmpty(result) {
		t.Fail()
		t.Logf("listAction.Usage() should return something")
	}
}

// If one of the given arguments is invalid an error should be returned
func Test_listAction_InvalidArguments_ErrorIsReturned(t *testing.T) {
	// arrange
	argumentsSet := [][]string{
		{"-domain-NAME", "example.com", "-subdomain", "www"},
		{"-domain", "example.com", "-sub-domain", "www"},
		{"-DOMAIN", "example.com", "-subdomain", "www"},
	}

	for _, arguments := range argumentsSet {
		dnsInfoProvider := testDNSInfoProvider{
			getDomainNamesFunc: func() ([]string, error) {
				return []string{}, nil
			},
			getSubdomainRecordFunc: func(domain, subdomain, recordType string) (dnsimple.Record, error) {
				return dnsimple.Record{}, nil
			},
		}

		infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

		list := listAction{infoProviderFactory}

		// act
		_, err := list.Execute(arguments)

		// assert
		if err == nil || strings.Contains(err.Error(), "flag provided but not") == false {
			t.Fail()
			t.Logf("list.Execute(%q) should return an error.", arguments)
		}
	}
}

// If no info provider factory is present an error is returned.
func Test_listAction_ValidArguments_NoInfoProviderFactory_ErrorIsReturned(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
	}

	list := listAction{}

	// act
	_, err := list.Execute(arguments)

	// assert
	if err == nil || strings.Contains(err.Error(), "No DNS info provider") == false {
		t.Fail()
		t.Logf("list.Execute(%q) should return an error.", arguments)
	}
}

// Case 2, happy path.
func Test_listAction_DomainAndSubDomainAreSet_SubdomainRecordsArePrinted(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
	}

	dnsInfoProvider := testDNSInfoProvider{
		getSubdomainRecordsFunc: func(domain, subdomain string) ([]dnsimple.Record, error) {
			records := []dnsimple.Record{
				dnsimple.Record{
					Name:       "www",
					Content:    "2001:0db8:0000:0042:0000:8a2e:0370:7334",
					RecordType: "AAAA",
				},
				dnsimple.Record{
					Name:       "www",
					Content:    "10.0.2.1",
					RecordType: "A",
				},
			}

			return records, nil
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	list := listAction{infoProviderFactory}

	// act
	result, _ := list.Execute(arguments)

	// assert
	if isEmpty(result.Text()) {
		t.Fail()
		t.Logf("list.Execute(%q) should not return an empty result.", arguments)
	}
}

// Case 2, error path.
func Test_listAction_DomainAndSubDomainAreSet_InfoProviderReturnsError_ErrorIsReturned(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
		"-subdomain",
		"www",
	}

	dnsInfoProvider := testDNSInfoProvider{
		getSubdomainRecordsFunc: func(domain, subdomain string) ([]dnsimple.Record, error) {
			return nil, fmt.Errorf("No records found")
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	list := listAction{infoProviderFactory}

	// act
	_, err := list.Execute(arguments)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("list.Execute(%q) return an error.", arguments)
	}
}

// Case 3, happy path.
func Test_listAction_DomainSet_DNSRecordsArePrinted(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
	}

	dnsInfoProvider := testDNSInfoProvider{
		getDomainRecordsFunc: func(domain string) ([]dnsimple.Record, error) {
			records := []dnsimple.Record{
				dnsimple.Record{
					Name:       "www",
					Content:    "2001:0db8:0000:0042:0000:8a2e:0370:7334",
					RecordType: "AAAA",
				},
				dnsimple.Record{
					Name:       "www",
					Content:    "10.0.2.1",
					RecordType: "A",
				},
			}

			return records, nil
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	list := listAction{infoProviderFactory}

	// act
	result, _ := list.Execute(arguments)

	// assert
	if isEmpty(result.Text()) {
		t.Fail()
		t.Logf("list.Execute(%q) should not return an empty result.", arguments)
	}
}

// Case 3, error path.
func Test_listAction_DomainSet_InfoProviderReturnsError_ErrorIsReturned(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
	}

	dnsInfoProvider := testDNSInfoProvider{
		getDomainRecordsFunc: func(domain string) ([]dnsimple.Record, error) {
			return nil, fmt.Errorf("Error DNS record")
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	list := listAction{infoProviderFactory}

	// act
	_, err := list.Execute(arguments)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("list.Execute(%q) return an error.", arguments)
	}
}

// Case 1, happy path.
func Test_listAction_NoArguments_DomainsArePrinted(t *testing.T) {
	// arrange
	arguments := []string{}

	dnsInfoProvider := testDNSInfoProvider{
		getDomainNamesFunc: func() ([]string, error) {
			return []string{
				"example.com",
				"example.co.uk",
			}, nil
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	list := listAction{infoProviderFactory}

	// act
	result, _ := list.Execute(arguments)

	// assert
	if isEmpty(result.Text()) {
		t.Fail()
		t.Logf("list.Execute(%q) should not return an empty result.", arguments)
	}
}

// Case 1, happy path.
func Test_listAction_NoArguments_InfoProviderReturnsError_ErrorIsReturned(t *testing.T) {
	// arrange
	arguments := []string{}

	dnsInfoProvider := testDNSInfoProvider{
		getDomainNamesFunc: func() ([]string, error) {
			return []string{}, fmt.Errorf("Domain error")
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	list := listAction{infoProviderFactory}

	// act
	_, err := list.Execute(arguments)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("list.Execute(%q) return an error.", arguments)
	}
}

// The first column should contain the subdomain and domain.
func Test_formatDNSRecords_SubdomainIsSet_ResultContainsSubdomain(t *testing.T) {
	// arrange
	records := []dnsimple.Record{
		dnsimple.Record{
			Name:       "www",
			Content:    "2001:0db8:0000:0042:0000:8a2e:0370:7334",
			RecordType: "AAAA",
		},
		dnsimple.Record{
			Name:       "www",
			Content:    "10.0.2.1",
			RecordType: "A",
		},
	}
	domain := "example.com"

	// act
	result := formatDNSRecords(records, domain)

	// assert
	expectedResult := `www.example.com   AAAA   2001:0db8:0000:0042:0000:8a2e:0370:7334
www.example.com   A      10.0.2.1`

	if result != expectedResult {
		t.Fail()
		t.Logf("formatDNSRecords(%q, %q)\n", records, domain)

		t.Logf("Should have returned:\n")
		t.Logf("%s\n", expectedResult)
		t.Log("\n")

		t.Logf("But returned this instead:\n")
		t.Logf("%s\n", result)
	}
}

// Case 3, happy path.
func Test_listAction_OnlyDomainIsSet_DomainRecordsArePrinted(t *testing.T) {
	// arrange
	arguments := []string{
		"-domain",
		"example.com",
	}

	dnsInfoProvider := testDNSInfoProvider{
		getDomainRecordsFunc: func(domain string) ([]dnsimple.Record, error) {
			records := []dnsimple.Record{
				dnsimple.Record{
					Name:       "www",
					Content:    "2001:0db8:0000:0042:0000:8a2e:0370:7334",
					RecordType: "AAAA",
				},
				dnsimple.Record{
					Name:       "www",
					Content:    "10.0.2.1",
					RecordType: "A",
				},
			}

			return records, nil
		},
	}

	infoProviderFactory := testInfoProviderFactory{dnsInfoProvider, nil}

	list := listAction{infoProviderFactory}

	// act
	result, _ := list.Execute(arguments)

	// assert
	if isEmpty(result.Text()) {
		t.Fail()
		t.Logf("list.Execute(%q) should not return an empty result.", arguments)
	}
}

// The first column should contain the domain name without the subdomain.
func Test_formatDNSRecords_SubdomainIsNotSet_ResultDoesNotContainSubdomain(t *testing.T) {
	// arrange
	records := []dnsimple.Record{
		dnsimple.Record{
			Name:       "",
			Content:    "2001:0db8:0000:0042:0000:8a2e:0370:7334",
			RecordType: "AAAA",
		},
		dnsimple.Record{
			Name:       "",
			Content:    "10.0.2.1",
			RecordType: "A",
		},
	}
	domain := "example.com"

	// act
	result := formatDNSRecords(records, domain)

	// assert
	expectedResult := `example.com   AAAA   2001:0db8:0000:0042:0000:8a2e:0370:7334
example.com   A      10.0.2.1`

	if result != expectedResult {
		t.Fail()
		t.Logf("formatDNSRecords(%q, %q)\n", records, domain)

		t.Logf("Should have returned:\n")
		t.Logf("%s\n", expectedResult)
		t.Log("\n")

		t.Logf("But returned this instead:\n")
		t.Logf("%s\n", result)
	}
}

// The result should be empty if the given DNS record list is empty.
func Test_formatDNSRecords_EmptyRecordList_ResultIsEmpty(t *testing.T) {
	// arrange
	records := []dnsimple.Record{}
	domain := "example.com"

	// act
	result := formatDNSRecords(records, domain)

	// assert
	expectedResult := ``
	if result != expectedResult {
		t.Fail()
		t.Logf("formatDNSRecords(%q, %q)\n", records, domain)

		t.Logf("Should have returned:\n")
		t.Logf("%s\n", expectedResult)
		t.Log("\n")

		t.Logf("But returned this instead:\n")
		t.Logf("%s\n", result)
	}
}

// The result should not end with a newline character.
func Test_formatDNSRecords_ResultDoesNotEndWithNewline(t *testing.T) {
	// arrange
	records := []dnsimple.Record{
		dnsimple.Record{
			Name:       "www",
			Content:    "2001:0db8:0000:0042:0000:8a2e:0370:7334",
			RecordType: "AAAA",
		},
		dnsimple.Record{
			Name:       "www",
			Content:    "10.0.2.1",
			RecordType: "A",
		},
	}
	domain := "example.com"

	// act
	result := formatDNSRecords(records, domain)

	// assert
	if strings.HasSuffix(result, "\n") {
		t.Fail()
		t.Logf("formatDNSRecords(%q, %q) should not end with a newline character", records, domain)
	}
}
