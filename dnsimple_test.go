// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type testAction struct {
	name           string
	description    string
	usage          string
	executeMessage string
}

func (action testAction) Name() string {
	return action.name
}

func (action testAction) Description() string {
	return action.description
}

func (action testAction) Usage() string {
	return action.usage
}

func (action testAction) Execute(arguments []string) (message, error) {
	return successMessage{action.executeMessage}, nil
}
