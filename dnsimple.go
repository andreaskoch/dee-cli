// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pearkes/dnsimple"
	"github.com/spf13/afero"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// GitInfo is either the empty string (the default)
// or is set to the git hash of the most recent commit
// using the -X linker flag (Example: "2015-01-11-284c030+")
var GitInfo string

const actionNameUpdate = "update"
const actionNameLogin = "login"

var (

	// action: login
	loginActionArguments = flag.NewFlagSet("login", flag.ContinueOnError)
	emailAddress         = loginActionArguments.String("email", "", "The e-mail address of the account to use")
	apiToken             = loginActionArguments.String("apitoken", "", "The API token")

	// action: update
	updateSubdomainArguments = flag.NewFlagSet("update-subdomain", flag.ContinueOnError)
	domain                   = updateSubdomainArguments.String("domain", "", "Domain (e.g. example.com")
	subdomain                = updateSubdomainArguments.String("subdomain", "", "Subdomain (e.g. wwww)")
	ipAddress                = updateSubdomainArguments.String("ip", "", "IP address (e.g. 127.0.0.1. ::1")
)

func init() {

	executablePath := os.Args[0]
	executableName := path.Base(executablePath)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s updates DNS records via DNSimple.\n", executableName)
		fmt.Fprintf(os.Stderr, "\n")

		fmt.Fprintf(os.Stderr, "Version: %s\n", version())
		fmt.Fprintf(os.Stderr, "\n")

		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s <action> [arguments ...]\n", executableName)
		fmt.Fprintf(os.Stderr, "\n")

		fmt.Fprintf(os.Stderr, "Actions:\n\n")
		fmt.Fprintf(os.Stderr, "%10s  %s\n", actionNameLogin, "Save DNSimple API credentials to disc")
		fmt.Fprintf(os.Stderr, "%10s  %s\n", actionNameUpdate, "Update the DNS record for a given sub domain")
		fmt.Fprintf(os.Stderr, "\n")

		fmt.Fprintf(os.Stderr, "Action: %s\n\n", actionNameLogin)
		loginActionArguments.PrintDefaults()

		fmt.Fprintf(os.Stderr, "\n")

		fmt.Fprintf(os.Stderr, "Action: %s\n\n", actionNameUpdate)
		updateSubdomainArguments.PrintDefaults()
	}

}

func main() {

	// get action
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

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

	action := strings.TrimSpace(strings.ToLower(os.Args[1]))
	switch action {
	case actionNameLogin:
		{
			// select the login arguments
			loginArguments := os.Args[2:]

			// perform the login
			if loginError := login(credentialStore, loginArguments); loginError != nil {
				fmt.Fprintf(os.Stderr, "%s\n", loginError.Error())
				os.Exit(1)
			}

			fmt.Fprintf(os.Stdout, "Your API credentials have been saved to %s\n", credentialFilePath)
			os.Exit(0)
		}

	case actionNameUpdate:
		{
			// get the credentials
			credentials, credentialError := credentialStore.GetCredentials()
			if credentialError != nil {
				fmt.Fprintf(os.Stderr, "%s\n", credentialError.Error())
				os.Exit(1)
			}

			// create a DNSimple client
			dnsimpleClient, dnsimpleClientError := dnsimple.NewClient(credentials.Email, credentials.Token)
			if dnsimpleClientError != nil {
				fmt.Fprintf(os.Stderr, "Unable to create DNSimple client. Error: %s\n", dnsimpleClientError.Error())
				os.Exit(1)
			}

			// create DNSimple info provider
			dnsimpleInfoProvider := newDNSimpleInfoProvider(dnsimpleClient)

			// create DNSimple domain updater
			dnsimpleUpdater := newDNSimpleUpdater(dnsimpleClient, dnsimpleInfoProvider)

			message, updateError := update(dnsimpleUpdater, os.Args[2:])
			if updateError != nil {
				fmt.Fprintf(os.Stderr, "%s\n", updateError.Error())
				os.Exit(1)
			}

			fmt.Fprintf(os.Stdout, "%s\n", message.Text())
			os.Exit(0)
		}

	default:
		{
			fmt.Fprintf(os.Stderr, "Unknown action: %q\n", action)
			os.Exit(1)
		}
	}
}

// login parses the e-mail address and API token
// from the given arguments and stores the credentials
// in the given credential store. If the credentials are
// invalid or the save failed and error is returned.
func login(credentialStore credentialStore, arguments []string) error {

	// parse the command line arguments
	if parseError := loginActionArguments.Parse(arguments); parseError != nil {
		return parseError
	}

	// perform the login action
	loginAction := newLoginAction(credentialStore)
	err := loginAction.Login(*emailAddress, *apiToken)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	return nil
}

// update updates the DNS record of the domain given from the supplied arguments.
// If the update fails an error is returned.
func update(domainUpdater updater, arguments []string) (message, error) {

	// parse the arguments
	if parseError := updateSubdomainArguments.Parse(arguments); parseError != nil {
		return nil, parseError
	}

	// domain
	if *domain == "" {
		return nil, fmt.Errorf("No domain supplied.")
	}

	// subdomain
	if *subdomain == "" {
		return nil, fmt.Errorf("No subdomain supplied.")
	}

	// take ip from stdin
	if *ipAddress == "" && stdinHasData() {
		ipAddressFromStdin := ""
		fmt.Fscanf(os.Stdin, "%s", &ipAddressFromStdin)
		ipAddress = &ipAddressFromStdin
	}

	if *ipAddress == "" {
		return nil, fmt.Errorf("No IP address supplied.")
	}

	ip := net.ParseIP(*ipAddress)
	updateError := domainUpdater.UpdateSubdomain(*domain, *subdomain, ip)
	if updateError != nil {
		return nil, fmt.Errorf("%s", updateError.Error())
	}

	return successMessage{fmt.Sprintf("Updated: %s.%s → %s", *subdomain, *domain, ip.String())}, nil
}

// stdinHasData returns true if there is data avaialble in os.Stdin, otherwise false.
// see: http://stackoverflow.com/questions/22744443/check-if-there-is-something-to-read-on-stdin-in-golang
func stdinHasData() bool {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return true
	}

	return false
}

// getSettingsFolder returns the path of the settings folder
// and ensures that the folder exists.
func getSettingsFolder(fs afero.Fs, baseFolder string) string {
	settingsFolder := filepath.Join(baseFolder, ".dnsimple-cli")
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

// version returns the git version of this binary (e.g. "2015-01-11-284c030+").
// If the linker flags were not provided, the return value is "unknown".
func version() string {
	if GitInfo != "" {
		return GitInfo
	}

	return "unknown"
}
