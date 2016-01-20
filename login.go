// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// newLoginAction create a new instance of the loginAction
// with the given credential store.
func newLoginAction(credentialStore credentialStore) loginAction {
	return loginAction{credentialStore}
}

type loginAction struct {
	credentialStore credentialStore
}

func (login *loginAction) Login(email, token string) error {
	credentials, credentialError := newAPICredentials(email, token)
	if credentialError != nil {
		return credentialError
	}

	return login.credentialStore.SaveCredentials(credentials)
}
