// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/andreaskoch/dee-ns"
)

var (
	actionNameDelete = "delete"

	deleteAddressRecordArguments = flag.NewFlagSet(actionNameDelete, flag.ContinueOnError)
	deleteDomain                 = deleteAddressRecordArguments.String("domain", "", "Domain (e.g. example.com)")
	deleteSubdomain              = deleteAddressRecordArguments.String("subdomain", "", "Subdomain (e.g. www)")
	deleteRecordType             = deleteAddressRecordArguments.String("type", "", "The address record type (e.g. \"AAAA\")")
)

type deleteAction struct {
	dnsEditorFactory dnsEditorCreator
}

func (action deleteAction) Name() string {
	return actionNameDelete
}

func (action deleteAction) Description() string {
	return "Delete an address record"
}

func (action deleteAction) Usage() string {
	buf := new(bytes.Buffer)
	deleteAddressRecordArguments.SetOutput(buf)
	deleteAddressRecordArguments.PrintDefaults()
	return buf.String()
}

// Execute deletes the DNS record of the domain given from the supplied arguments.
// If the delete fails an error is returned.
func (action deleteAction) Execute(arguments []string) (message, error) {

	// parse the arguments
	*deleteDomain = ""
	*deleteSubdomain = ""
	*deleteRecordType = ""
	if parseError := deleteAddressRecordArguments.Parse(arguments); parseError != nil {
		return nil, parseError
	}

	// domain
	if *deleteDomain == "" {
		return nil, fmt.Errorf("No domain supplied")
	}

	// subdomain
	if *deleteSubdomain == "" {
		return nil, fmt.Errorf("No subdomain supplied")
	}

	// subdomain
	if *deleteRecordType == "" {
		return nil, fmt.Errorf("No record type supplied")
	}

	// create a DNS editor
	var addressRecordDeleter deens.DNSRecordDeleter
	addressRecordDeleter, dnsEditorError := action.dnsEditorFactory.CreateDNSEditor()
	if dnsEditorError != nil {
		return nil, fmt.Errorf("Cannot create DNS editor: %s", dnsEditorError.Error())
	}

	deleteError := addressRecordDeleter.DeleteSubdomain(*deleteDomain, *deleteSubdomain, *deleteRecordType)
	if deleteError != nil {
		return nil, fmt.Errorf("%s", deleteError.Error())
	}

	return successMessage{fmt.Sprintf("Deleted: %s.%s (%s)", *deleteSubdomain, *deleteDomain, *deleteRecordType)}, nil
}
