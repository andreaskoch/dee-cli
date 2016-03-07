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
	actionNameCreateOrUpdate = "createorupdate"

	createOrUpdateAddressRecordArguments = flag.NewFlagSet(actionNameCreateOrUpdate, flag.ContinueOnError)
	createOrUpdateDomain                 = createOrUpdateAddressRecordArguments.String("domain", "", "Domain (e.g. example.com)")
	createOrUpdateSubdomain              = createOrUpdateAddressRecordArguments.String("subdomain", "", "Subdomain (e.g. www)")
	createOrUpdateIP                     = createOrUpdateAddressRecordArguments.String("ip", "", "IP address (e.g. ::1, 127.0.0.1)")
	createOrUpdateTTL                    = createOrUpdateAddressRecordArguments.Int("ttl", defaultTTL, "The time to live in seconds")
)

type createOrUpdateAction struct {
	dnsEditorFactory    dnsEditorCreator
	infoProviderFactory dnsInfoProviderCreator
	stdin               *os.File
}

func (action createOrUpdateAction) Name() string {
	return actionNameCreateOrUpdate
}

func (action createOrUpdateAction) Description() string {
	return "Create or update an address record"
}

func (action createOrUpdateAction) Usage() string {
	buf := new(bytes.Buffer)
	createAddressRecordArguments.SetOutput(buf)
	createAddressRecordArguments.PrintDefaults()
	return buf.String()
}

// Execute creates the DNS record of the domain given from the supplied arguments.
// If the create fails an error is returned.
func (action createOrUpdateAction) Execute(arguments []string) (message, error) {

	// parse the arguments
	*createOrUpdateDomain = ""
	*createOrUpdateSubdomain = ""
	*createOrUpdateIP = ""
	*createOrUpdateTTL = defaultTTL
	if parseError := createOrUpdateAddressRecordArguments.Parse(arguments); parseError != nil {
		return nil, parseError
	}

	// domain
	if *createOrUpdateDomain == "" {
		return nil, fmt.Errorf("No domain supplied")
	}

	// subdomain
	noSubdomainGiven := *createOrUpdateSubdomain == ""

	// TTL
	if *createOrUpdateTTL < 0 {
		return nil, fmt.Errorf("The given TTL cannot be negative")
	}

	// take ip from stdin
	if *createOrUpdateIP == "" && stdinHasData(action.stdin) {
		ipAddressFromStdin := ""
		fmt.Fscanf(os.Stdin, "%s", &ipAddressFromStdin)
		createIP = &ipAddressFromStdin
	}

	if *createOrUpdateIP == "" {
		return nil, fmt.Errorf("No IP address supplied")
	}

	ip := net.ParseIP(*createOrUpdateIP)
	if ip == nil {
		return nil, fmt.Errorf("Cannot parse IP %q", ip)
	}

	// create a DNS editor
	var addressRecordEditor deens.DNSRecordEditor
	addressRecordEditor, dnsEditorError := action.dnsEditorFactory.CreateDNSEditor()
	if dnsEditorError != nil {
		return nil, fmt.Errorf("Cannot create DNS editor: %s", dnsEditorError.Error())
	}

	// info provider
	infoProvider, infoProviderError := action.getInfoProvider()
	if infoProviderError != nil {
		return nil, fmt.Errorf("No DNS info provider available")
	}

	// determine the record type
	dnsRecordType := getDNSRecordTypeByIP(ip)

	_, domainRecordError := infoProvider.GetSubdomainRecord(*createOrUpdateDomain, *createOrUpdateSubdomain, dnsRecordType)
	if domainRecordError == nil || noSubdomainGiven {

		// update
		updateError := addressRecordEditor.UpdateSubdomain(*createOrUpdateDomain, *createOrUpdateSubdomain, ip)
		if updateError != nil {
			return nil, fmt.Errorf("%s", updateError.Error())
		}

		return successMessage{fmt.Sprintf("Updated: %s → %s", getFormattedDomainName(*createOrUpdateSubdomain, *createOrUpdateDomain), ip.String())}, nil

	}

	// create
	createError := addressRecordEditor.CreateSubdomain(*createOrUpdateDomain, *createOrUpdateSubdomain, *createOrUpdateTTL, ip)
	if createError != nil {
		return nil, fmt.Errorf("%s", createError.Error())
	}

	return successMessage{fmt.Sprintf("Created: %s → %s", getFormattedDomainName(*createOrUpdateSubdomain, *createOrUpdateDomain), ip.String())}, nil
}

// getInfoProvider returns a DNS info provider instance or an error if the creation of the provider failed.
func (action createOrUpdateAction) getInfoProvider() (deens.DNSInfoProvider, error) {
	if action.infoProviderFactory == nil {
		return nil, fmt.Errorf("No DNS info provider factory available")
	}

	return action.infoProviderFactory.CreateInfoProvider()
}
