// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
)

var (
	actionNameLogin      = "login"
	loginActionArguments = flag.NewFlagSet(actionNameLogin, flag.ContinueOnError)
	emailAddress         = loginActionArguments.String("email", "", "The e-mail address of the account to use")
	apiToken             = loginActionArguments.String("apitoken", "", "The API token")
)

type loginAction struct {
	credentialStore credentialStore
}

func (action loginAction) Name() string {
	return actionNameLogin
}

func (action loginAction) Description() string {
	return "Save DNSimple API credentials to disc"
}

func (action loginAction) Usage() string {
	buf := new(bytes.Buffer)
	loginActionArguments.SetOutput(buf)
	loginActionArguments.PrintDefaults()
	return buf.String()
}

// Execute parses the e-mail address and API token
// from the given arguments and stores the credentials
// in the given credential store. If the credentials are
// invalid or the save failed and error is returned.
func (action loginAction) Execute(arguments []string) (message, error) {

	if action.credentialStore == nil {
		return nil, fmt.Errorf("No credential store present")
	}

	// parse the command line arguments
	if parseError := loginActionArguments.Parse(arguments); parseError != nil {
		return nil, parseError
	}

	// perform the login action
	credentials, credentialError := newAPICredentials(*emailAddress, *apiToken)
	if credentialError != nil {
		return nil, credentialError
	}

	if saveErr := action.credentialStore.SaveCredentials(credentials); saveErr != nil {
		return nil, saveErr
	}

	return successMessage{"Login succeeded"}, nil
}
