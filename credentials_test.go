// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/spf13/afero"
	"os"
	"path"
	"testing"
)

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

// getFileWithContent returns an afero.File with the given name and content.
func getFileWithContent(name, content string) (afero.File, error) {
	fs := afero.NewMemMapFs()
	fs.MkdirAll(path.Base(name), 0755)
	f, createFileError := fs.OpenFile(name, os.O_CREATE|os.O_RDWR, 0744)
	if createFileError != nil {
		return nil, createFileError
	}

	_, writeError := f.WriteString(content)
	if writeError != nil {
		return nil, writeError
	}
	f.Close()

	f2, openFileError := fs.OpenFile(name, os.O_CREATE|os.O_RDWR, 0744)
	if openFileError != nil {
		return nil, openFileError
	}

	return f2, nil
}

func Test_newAPICredentials_ValidEmailAndToken_NoErrorIsReturned(t *testing.T) {
	// arrange
	var inputs = []struct {
		email string
		token string
	}{
		{"example@example.com", "1234"},
		{"example@example", "a"},
		{"test+test@example.co.uk", "ölö23p4k23lö4köl23k4öä"},
	}

	// act
	for _, input := range inputs {
		_, err := newAPICredentials(input.email, input.token)

		// assert
		if err != nil {
			t.Fail()
			t.Logf("newAPICredentials(%q, %q) should not return an error because the input is valid. But it returned: %s", input.email, input.token, err.Error())
		}
	}
}

func Test_newAPICredentials_InvalidValidEmailOrToken_ErrorIsReturned(t *testing.T) {
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

	// act
	for _, input := range inputs {
		_, err := newAPICredentials(input.email, input.token)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("newAPICredentials(%q, %q) should return an error because the given input is invalid.", input.email, input.token)
		}
	}
}

func Test_filesystemCredentialStore_WithoutFilesystem_SaveCredentials_ErrorIsReturned(t *testing.T) {

	// arrange
	credentialFilePath := "/home/user/.dnsimple-cli/credentials.json"
	credentialStore := filesystemCredentialStore{
		filePath: credentialFilePath,
	}

	// act
	err := credentialStore.SaveCredentials(apiCredentials{"john@example.com", "123456"})

	// assert
	if err == nil {
		t.Fail()
		t.Logf("SaveCredentials should return an error if the credential store has no file system.")
	}
}

func Test_filesystemCredentialStore_WithoutFilePath_SaveCredentials_ErrorIsReturned(t *testing.T) {

	// arrange
	filesystem := afero.NewMemMapFs()
	credentialStore := filesystemCredentialStore{
		fs: filesystem,
	}

	// act
	err := credentialStore.SaveCredentials(apiCredentials{"john@example.com", "123456"})

	// assert
	if err == nil {
		t.Fail()
		t.Logf("SaveCredentials should return an error if the credential store file path configured.")
	}
}

func Test_filesystemCredentialStore_SaveCredentials_CredentialsAreValid_NoErrorIsReturned(t *testing.T) {

	// arrange
	fs := afero.NewMemMapFs()
	credentialFilePath := "/home/user/.dnsimple-cli/credentials.json"
	credentialStore := filesystemCredentialStore{fs, credentialFilePath}

	// act
	err := credentialStore.SaveCredentials(apiCredentials{"john@example.com", "123456"})

	// assert
	if err != nil {
		t.Fail()
		t.Logf("SaveCredentials should not return an error (%q)", err.Error())
	}
}

func Test_filesystemCredentialStore_SaveCredentials_CredentialsAreValid_JSONIsWrittenToFile(t *testing.T) {

	// arrange
	fs := afero.NewMemMapFs()
	credentialFilePath := "/home/user/.dnsimple-cli/credentials.json"
	credentialStore := filesystemCredentialStore{fs, credentialFilePath}

	// act
	credentialStore.SaveCredentials(apiCredentials{"john@example.com", "123456"})

	// assert
	expectedResult := `{"Email":"john@example.com","Token":"123456"}`
	content, _ := afero.ReadFile(fs, credentialFilePath)
	if string(content) != expectedResult {
		t.Fail()
		t.Logf("SaveCredentials should have written the credentials as JSON to %q. Expected: %q, Actual: %q", credentialFilePath, expectedResult, string(content))
	}
}

func Test_filesystemCredentialStore_SaveCredentials_FileExists_FileIsOverridden(t *testing.T) {

	// arrange
	fs := afero.NewMemMapFs()
	credentialFilePath := "/home/user/.dnsimple-cli/credentials.json"

	// create the initial file
	f1, _ := fs.Create(credentialFilePath)
	f1.WriteString(`{"Email":"previous@example.com","Token":"543"}`)
	f1.Close()

	credentialStore := filesystemCredentialStore{fs, credentialFilePath}

	// act
	credentialStore.SaveCredentials(apiCredentials{"new@example.com", "123456"})

	// assert
	expectedResult := `{"Email":"new@example.com","Token":"123456"}`
	content, _ := afero.ReadFile(fs, credentialFilePath)
	if string(content) != expectedResult {
		t.Fail()
		t.Logf("SaveCredentials should have written the credentials as JSON to %q. Expected: %q, Actual: %q", credentialFilePath, expectedResult, string(content))
	}
}

func Test_filesystemCredentialStore_GetCredentials_SavedCredentialsAreValid_CredentialsAreReturned(t *testing.T) {

	// arrange
	fs := afero.NewMemMapFs()
	credentialFilePath := "/home/user/.dnsimple-cli/credentials.json"

	// create the initial file
	f1, _ := fs.Create(credentialFilePath)
	f1.WriteString(`{"Email":"john@example.com","Token":"123456"}`)
	f1.Close()

	credentialStore := filesystemCredentialStore{fs, credentialFilePath}

	// act
	credentials, err := credentialStore.GetCredentials()

	// assert
	if err != nil {
		t.Fail()
		t.Logf("GetCredentials returned an error: %s", err.Error())
	}

	if credentials.Email != "john@example.com" || credentials.Token != "123456" {
		t.Fail()
		t.Logf("GetCredentials did not return the correct credentials. Email: %q, Token: %q", credentials.Email, credentials.Token)
	}
}

func Test_filesystemCredentialStore_GetCredentials_SavedCredentialsAreEmpty_ErrorIsReturned(t *testing.T) {

	// arrange
	fs := afero.NewMemMapFs()
	credentialFilePath := "/home/user/.dnsimple-cli/credentials.json"

	// create the initial file
	f1, _ := fs.Create(credentialFilePath)
	f1.WriteString(``)
	f1.Close()

	credentialStore := filesystemCredentialStore{fs, credentialFilePath}

	// act
	_, err := credentialStore.GetCredentials()

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetCredentials should return an error if the given file is empty.")
	}
}

func Test_filesystemCredentialStore_GetCredentials_JSONIsInvalid_ErrorIsReturned(t *testing.T) {

	// arrange
	fs := afero.NewMemMapFs()
	credentialFilePath := "/home/user/.dnsimple-cli/credentials.json"

	// create the initial file
	f1, _ := fs.Create(credentialFilePath)
	f1.WriteString("dsakldjasl ---")
	f1.Close()

	credentialStore := filesystemCredentialStore{fs, credentialFilePath}

	// act
	_, err := credentialStore.GetCredentials()

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetCredentials should return an error if the given file content is invalid.")
	}
}

func Test_filesystemCredentialStore_WithoutSourceFile_GetCredentials_ErrorIsReturned(t *testing.T) {

	// arrange
	credentialStore := filesystemCredentialStore{}

	// act
	_, err := credentialStore.GetCredentials()

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetCredentials should return an error if the credential store has no source file.")
	}
}

func Test_filesystemCredentialStore_WithoutSourceFile_DeleteCredentials_credentialsNotFoundErrorIsReturned(t *testing.T) {

	// arrange
	fs := afero.NewMemMapFs()
	credentialFilePath := "/home/user/.dnsimple-cli/credentials.json"
	credentialStore := filesystemCredentialStore{fs, credentialFilePath}

	// act
	err := credentialStore.DeleteCredentials()

	// assert
	if !isNoCredentialsError(err) {
		t.Fail()
		t.Logf("DeleteCredentials should return a noCredentials-error if the credential store has no source file.")
	}
}

func Test_filesystemCredentialStore_SourceFileExists_DeleteCredentials_NoErrorIsReturned(t *testing.T) {

	// arrange
	fs := afero.NewMemMapFs()
	credentialFilePath := "/home/user/.dnsimple-cli/credentials.json"

	// create the initial file
	f1, _ := fs.Create(credentialFilePath)
	f1.WriteString("Some content")

	credentialStore := filesystemCredentialStore{fs, credentialFilePath}

	// act
	err := credentialStore.DeleteCredentials()

	// assert
	if err != nil {
		t.Fail()
		t.Logf("DeleteCredentials should not return an error if the credential file exists and was successfully deleted.")
	}
}

func Test_filesystemCredentialStore_SourceFileExists_DeleteCredentials_FileIsDeleted(t *testing.T) {

	// arrange
	fs := afero.NewMemMapFs()
	credentialFilePath := "/home/user/.dnsimple-cli/credentials.json"

	// create the initial file
	f1, _ := fs.Create(credentialFilePath)
	f1.WriteString("Some content")

	credentialStore := filesystemCredentialStore{fs, credentialFilePath}

	// act
	credentialStore.DeleteCredentials()

	// assert
	fileInfo, _ := fs.Stat(credentialFilePath)
	fileExists := fileInfo != nil
	if fileExists {
		t.Fail()
		t.Logf("The file %q should be deleted after DeleteCredentials is executed.", credentialFilePath)
	}
}
