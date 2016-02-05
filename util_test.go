// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func Test_isEmpty_EmptyString_ResultIsTrue(t *testing.T) {
	// arrange
	inputs := []string{
		"",
		" ",
		"    ",
		" ",
		" ",
		" ",
	}

	// act
	for _, input := range inputs {
		result := isEmpty(input)

		// assert
		if result == false {
			t.Fail()
			t.Logf("isEmpty(%q) should return true", input)
		}
	}
}

func Test_isEmpty_NotEmptyString_ResultIsFalse(t *testing.T) {
	// arrange
	inputs := []string{
		"-",
		".",
		" a ",
		" _ ",
	}

	// act
	for _, input := range inputs {
		result := isEmpty(input)

		// assert
		if result == true {
			t.Fail()
			t.Logf("isEmpty(%q) should return false", input)
		}
	}
}
