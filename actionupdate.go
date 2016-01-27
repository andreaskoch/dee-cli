// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
)

var (
	actionNameUpdate = "update"

	updateSubdomainArguments = flag.NewFlagSet(actionNameUpdate, flag.ContinueOnError)
	updateDomain             = updateSubdomainArguments.String("domain", "", "Domain (e.g. example.com")
	updateSubdomain          = updateSubdomainArguments.String("subdomain", "", "Subdomain (e.g. www)")
	updateIP                 = updateSubdomainArguments.String("ip", "", "IP address (e.g. 127.0.0.1, ::1)")
)

type updateAction struct {
	domainUpdater updater
	stdin         *os.File
}

func (action updateAction) Name() string {
	return actionNameUpdate
}

func (action updateAction) Description() string {
	return "Update the DNS record for a given sub domain"
}

func (action updateAction) Usage() string {
	buf := new(bytes.Buffer)
	updateSubdomainArguments.SetOutput(buf)
	updateSubdomainArguments.PrintDefaults()
	return buf.String()
}

// Execute updates the DNS record of the domain given from the supplied arguments.
// If the update fails an error is returned.
func (action updateAction) Execute(arguments []string) (message, error) {

	// parse the arguments
	if parseError := updateSubdomainArguments.Parse(arguments); parseError != nil {
		return nil, parseError
	}

	// domain
	if *updateDomain == "" {
		return nil, fmt.Errorf("No domain supplied.")
	}

	// subdomain
	if *updateSubdomain == "" {
		return nil, fmt.Errorf("No subdomain supplied.")
	}

	// take ip from stdin
	if *updateIP == "" && stdinHasData(action.stdin) {
		ipAddressFromStdin := ""
		fmt.Fscanf(os.Stdin, "%s", &ipAddressFromStdin)
		updateIP = &ipAddressFromStdin
	}

	if *updateIP == "" {
		return nil, fmt.Errorf("No IP address supplied.")
	}

	ip := net.ParseIP(*updateIP)
	updateError := action.domainUpdater.UpdateSubdomain(*updateDomain, *updateSubdomain, ip)
	if updateError != nil {
		return nil, fmt.Errorf("%s", updateError.Error())
	}

	return successMessage{fmt.Sprintf("Updated: %s.%s → %s", *updateSubdomain, *updateDomain, ip.String())}, nil
}
