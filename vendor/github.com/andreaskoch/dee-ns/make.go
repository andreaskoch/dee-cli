// +build ignore

// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program builds dee.
//
// $ go run make.go -install
//
// View the README.md for further details.
//
// The output binaries go into the ./bin/ directory (under the GOPATH, where make.go is)
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ProjectName contains the name of the project
const ProjectName = "dee-ns"

// GOPATH environment variable name.
const GOPATH_ENVIRONMENT_VARIABLE = "GOPATH"

// GO15VENDOREXPERIMENT environment variable name
const GO15VENDOREXPERIMENT_ENVIRONMENT_VARIBALE = "GO15VENDOREXPERIMENT"

var (

	// command line flags
	verboseFlagIsSet  = flag.Bool("v", false, "Verbose mode")
	buildFlagIsSet    = flag.Bool("build", false, fmt.Sprintf("Builds the %s package", ProjectName))
	coverageFlagIsSet = flag.Bool("coverage", false, "Run the test and create a code coverage report")

	// The GOPATH for the current project
	goPath = getWorkingDirectory()

	// The GOBIN for the current project
	goBin = filepath.Join(goPath, "bin")
)

// Compilation Target Definition
type compilationTarget struct {
	OperatingSystem string
	Architecture    string
	OtherVariables  []string
}

func (target *compilationTarget) String() string {
	return fmt.Sprintf("%s/%s", target.OperatingSystem, target.Architecture)
}

func init() {

	executableName := "make.go"

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s provides functions for compiling %s.\n", executableName, ProjectName)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  go run %s [options]\n", executableName)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		return
	}

	if *verboseFlagIsSet {
		fmt.Printf("%s: %s\n", GOPATH_ENVIRONMENT_VARIABLE, goPath)
	}

	if *buildFlagIsSet {
		if err := build(); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		return
	}

	if *coverageFlagIsSet {
		if err := codecoverage(); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		return
	}

	flag.Usage()
}

// Build the package.
func build() error {

	// prepare the environment variables
	environmentVariables := cleanGoEnv()
	environmentVariables = setEnv(environmentVariables, GO15VENDOREXPERIMENT_ENVIRONMENT_VARIBALE, "1")

	return runCommand(os.Stdout, os.Stderr, goPath, environmentVariables, "go", "build")

}

// codecoverage runs all tests and creates a code coverage report
func codecoverage() error {

	// prepare the environment variables
	environmentVariables := cleanGoEnv()
	environmentVariables = setEnv(environmentVariables, GO15VENDOREXPERIMENT_ENVIRONMENT_VARIBALE, "1")

	coverageError := runCommand(os.Stdout, os.Stderr, goPath, environmentVariables, "go", "test", "-coverprofile=coverage.out")
	if coverageError != nil {
		return coverageError
	}

	reportError := runCommand(os.Stdout, os.Stderr, goPath, environmentVariables, "go", "tool", "cover", "-html=coverage.out", "-o", "coverage.html")
	if reportError != nil {
		return reportError
	}

	return nil
}

// getWorkingDirectory returns the current working directory path or fails.
func getWorkingDirectory() string {
	goPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	return goPath
}

// Execute go in the specified go path with the supplied command arguments.
func runCommand(stdout, stderr io.Writer, workingDirectory string, environmentVariables []string, command string, args ...string) error {

	// Create the command
	cmdName := fmt.Sprintf("%s %s", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)

	// Set the working directory
	cmd.Dir = workingDirectory

	// set environment variables
	cmd.Env = environmentVariables

	// Capture the output
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if *verboseFlagIsSet {
		log.Printf("Running %s", cmdName)
	}

	// execute the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error running %s: %v", cmdName, err.Error())
	}

	return nil
}

// cleanGoEnv returns a copy of the current environment.
func cleanGoEnv() (clean []string) {
	return os.Environ()
}

// setEnv sets the given key & value in the provided environment.
// Each value in the env list should be of the form key=value.
func setEnv(env []string, key, value string) []string {
	for i, s := range env {
		if strings.HasPrefix(s, fmt.Sprintf("%s=", key)) {
			env[i] = envPair(key, value)
			return env
		}
	}
	env = append(env, envPair(key, value))
	return env
}

// Create an environment variable of the form key=value.
func envPair(key, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}
