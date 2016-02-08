// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"
)

// isEmpty returns true if the given text is empty or contains
// nothing but white space characters.
func isEmpty(text string) bool {
	return strings.TrimSpace(text) == ""
}

// stdinHasData returns true if there is data avaialble in the given file (os.Stdin), otherwise false.
// see: http://stackoverflow.com/questions/22744443/check-if-there-is-something-to-read-on-stdin-in-golang
func stdinHasData(stdin *os.File) bool {
	if stdin == nil {
		return false
	}

	stat, _ := stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return true
	}

	return false
}

// getFormattedDomainName returns the formatted domain name for
// the given subdomain and domain names.
func getFormattedDomainName(subdomain, domain string) string {
	if domain == "" {
		return ""
	}

	if subdomain == "" {
		return domain
	}

	return fmt.Sprintf("%s.%s", subdomain, domain)
}
