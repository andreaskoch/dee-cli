// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deens

import (
	"fmt"
)

// NewAPICredentials creates a new credentials model from the given
// e-mail address and API token. If the given parameters are invalid
// an error will be returned.
func NewAPICredentials(email, token string) (APICredentials, error) {
	if isEmpty(email) {
		return APICredentials{}, fmt.Errorf("No e-mail address given")
	}

	if isEmpty(token) {
		return APICredentials{}, fmt.Errorf("No API token given")
	}

	return APICredentials{email, token}, nil
}

// APICredentials contains the credentials for accessing the DNSimple API.
type APICredentials struct {
	// Email is the E-Mail address that is used for accessing the DNSimple API
	Email string

	// Token is the API token used for accessing the DNSimple API
	Token string
}

// CredentialProvider returns credentials.
type CredentialProvider interface {
	// GetCredentials returns any stored credentials if there are any.
	// Otherwise GetCredentials will return an error.
	GetCredentials() (APICredentials, error)
}

// CredentialSaver persists credentials.
type CredentialSaver interface {
	// SaveCredentials saves the given credentials.
	SaveCredentials(credentials APICredentials) error
}

// CredentialDeleter deletes credentials.
type CredentialDeleter interface {
	// DeleteCredentials deletes any saved credentials.
	DeleteCredentials() error
}

// CredentialStore provides functions for reading and
// persisting APICredentials.
type CredentialStore interface {
	CredentialProvider
	CredentialSaver
	CredentialDeleter
}
