// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"testing"
)

func Test_PrintUsageInformation_ResultIsNotEmpty(t *testing.T) {
	// arrange
	actions := []action{
		testAction{
			name:           "login",
			description:    "login",
			usage:          "yada yada login",
			executeMessage: "success",
		},
		testAction{
			name:           "logout",
			description:    "logout",
			usage:          "yada yada logout",
			executeMessage: "success",
		},
	}
	usagePrinter := newUsagePrinter("dnsimple-cli", "v0.1.0", actions)

	// act
	buf := new(bytes.Buffer)
	usagePrinter.PrintUsageInformation(buf)

	// assert
	if len(buf.String()) == 0 {
		t.Fail()
		t.Logf("PrintUsageInformation did not write to the given writer.")
	}
}
