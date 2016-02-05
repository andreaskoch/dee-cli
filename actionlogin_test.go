// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/andreaskoch/dee-ns"
	"github.com/spf13/afero"
	"testing"
)

// The Name function should return "login"
func Test_loginAction_Name_ResultIsLogin(t *testing.T) {
	// arrange

	credentialStore := testCredentialsStore{saveFunc: func(credentials deens.APICredentials) error {
		return nil
	}}

	login := loginAction{credentialStore}

	// act
	result := login.Name()

	// assert
	if result != "login" {
		t.Fail()
		t.Logf("login.Name() should return %q but returned %q instead", "login", result)
	}
}

// The Description function return something
func Test_loginAction_Description_ResultIsNotEmpty(t *testing.T) {
	// arrange

	credentialStore := testCredentialsStore{saveFunc: func(credentials deens.APICredentials) error {
		return nil
	}}

	login := loginAction{credentialStore}

	// act
	result := login.Description()

	// assert
	if isEmpty(result) {
		t.Fail()
		t.Logf("login.Description() should return something")
	}
}

// The Usage function return something
func Test_loginAction_Usage_ResultIsNotEmpty(t *testing.T) {
	// arrange

	credentialStore := testCredentialsStore{saveFunc: func(credentials deens.APICredentials) error {
		return nil
	}}

	login := loginAction{credentialStore}

	// act
	result := login.Usage()

	// assert
	if isEmpty(result) {
		t.Fail()
		t.Logf("login.Usage() should return something")
	}
}

// If invalid arguments are supplied an error should be returned.
func Test_loginAction_InvalidArguments_ErrorIsReturned(t *testing.T) {
	// arrange
	validArgumentsSet := [][]string{
		{
			"-e-mail",
			"user@example.com",
			"-apitoken",
			"123456",
		},
		{
			"--eMail",
			"user@example.com",
			"--apitoken",
			"123456",
		},
		{
			"-email=user@example.com",
			"-token=123456",
		},
		{
			"--email=user@example.com",
			"----apitoken=123456",
		},
	}

	for _, arguments := range validArgumentsSet {
		credentialStore := testCredentialsStore{saveFunc: func(credentials deens.APICredentials) error {
			return nil
		}}

		login := loginAction{credentialStore}

		// act
		_, err := login.Execute(arguments)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("login.Execute(%q) should return an error", arguments)
		}
	}
}

// If valid arguments are supplied to the login method no error should be returned.
func Test_loginAction_ValidArguments_NoErrorIsReturned(t *testing.T) {
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
		credentialStore := testCredentialsStore{saveFunc: func(credentials deens.APICredentials) error {
			return nil
		}}

		login := loginAction{credentialStore}

		// act
		_, err := login.Execute(arguments)

		// assert
		if err != nil {
			t.Fail()
			t.Logf("login.Execute(%q) should not return an error: %q", arguments, err.Error())
		}
	}
}

// If the supplied argument values are invalid an error should be returned
func Test_loginAction_ValidArguments_ButInvalidValues_ErrorIsReturned(t *testing.T) {
	// arrange
	validArgumentsSet := [][]string{
		{
			"-email",
			"",
			"-apitoken",
			"123456",
		},
		{
			"--email",
			"user@example.com",
			"--apitoken",
			"",
		},
		{
			"-email=  ",
			"-apitoken=123456",
		},
		{
			"--email=user@example.com",
			"--apitoken=  ",
		},
	}

	for _, arguments := range validArgumentsSet {
		credentialStore := testCredentialsStore{saveFunc: func(credentials deens.APICredentials) error {
			return nil
		}}

		login := loginAction{credentialStore}

		// act
		_, err := login.Execute(arguments)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("login.Execute(%q) should return an error", arguments)
		}
	}
}

// Integration test.
func Test_loginAction_ValidArguments_CredentialFileIsCreated(t *testing.T) {
	// arrange
	filesystem := afero.NewMemMapFs()
	credentialStore := newFilesystemCredentialStore(filesystem, "/home/testuser/.dee/credentials.json")
	arguments := []string{
		"-email",
		"user@example.com",
		"-apitoken",
		"123456",
	}

	login := loginAction{credentialStore}

	// act
	login.Execute(arguments)

	// assert
	fileInfo, err := filesystem.Stat("/home/testuser/.dee/credentials.json")
	if fileInfo == nil || err != nil {
		t.Fail()
		t.Logf("login.Execute(%q) create the credential file and should not return an error: %s", arguments, err.Error())
	}
}

func Test_loginAction_Login_InvalidCredentials_ErrorIsReturned(t *testing.T) {
	// arrange
	var inputs = [][]string{
		{"--email=example@example.com", "--apitoken"},
		{"--email", "--apitoken=12456"},
		{"--email", "--apitoken"},
		{"--email=\" \"", "--apitoken=\" \""},
	}

	loginAction := loginAction{}

	// act
	for _, arguments := range inputs {

		_, err := loginAction.Execute(arguments)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("login.Execute(%q) should return an error because the input is invalid.", arguments)
		}
	}
}

func Test_loginAction_Login_ValidCredentials_CredentialsArePassedToCredentialStore(t *testing.T) {
	// arrange
	var inputs = [][]string{
		{"-email", "example@example.com", "-apitoken", "1234"},
		{"-email", "example@example", "-apitoken", "a"},
		{"-email", "test+test@example.co.uk", "-apitoken", "ölö23p4k23lö4köl23k4öä"},
	}

	for _, arguments := range inputs {

		credStore := testCredentialsStore{
			saveFunc: func(credentials deens.APICredentials) error {

				// assert
				if credentials.Email != arguments[1] || credentials.Token != arguments[3] {
					t.Fail()
					t.Logf("Login(%q, %q) passed invalid credentials to the Save function of the credential store: %s", arguments[1], arguments[3], credentials)
				}

				return nil
			},
		}
		loginAction := loginAction{credStore}

		// act
		loginAction.Execute(arguments)
	}
}

func Test_loginAction_Login_ValidCredentials_CredentialStoreSaveFails_ErrorIsReturned(t *testing.T) {
	// arrange
	credStore := testCredentialsStore{
		saveFunc: func(credentials deens.APICredentials) error {
			return fmt.Errorf("Save failed")
		},
	}
	loginAction := loginAction{credStore}

	// act
	_, err := loginAction.Execute([]string{"-email", "example@example.com", "-apitoken", "1234"})

	// assert
	if err == nil {
		t.Fail()
		t.Logf("If the save at the credential store fails Login should return an error.")
	}
}
