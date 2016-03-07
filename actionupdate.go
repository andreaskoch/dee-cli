// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/andreaskoch/dee-ns"
	"net"
	"os"
)

var (
	actionNameUpdate = "update"

	updateAddressRecordArguments = flag.NewFlagSet(actionNameUpdate, flag.ContinueOnError)
	updateDomain                 = updateAddressRecordArguments.String("domain", "", "Domain (e.g. example.com)")
	updateSubdomain              = updateAddressRecordArguments.String("subdomain", "", "Subdomain (e.g. www)")
	updateIP                     = updateAddressRecordArguments.String("ip", "", "IP address (e.g. ::1, 127.0.0.1)")
)

type updateAction struct {
	dnsEditorFactory dnsEditorCreator
	stdin            *os.File
}

func (action updateAction) Name() string {
	return actionNameUpdate
}

func (action updateAction) Description() string {
	return "Update the address record for a given sub domain"
}

func (action updateAction) Usage() string {
	buf := new(bytes.Buffer)
	updateAddressRecordArguments.SetOutput(buf)
	updateAddressRecordArguments.PrintDefaults()
	return buf.String()
}

// Execute updates the DNS record of the domain given from the supplied arguments.
// If the update fails an error is returned.
func (action updateAction) Execute(arguments []string) (message, error) {

	// parse the arguments
	*updateDomain = ""
	*updateSubdomain = ""
	*updateIP = ""
	if parseError := updateAddressRecordArguments.Parse(arguments); parseError != nil {
		return nil, parseError
	}

	// domain
	if *updateDomain == "" {
		return nil, fmt.Errorf("No domain supplied")
	}

	// take ip from stdin
	if *updateIP == "" && stdinHasData(action.stdin) {
		ipAddressFromStdin := ""
		fmt.Fscanf(os.Stdin, "%s", &ipAddressFromStdin)
		updateIP = &ipAddressFromStdin
	}

	if *updateIP == "" {
		return nil, fmt.Errorf("No IP address supplied")
	}

	ip := net.ParseIP(*updateIP)
	if ip == nil {
		return nil, fmt.Errorf("Cannot parse IP %q", ip)
	}

	// create a DNS editor
	var addressRecordUpdater deens.DNSRecordUpdater
	addressRecordUpdater, dnsEditorError := action.dnsEditorFactory.CreateDNSEditor()
	if dnsEditorError != nil {
		return nil, fmt.Errorf("Cannot create DNS editor: %s", dnsEditorError.Error())
	}

	updateError := addressRecordUpdater.UpdateSubdomain(*updateDomain, *updateSubdomain, ip)
	if updateError != nil {
		return nil, fmt.Errorf("%s", updateError.Error())
	}

	return successMessage{fmt.Sprintf("Updated: %s → %s", getFormattedDomainName(*updateSubdomain, *updateDomain), ip.String())}, nil
}
