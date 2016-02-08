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
	actionNameCreate = "create"

	createAddressRecordArguments = flag.NewFlagSet(actionNameCreate, flag.ContinueOnError)
	createDomain                 = createAddressRecordArguments.String("domain", "", "Domain (e.g. example.com)")
	createSubdomain              = createAddressRecordArguments.String("subdomain", "", "Subdomain (e.g. www)")
	createIP                     = createAddressRecordArguments.String("ip", "", "IP address (e.g. ::1, 127.0.0.1)")
	createTTL                    = createAddressRecordArguments.Int("ttl", 600, "The time to live in seconds")
)

type createAction struct {
	dnsEditorFactory dnsEditorCreator
	stdin            *os.File
}

func (action createAction) Name() string {
	return actionNameCreate
}

func (action createAction) Description() string {
	return "Create an address record"
}

func (action createAction) Usage() string {
	buf := new(bytes.Buffer)
	createAddressRecordArguments.SetOutput(buf)
	createAddressRecordArguments.PrintDefaults()
	return buf.String()
}

// Execute creates the DNS record of the domain given from the supplied arguments.
// If the create fails an error is returned.
func (action createAction) Execute(arguments []string) (message, error) {

	// parse the arguments
	*createDomain = ""
	*createSubdomain = ""
	*createIP = ""
	*createTTL = 0
	if parseError := createAddressRecordArguments.Parse(arguments); parseError != nil {
		return nil, parseError
	}

	// domain
	if *createDomain == "" {
		return nil, fmt.Errorf("No domain supplied")
	}

	// TTL
	if *createTTL < 0 {
		return nil, fmt.Errorf("The given TTL cannot be negative")
	}

	// take ip from stdin
	if *createIP == "" && stdinHasData(action.stdin) {
		ipAddressFromStdin := ""
		fmt.Fscanf(os.Stdin, "%s", &ipAddressFromStdin)
		createIP = &ipAddressFromStdin
	}

	if *createIP == "" {
		return nil, fmt.Errorf("No IP address supplied")
	}

	ip := net.ParseIP(*createIP)
	if ip == nil {
		return nil, fmt.Errorf("Cannot parse IP %q", ip)
	}

	// create a DNS editor
	var addressRecordCreator deens.DNSRecordCreator
	addressRecordCreator, dnsEditorError := action.dnsEditorFactory.CreateDNSEditor()
	if dnsEditorError != nil {
		return nil, fmt.Errorf("Cannot create DNS editor: %s", dnsEditorError.Error())
	}

	createError := addressRecordCreator.CreateSubdomain(*createDomain, *createSubdomain, *createTTL, ip)
	if createError != nil {
		return nil, fmt.Errorf("%s", createError.Error())
	}

	return successMessage{fmt.Sprintf("Created: %s → %s", getFormattedDomainName(*createSubdomain, *createDomain), ip.String())}, nil
}
