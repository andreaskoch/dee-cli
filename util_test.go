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

func Test_isValidSubdomain_GivenTextIsValid_ResultIsTrue(t *testing.T) {

	// arrange
	inputs := []string{
		"www",
		"w-w-w",
		"w.w.w",
		"a",
		"1",
		"123",
		"abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk",
		"abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk",
	}

	for _, input := range inputs {
		// act
		result := isValidSubdomain(input)

		// assert
		if result == false {
			t.Fail()
			t.Logf("isValidSubdomain(%q) should have returned true", input)
		}
	}
}

func Test_isValidSubdomain_GivenTextIsInvalid_ResultIsFalse(t *testing.T) {

	// arrange
	inputs := []string{
		" www",
		"www ",
		"w ww",
		"-a",
		"-hi-",
		"_hi_",
		"*hi*",
		"abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk",
	}

	for _, input := range inputs {
		// act
		result := isValidSubdomain(input)

		// assert
		if result == true {
			t.Fail()
			t.Logf("isValidSubdomain(%q) should have returned false", input)
		}
	}
}
