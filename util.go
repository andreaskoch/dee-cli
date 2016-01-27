// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net"
	"os"
	"regexp"
	"strings"
)

// subDomainPattern defines a pattern for valid subdomain names.
// see:
// - http://stackoverflow.com/questions/7930751/regexp-for-subdomain
// - http://webmasters.stackexchange.com/questions/16996/maximum-domain-name-length
// - https://en.wikipedia.org/wiki/Hostname#Restrictions_on_valid_host_names
var subDomainPattern = regexp.MustCompile(`^(?:[A-Za-z0-9][A-Za-z0-9\-]{0,61}[A-Za-z0-9]|[A-Za-z0-9])$`)

// isEmpty returns true if the given text is empty or contains
// nothing but white space characters.
func isEmpty(text string) bool {
	return strings.TrimSpace(text) == ""
}

// isValidDomain returns true if the given domain name is valid; otherwise false.
// Note: This is not a real validation. I just want to exclude total garbage.
func isValidDomain(domain string) bool {
	if len(domain) > 255 {
		// too long.
		return false
	}

	return isEmpty(domain) == false
}

// isValidSubdomain returns true if the given subdomain name is valid; otherwise false.
// Note: This is not a real validation. I just want to exclude total garbage.
func isValidSubdomain(subdomain string) bool {
	if len(subdomain) > 253 {
		// too long
		return false
	}

	// each part must be valid (if there are multiple parts)
	for _, part := range strings.Split(subdomain, ".") {
		if len(part) > 63 {
			// too long
			return false
		}

		if !subDomainPattern.MatchString(part) {
			return false
		}
	}

	return true
}

// getDNSRecordTypeByIP returns the DNS record type for the given IP.
// It will return "A" for an IPv4 address and "AAAA" for an IPv6 address.
func getDNSRecordTypeByIP(ip net.IP) string {
	if ip.To4() == nil {
		return "AAAA"
	}

	return "A"
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
