// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/andreaskoch/dee-ns"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// GitInfo is either the empty string (the default)
// or is set to the git hash of the most recent commit
// using the -X linker flag (Example: "2015-01-11-284c030+")
var GitInfo string

// version returns the git version of this binary (e.g. "2015-01-11-284c030+").
// If the linker flags were not provided, the return value is "unknown".
func version() string {
	if GitInfo != "" {
		return GitInfo
	}

	return "unknown"
}

var actions []action

type action interface {
	Name() string
	Description() string
	Usage() string
	Execute(arguments []string) (message, error)
}

func init() {

	// file system
	filesystem := afero.NewOsFs()

	// locate the home directory
	userHomeDir, homeDirError := homedir.Dir()
	if homeDirError != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine home directory: %s\n", homeDirError.Error())
		os.Exit(1)
	}

	// base folder
	baseFolder := getSettingsFolder(filesystem, userHomeDir)

	// credential store
	credentialFilePath := filepath.Join(baseFolder, "credentials.json")
	credentialStore := filesystemCredentialStore{filesystem, credentialFilePath}

	// DNS client factory
	dnsClientFactory := dnsimpleClientFactory{credentialStore}

	// create DNSimple info provider
	dnsInfoProviderFactory := dnsimpleInfoProviderFactory{dnsClientFactory}

	// create a DNS editor instance
	dnsEditorFactory := dnsEditorFactory{dnsClientFactory, dnsInfoProviderFactory}

	actions = []action{
		loginAction{credentialStore},
		logoutAction{credentialStore},
		listAction{dnsInfoProviderFactory},
		createAction{dnsEditorFactory, os.Stdin},
		updateAction{dnsEditorFactory, os.Stdin},
		deleteAction{dnsEditorFactory},
		createOrUpdateAction{dnsEditorFactory, dnsInfoProviderFactory, os.Stdin},
	}

	// override the help information printer
	// of the flag package
	executablePath := os.Args[0]
	executableName := path.Base(executablePath)
	usagePrinter := newUsagePrinter(executableName, version(), actions)

	flag.Usage = func() {
		usagePrinter.PrintUsageInformation(os.Stdout)
	}

}

func main() {

	// get action
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	// get the action name
	selectedActionName := strings.TrimSpace(strings.ToLower(os.Args[1]))

	// find a matching action
	selectedAction := getActionByName(selectedActionName, actions)
	if selectedAction == nil {
		fmt.Fprintf(os.Stderr, "Unknown action: %q\n", selectedActionName)
		os.Exit(1)
	}

	// execute the action
	message, err := selectedAction.Execute(os.Args[2:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s\n", message.Text())
	os.Exit(0)

}

// getActionByName returns the action which matches the given name from the list.
func getActionByName(actionName string, actions []action) action {

	for _, action := range actions {
		if action.Name() != actionName {
			continue
		}
		return action
	}

	return nil
}

// getSettingsFolder returns the path of the settings folder
// and ensures that the folder exists.
func getSettingsFolder(fs afero.Fs, baseFolder string) string {
	settingsFolder := filepath.Join(baseFolder, ".dee")
	createFolderError := fs.MkdirAll(settingsFolder, 0700)
	if createFolderError != nil {
		panic(createFolderError)
	}

	return settingsFolder
}

type message interface {
	Text() string
}

// successMessage contains a text-message indicating success.
type successMessage struct {
	text string
}

// Text returns the text of the current message.
func (m successMessage) Text() string {
	return m.text
}

// dnsClientFactory provides the ability to create DNS clients.
type dnsClientFactory interface {
	// CreateClient create a new dnsClient client instance.
	CreateClient() (deens.DNSClient, error)
}

// dnsimpleClientFactory creates DNSimple clients.
type dnsimpleClientFactory struct {
	credentialStore deens.CredentialStore
}

// CreateClient create a new DNSimple client instance.
func (clientFactory dnsimpleClientFactory) CreateClient() (deens.DNSClient, error) {

	// get the credentials
	credentials, credentialError := clientFactory.credentialStore.GetCredentials()
	if credentialError != nil {
		return nil, fmt.Errorf("%s", credentialError.Error())
	}

	// create a DNSimple client
	dnsimpleClient, dnsimpleClientError := deens.NewDNSClient(credentials)
	if dnsimpleClientError != nil {
		return nil, fmt.Errorf("Unable to create DNSimple client. Error: %s", dnsimpleClientError.Error())
	}

	return dnsimpleClient, nil
}

type dnsInfoProviderCreator interface {
	CreateInfoProvider() (deens.DNSInfoProvider, error)
}

type dnsimpleInfoProviderFactory struct {
	clientFactory dnsClientFactory
}

func (infoFactory dnsimpleInfoProviderFactory) CreateInfoProvider() (deens.DNSInfoProvider, error) {
	client, err := infoFactory.clientFactory.CreateClient()
	if err != nil {
		return nil, err
	}

	return deens.NewDNSInfoProvider(client), nil
}

type dnsEditorCreator interface {
	CreateDNSEditor() (deens.DNSRecordEditor, error)
}

type dnsEditorFactory struct {
	clientFactory       dnsClientFactory
	infoProviderFactory dnsInfoProviderCreator
}

func (editorFactory dnsEditorFactory) CreateDNSEditor() (deens.DNSRecordEditor, error) {
	client, err := editorFactory.clientFactory.CreateClient()
	if err != nil {
		return nil, err
	}

	infoProvider, err := editorFactory.infoProviderFactory.CreateInfoProvider()
	if err != nil {
		return nil, err
	}

	return deens.NewDNSEditor(client, infoProvider), nil
}
