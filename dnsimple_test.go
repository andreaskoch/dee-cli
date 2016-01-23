// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/spf13/afero"
	"net"
	"strings"
	"testing"
)

// If valid arguments are supplied to the login method no error should be returned.
func Test_login_ValidArguments_NoErrorIsReturned(t *testing.T) {
	// arrange
	validArgumentsSet := [][]string{
		{
			"-email",
			"user@example.com",
			"-apitoken",
			"123456",
		},
		{
			"--email",
			"user@example.com",
			"--apitoken",
			"123456",
		},
		{
			"-email=user@example.com",
			"-apitoken=123456",
		},
		{
			"--email=user@example.com",
			"--apitoken=123456",
		},
	}

	for _, arguments := range validArgumentsSet {
		credentialStore := testCredentialsStore{saveFunc: func(credentials apiCredentials) error {
			return nil
		}}

		// act
		err := login(credentialStore, arguments)

		// assert
		if err != nil {
			t.Fail()
			t.Logf("login(credentialStore, %q) should not return an error: %q", arguments, err.Error())
		}
	}
}

// Integration test.
func Test_login_ValidArguments_CredentialFileIsCreated(t *testing.T) {
	// arrange
	filesystem := afero.NewMemMapFs()
	credentialStore := newFilesystemCredentialStore(filesystem, "/home/testuser/.dnsimple-cli/credentials.json")
	arguments := []string{
		"-email",
		"user@example.com",
		"-apitoken",
		"123456",
	}

	// act
	login(credentialStore, arguments)

	// assert
	fileInfo, err := filesystem.Stat("/home/testuser/.dnsimple-cli/credentials.json")
	if fileInfo == nil || err != nil {
		t.Fail()
		t.Logf("login(credentialStore, %q) create the credential file and should not return an error: %s", arguments, err.Error())
	}
}

// update should not return an error if the arguments are valid and the subdomain update succeeds.
func Test_update_ValidArguments_NoErrorIsReturned(t *testing.T) {
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

	for _, arguments := range validArgumentsSet {

		// act
		_, err := update(dnsUpdater, arguments)

		// assert
		if err != nil {
			t.Fail()
			t.Logf("update(dnsUpdater, %q) should not return an error: %q", arguments, err.Error())
		}
	}
}

// update should return an error if the DNS Updater responds with one.
func Test_update_ValidArguments_DNSUpdaterRespondsWithError_ErrorIsReturned(t *testing.T) {
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

	// act
	_, err := update(dnsUpdater, arguments)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("update(dnsUpdater, %q) should return an error because the DNS updater returned one.", arguments)
	}
}

// update should return a success message if the DNS updater succeeds.
func Test_update_ValidArguments_DNSUpdaterSucceeds_SuccessMessageIsReturned(t *testing.T) {
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

	// act
	response, _ := update(dnsUpdater, arguments)

	// assert
	if response == nil {
		t.Fail()
		t.Logf("update(dnsUpdater, %q) should respond with a success message if the DNS updater succeeds.", arguments)
	}
}

// update should return a success message that contains the subdomain, domain and ip.
func Test_update_ValidArguments_DNSUpdaterSucceeds_SuccessMessageContainsTheSubdomainAndNewIP(t *testing.T) {
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

	// act
	response, _ := update(dnsUpdater, arguments)

	// assert
	containsIP := strings.Contains(response.Text(), "2001:db8:0:42:0:8a2e:370:7334")
	containsSubdomain := strings.Contains(response.Text(), "www")
	containsDomain := strings.Contains(response.Text(), "example.com")

	if !containsIP || !containsSubdomain || !containsDomain {
		t.Fail()
		t.Logf("update(dnsUpdater, %q) should respond with a success message that contains the domain, subdomain and ip but responded with %q instead.", arguments, response.Text())
	}
}

func Test_logout_CredentialStoreSucceedsInDeletingTheCredentials_NoErrorIsReturned(t *testing.T) {

	// arrange
	credentialStore := testCredentialsStore{deleteFunc: func() error {
		return nil
	}}

	// act
	err := logout(credentialStore)

	// assert
	if err != nil {
		t.Fail()
		t.Logf("logout(credentialStore) should not return an error if the credential store succeeds in deleting the credentials.")
	}

}

func Test_logout_CredentialStoreReturnsNoCredentialsError_NoLogoutRequiredErrorIsReturned(t *testing.T) {

	// arrange
	credentialStore := testCredentialsStore{deleteFunc: func() error {
		return noCredentialsError{"file does not exist"}
	}}

	// act
	err := logout(credentialStore)

	// assert
	if !strings.Contains(err.Error(), "No logout required") {
		t.Fail()
		t.Logf("logout(credentialStore) should return an error stating that no logout was required.")
	}

}

func Test_logout_CredentialStoreReturnsGenericError_LogoutFailedErrorIsReturned(t *testing.T) {

	// arrange
	credentialStore := testCredentialsStore{deleteFunc: func() error {
		return fmt.Errorf("Some error")
	}}

	// act
	err := logout(credentialStore)

	// assert
	if !strings.Contains(err.Error(), "Logout failed") {
		t.Fail()
		t.Logf("logout(credentialStore) should return an error stating that the logout failed.")
	}

}

// dnsimpleUpdater updates DNSimple domain records.
type testDNSUpdater struct {
	updateSubdomainFunc func(domain, subdomain string, ip net.IP) error
}

func (updater *testDNSUpdater) UpdateSubdomain(domain, subdomain string, ip net.IP) error {
	return updater.updateSubdomainFunc(domain, subdomain, ip)
}
