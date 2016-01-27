// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"testing"
)

// testDNSClientFactory creates test DNS clients.
type testDNSClientFactory struct {
	client dnsClient
}

// CreateClient create a new DNSimple client instance.
func (clientFactory testDNSClientFactory) CreateClient() (dnsClient, error) {
	if clientFactory.client == nil {
		return nil, fmt.Errorf("No client available")
	}

	return clientFactory.client, nil
}

// If the crendential store returns proper credentials the DNS client factory should return a DNS client.
func Test_dnsimpleClientFactory_CreateClient_CredentialStoreReturnsCredentials_ClientIsReturned(t *testing.T) {
	// arrange
	credentials := apiCredentials{"johndoe@example.com", "abcdefg123"}
	credentialStore := testCredentialsStore{
		getFunc: func() (apiCredentials, error) {
			return credentials, nil
		},
	}

	clientFactory := dnsimpleClientFactory{credentialStore}

	// act
	client, _ := clientFactory.CreateClient()

	// assert
	if client == nil {
		t.Fail()
		t.Logf("clientFactory.CreateClient() should not return nil if proper credentials (%q) were given.", credentials)
	}
}

// The DNS client factory should return an error if the credential store returns an error instead of credentials.
func Test_dnsimpleClientFactory_CreateClient_CredentialStoreReturnsError_ErrorIsReturned(t *testing.T) {
	// arrange
	credentialStore := testCredentialsStore{
		getFunc: func() (apiCredentials, error) {
			return apiCredentials{}, fmt.Errorf("Unable to get credentials")
		},
	}

	clientFactory := dnsimpleClientFactory{credentialStore}

	// act
	_, err := clientFactory.CreateClient()

	// assert
	if err == nil {
		t.Fail()
		t.Logf("clientFactory.CreateClient() should return an error if the credential store does not return credentials.")
	}
}
