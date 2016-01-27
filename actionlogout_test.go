// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
	"testing"
)

func Test_logoutAction_Name_LogoutIsReturned(t *testing.T) {

	// arrange
	credentialStore := testCredentialsStore{deleteFunc: func() error {
		return nil
	}}

	logoutAction := logoutAction{credentialStore}

	// act
	result := logoutAction.Name()

	// assert
	if result != "logout" {
		t.Fail()
		t.Logf("logoutAction.Name() should have returned %q but returned %q instead.", "logout", result)
	}

}

func Test_logoutAction_Description_ResultIsNotEmpty(t *testing.T) {

	// arrange
	credentialStore := testCredentialsStore{deleteFunc: func() error {
		return nil
	}}

	logoutAction := logoutAction{credentialStore}

	// act
	result := logoutAction.Description()

	// assert
	if isEmpty(result) {
		t.Fail()
		t.Logf("logoutAction.Description() not be empty.")
	}

}

func Test_logoutAction_Usage_ResultIsNotEmpty(t *testing.T) {

	// arrange
	credentialStore := testCredentialsStore{deleteFunc: func() error {
		return nil
	}}

	logoutAction := logoutAction{credentialStore}

	// act
	result := logoutAction.Usage()

	// assert
	if isEmpty(result) {
		t.Fail()
		t.Logf("logoutAction.Usage() not be empty.")
	}

}

func Test_logoutAction_NoCredentialStore_ErrorIsReturned(t *testing.T) {

	// arrange
	logoutAction := logoutAction{}

	// act
	_, err := logoutAction.Execute([]string{})

	// assert
	if err == nil {
		t.Fail()
		t.Logf("logoutAction.Execute() should return an error if no credential store is present.")
	}

}

func Test_logoutAction_CredentialStoreSucceedsInDeletingTheCredentials_NoErrorIsReturned(t *testing.T) {

	// arrange
	credentialStore := testCredentialsStore{deleteFunc: func() error {
		return nil
	}}

	logout := logoutAction{credentialStore}

	// act
	_, err := logout.Execute([]string{})

	// assert
	if err != nil {
		t.Fail()
		t.Logf("logout(credentialStore) should not return an error if the credential store succeeds in deleting the credentials.")
	}

}

func Test_logoutAction_CredentialStoreReturnsNoCredentialsError_NoLogoutRequiredErrorIsReturned(t *testing.T) {

	// arrange
	credentialStore := testCredentialsStore{deleteFunc: func() error {
		return noCredentialsError{"file does not exist"}
	}}

	logout := logoutAction{credentialStore}

	// act
	_, err := logout.Execute([]string{})

	// assert
	if !strings.Contains(err.Error(), "No logout required") {
		t.Fail()
		t.Logf("logout(credentialStore) should return an error stating that no logout was required.")
	}

}

func Test_logoutAction_CredentialStoreReturnsGenericError_LogoutFailedErrorIsReturned(t *testing.T) {

	// arrange
	credentialStore := testCredentialsStore{deleteFunc: func() error {
		return fmt.Errorf("Some error")
	}}

	logout := logoutAction{credentialStore}

	// act
	_, err := logout.Execute([]string{})

	// assert
	if !strings.Contains(err.Error(), "Logout failed") {
		t.Fail()
		t.Logf("logout(credentialStore) should return an error stating that the logout failed.")
	}

}
