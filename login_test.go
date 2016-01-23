// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"testing"
)

func Test_loginAction_Login_InvalidCredentials_ErrorIsReturned(t *testing.T) {
	// arrange
	var inputs = []struct {
		email string
		token string
	}{
		{"example@example.com", ""},
		{"", "12456"},
		{"", ""},
		{" ", " "},
	}

	loginAction := loginAction{}

	// act
	for _, input := range inputs {

		err := loginAction.Login(input.email, input.token)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("loginAction.Login(%q, %q) should return an error because the input is invalid.", input.email, input.token)
		}
	}
}

func Test_loginAction_Login_ValidCredentials_CredentialsArePassedToCredentialStore(t *testing.T) {
	// arrange
	var inputs = []struct {
		email string
		token string
	}{
		{"example@example.com", "1234"},
		{"example@example", "a"},
		{"test+test@example.co.uk", "ölö23p4k23lö4köl23k4öä"},
	}

	for _, input := range inputs {

		credStore := testCredentialsStore{
			saveFunc: func(credentials apiCredentials) error {

				// assert
				if credentials.Email != input.email || credentials.Token != input.token {
					t.Fail()
					t.Logf("Login(%q, %q) passed invalid credentials to the Save function of the credential store: %s", input.email, input.token, credentials)
				}

				return nil
			},
		}
		loginAction := loginAction{credStore}

		// act
		loginAction.Login(input.email, input.token)
	}
}

func Test_loginAction_Login_ValidCredentials_CredentialStoreSaveFails_ErrorIsReturned(t *testing.T) {
	// arrange
	credStore := testCredentialsStore{
		saveFunc: func(credentials apiCredentials) error {
			return fmt.Errorf("Save failed")
		},
	}
	loginAction := loginAction{credStore}

	// act
	err := loginAction.Login("example@example.com", "1234")

	// assert
	if err == nil {
		t.Fail()
		t.Logf("If the save at the credential store fails Login should return an error.")
	}
}

type testCredentialsStore struct {
	saveFunc   func(credentials apiCredentials) error
	getFunc    func() (apiCredentials, error)
	deleteFunc func() error
}

func (credStore testCredentialsStore) SaveCredentials(credentials apiCredentials) error {
	return credStore.saveFunc(credentials)
}

func (credStore testCredentialsStore) GetCredentials() (apiCredentials, error) {
	return credStore.getFunc()
}

func (credStore testCredentialsStore) DeleteCredentials() error {
	return credStore.deleteFunc()
}
