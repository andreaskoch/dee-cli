// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
)

// newUsagePrinter creates a new instance of the usage printer.
func newUsagePrinter(executableName string, version string, actions []action) usagePrinter {
	return usagePrinter{executableName, version, actions}
}

// usagePrinter prints usage information for the command line utility.
type usagePrinter struct {
	executableName string
	version        string
	actions        []action
}

// PrintUsageInformation prints the applications usage information for all available actions
// into the given output writer.
func (printer usagePrinter) PrintUsageInformation(output io.Writer) {
	fmt.Fprintf(output, "%s updates DNS records via DNSimple.\n", printer.executableName)
	fmt.Fprintf(output, "\n")

	fmt.Fprintf(output, "Version: %s\n", printer.version)
	fmt.Fprintf(output, "\n")

	fmt.Fprintf(output, "Usage:\n")
	fmt.Fprintf(output, "\n")
	fmt.Fprintf(output, "  %s <action> [arguments ...]\n", printer.executableName)
	fmt.Fprintf(output, "\n")

	// List of all actions
	fmt.Fprintf(output, "Actions:\n")

	for _, action := range printer.actions {
		fmt.Fprintf(output, "%10s  %s\n", action.Name(), action.Description())
	}

	fmt.Fprintf(output, "\n")

	// Action details
	for _, action := range printer.actions {
		fmt.Fprintf(output, "Action: %s\n", action.Name())
		fmt.Fprintf(output, "%s\n", action.Usage())
	}
}
