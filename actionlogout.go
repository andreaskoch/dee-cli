// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/andreaskoch/dee-ns"
)

var (
	actionNameLogout = "logout"
)

type logoutAction struct {
	credentialStore deens.CredentialStore
}

func (action logoutAction) Name() string {
	return actionNameLogout
}

func (action logoutAction) Description() string {
	return "Remove any stored DNSimple API credentials from disc"
}

func (action logoutAction) Usage() string {
	return "  <no options required>\n"
}

// Execute deletes the API credentials.
func (action logoutAction) Execute(arguments []string) (message, error) {

	if action.credentialStore == nil {
		return nil, fmt.Errorf("No credential store present")
	}

	err := action.credentialStore.DeleteCredentials()
	if err == nil {
		return successMessage{"Logout succeeded"}, nil
	}

	if isNoCredentialsError(err) {
		return nil, fmt.Errorf("No logout required: %s", err.Error())
	}

	return nil, fmt.Errorf("Logout failed: %s", err.Error())
}
