// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/afero"
	"io/ioutil"
	"os"
)

// newAPICredentials creates a new credentials model from the given
// e-mail address and API token. If the given parameters are invalid
// an error will be returned.
func newAPICredentials(email, token string) (apiCredentials, error) {
	if isEmpty(email) {
		return apiCredentials{}, fmt.Errorf("No e-mail address given")
	}

	if isEmpty(token) {
		return apiCredentials{}, fmt.Errorf("No API token given")
	}

	return apiCredentials{email, token}, nil
}

// apiCredentials contains the credentials for accessing the DNSimple API.
type apiCredentials struct {
	// Email is the E-Mail address that is used for accessing the DNSimple API
	Email string

	// Token is the API token used for accessing the DNSimple API
	Token string
}

// credentialProvider returns credentials.
type credentialProvider interface {
	// GetCredentials returns any stored credentials if there are any.
	// Otherwise GetCredentials will return an error.
	GetCredentials() (apiCredentials, error)
}

// credentialSaver persists credentials.
type credentialSaver interface {
	// SaveCredentials saves the given credentials.
	SaveCredentials(credentials apiCredentials) error
}

// credentialDeleter deletes credentials.
type credentialDeleter interface {
	// DeleteCredentials deletes any saved credentials.
	DeleteCredentials() error
}

// credentialStore provides functions for reading and
// persisting apiCredentials.
type credentialStore interface {
	credentialProvider
	credentialSaver
	credentialDeleter
}

// newFilesystemCredentialStore creates a new filesystem credential store instance.
func newFilesystemCredentialStore(filesystem afero.Fs, filePath string) filesystemCredentialStore {
	return filesystemCredentialStore{
		fs:       filesystem,
		filePath: filePath,
	}
}

// filesystemCredentialStore reads and persists apiCredentials from and to disc.
type filesystemCredentialStore struct {
	fs       afero.Fs
	filePath string
}

// SaveCredentials saves the given credentials to disc.
func (c filesystemCredentialStore) SaveCredentials(credentials apiCredentials) error {

	// check if the file system is initialized
	if c.fs == nil {
		return fmt.Errorf("No filesystem provided")
	}

	// check if the file path is set
	if c.filePath == "" {
		return fmt.Errorf("No file path specified")
	}

	// open the source file for writing
	file, openError := c.fs.OpenFile(c.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if openError != nil {
		return openError
	}

	defer file.Close()

	// convert credentials to JSON
	json, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	// write JSON to file
	fmt.Fprintf(file, "%s", json)

	return nil
}

// DeleteCredentials removes the saved credentials from disc.
func (c filesystemCredentialStore) DeleteCredentials() error {

	// check if the file system is initialized
	if c.fs == nil {
		return fmt.Errorf("No filesystem provided")
	}

	err := c.fs.Remove(c.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return noCredentialsError{fmt.Sprintf("There are no credentials stored at %q", c.filePath)}
		}

		return fmt.Errorf("Deleting %q failed: %s", c.filePath, err.Error())
	}

	return nil
}

// GetCredentials returns any stored credentials if there are any.
// Otherwise GetCredentials will return an error.
func (c filesystemCredentialStore) GetCredentials() (apiCredentials, error) {

	// check if the file system is initialized
	if c.fs == nil {
		return apiCredentials{}, fmt.Errorf("No filesystem specified")
	}

	// check if the file path is set
	if c.filePath == "" {
		return apiCredentials{}, fmt.Errorf("No file path specified")
	}

	// open the source file for reading
	file, openError := c.fs.Open(c.filePath)
	if openError != nil {
		return apiCredentials{}, openError
	}

	defer file.Close()

	// read the source file
	reader := bufio.NewReader(file)
	content, readError := ioutil.ReadAll(reader)
	if readError != nil {
		return apiCredentials{}, readError
	}

	// check if there is content in the file
	if len(content) == 0 {
		return apiCredentials{}, fmt.Errorf("The source file is empty")
	}

	// convert JSON to credentials mopdel
	var credentials apiCredentials
	if unmarshalErr := json.Unmarshal(content, &credentials); unmarshalErr != nil {
		return apiCredentials{}, unmarshalErr
	}

	return credentials, nil
}

type noCredentialsError struct {
	message string
}

func (err noCredentialsError) Error() string {
	return err.message
}

func isNoCredentialsError(err error) bool {
	switch err.(type) {
	case noCredentialsError:
		return true
	}

	return false
}
