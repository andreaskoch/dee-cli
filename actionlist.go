// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/pearkes/dnsimple"
	"strings"
	"text/tabwriter"
)

var (
	actionNameList = "list"

	listArguments = flag.NewFlagSet(actionNameList, flag.ContinueOnError)
	listDomain    = listArguments.String("domain", "", "Domain (optional")
	listSubdomain = listArguments.String("subdomain", "", "Subdomain (optional)")
)

type listAction struct {
	infoProviderFactory dnsInfoProviderFactory
}

func (action listAction) Name() string {
	return actionNameList
}

func (action listAction) Description() string {
	return "List all available domains, subdomains or DNS records"
}

func (action listAction) Usage() string {
	buf := new(bytes.Buffer)
	listArguments.SetOutput(buf)
	listArguments.PrintDefaults()
	return buf.String()
}

// Execute lists the list of all domains, subdomains or DNS records
// based on the supplied arguments.
func (action listAction) Execute(arguments []string) (message, error) {

	// parse the arguments
	*listDomain = ""
	*listSubdomain = ""

	if parseError := listArguments.Parse(arguments); parseError != nil {
		return nil, parseError
	}

	infoProvider, infoProviderError := action.getInfoProvider()
	if infoProviderError != nil {
		return nil, fmt.Errorf("No DNS info provider available")
	}

	domainParamIsSet := isEmpty(*listDomain) == false
	subdomainParamIsSet := isEmpty(*listSubdomain) == false

	// case: 2 get DNS records for the given subdomain
	if domainParamIsSet && subdomainParamIsSet {
		records, err := infoProvider.GetSubdomainDNSRecords(*listDomain, *listSubdomain)
		if err != nil {
			return nil, fmt.Errorf("Unable to fetch DNS records for subdomain %s.%s", *listSubdomain, *listDomain)
		}

		return successMessage{formatDNSRecords(records, *listDomain)}, nil
	}

	// case 3: get all subdomains
	if domainParamIsSet && !subdomainParamIsSet {
		records, err := infoProvider.GetAllDNSRecords(*listDomain)
		if err != nil {
			return nil, fmt.Errorf("Unable to fetch DNS records for domain %s", *listDomain)
		}

		return successMessage{formatDNSRecords(records, *listDomain)}, nil
	}

	// case 1: get all domain names
	names, err := infoProvider.GetDomainNames()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve domain names: %s", err.Error())
	}

	return successMessage{strings.Join(names, "\n")}, nil
}

// getInfoProvider returns a DNS info provider instance or an error if the creation of the provider failed.
func (action listAction) getInfoProvider() (dnsInfoProvider, error) {
	if action.infoProviderFactory == nil {
		return nil, fmt.Errorf("No DNS info provider factory available.")
	}

	infoProvider := action.infoProviderFactory.CreateInfoProvider()
	return infoProvider, nil
}

// formatDNSRecords takes a list of DNS records and formats them as a table.
func formatDNSRecords(records []dnsimple.Record, domainName string) string {
	buf := new(bytes.Buffer)

	// initialize the tabwriter
	w := new(tabwriter.Writer)
	minWidth := 0
	tabWidth := 8
	padding := 3
	w.Init(buf, minWidth, tabWidth, padding, ' ', 0)

	for index, record := range records {

		// assemble the subdomain / domain name
		domainName := domainName
		if !isEmpty(record.Name) {
			domainName = record.Name + "." + domainName
		}

		fmt.Fprintf(w, "%s\t%s\t%s", domainName, record.RecordType, record.Content)

		// append newline if we are not
		// formatting the last record
		if index < len(records)-1 {
			fmt.Fprintf(w, "\n")
		}

	}

	w.Flush()

	return buf.String()
}
