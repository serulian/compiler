// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package karma implements support for testing Serulian code via the karma testing framework.
package karma

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/serulian/compiler/tester"

	"github.com/spf13/cobra"
)

var (
	browsers []string
)

// _KARMA_PACKAGES defines the NPM packages necessary to run tests via the Karma runner.
var _KARMA_PACKAGES = []string{"karma", "jasmine-core", "karma-jasmine", "karma-sourcemap-loader"}

// _KARMA_RUNNER_PACKAGES defines the packages necessary for each browser runner.
var _KARMA_RUNNER_PACKAGES = map[string]string{
	"Chrome": "karma-chrome-launcher",
}

// karmaTestRunner defines the karma test runner.
type karmaTestRunner struct{}

func (ktr *karmaTestRunner) Title() string {
	return "Karma"
}

func (ktr *karmaTestRunner) DecorateCommand(command *cobra.Command) {
	command.PersistentFlags().StringSliceVar(&browsers, "browser", []string{"Chrome"},
		"If specified, browsers to use to run the tests")

}

func (ktr *karmaTestRunner) SetupIfNecessary() error {
	log.Printf("Verifying installation of NPM")

	// Ensure NPM is installed and available.
	_, err := runCommand("npm", "-v")
	if err != nil {
		return err
	}

	// Collect the packages that need to be installed.
	var packages = make([]string, 0)
	for _, packageName := range _KARMA_PACKAGES {
		packages = append(packages, packageName)
	}

	for _, browser := range browsers {
		browserPackage, exists := _KARMA_RUNNER_PACKAGES[browser]
		if !exists {
			return fmt.Errorf("Unknown browser %s", browser)
		}

		packages = append(packages, browserPackage)
	}

	// Filter out any packages already installed.
	log.Printf("Verifying installation of NPM packages")
	var requiredPackages = make([]string, 0, len(packages))
	for _, packageName := range packages {
		log.Printf("Verifying NPM package %s...\n", packageName)
		_, err := runCommand("npm", "list", packageName)
		if err != nil {
			requiredPackages = append(requiredPackages, packageName)
		}
	}

	// Check if we need to install some packages.
	if len(requiredPackages) > 0 {
		fmt.Printf("The following NPM packages are required to coninue: %v. Install now? [Y/n]\n", requiredPackages)
		var input string
		fmt.Scanln(&input)

		if input != "Y" && input != "" {
			return fmt.Errorf("Installation of required packages canceled")
		}

		// Install the necessary packages.
		for _, packageName := range requiredPackages {
			log.Printf("Installing NPM package %s...\n", packageName)
			err := runStreamingCommand("npm", "install", packageName, "--save-dev")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

const _KARMA_CONFIG_TEMPLATE = `
module.exports = function(config) {
  config.set(%s);
};
`

const _KARMA_CONFIG_FILENAME = "karma.conf.js"

// karmaConfig defines the structural view of Karma configuration.
type karmaConfig struct {
	BasePath      string                 `json:"basePath"`
	Frameworks    []string               `json:"frameworks"`
	AutoWatch     bool                   `json:"autoWatch"`
	Browsers      []string               `json:"browsers"`
	Files         []string               `json:"files"`
	SingleRun     bool                   `json:"singleRun"`
	Preprocessors map[string]interface{} `json:"preprocessors"`
}

func (ktr *karmaTestRunner) Run(generatedFilePath string) (bool, error) {
	directory := path.Dir(generatedFilePath)
	config := karmaConfig{
		BasePath:   ".",
		Frameworks: []string{"jasmine"},
		AutoWatch:  false,
		Browsers:   browsers,
		Files:      []string{"*_test.seru.js"},
		SingleRun:  true,
		Preprocessors: map[string]interface{}{
			"**/*.js": []string{"sourcemap"},
		},
	}

	marshalledConfig, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}

	// Write karma configuration into the directory of the generated path.
	karmaConfigPath := path.Join(directory, _KARMA_CONFIG_FILENAME)
	configStr := fmt.Sprintf(_KARMA_CONFIG_TEMPLATE, string(marshalledConfig))
	err = ioutil.WriteFile(karmaConfigPath, []byte(configStr), 0777)
	if err != nil {
		log.Fatal(err)
	}

	// Execute karma over the config.
	log.Printf("Running tests via Karma...")
	err = runStreamingCommand("node_modules/karma/bin/karma", "start", karmaConfigPath)
	if err != nil {
		// Handle exits that are expected.
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func runStreamingCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func runCommand(command string, args ...string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func init() {
	tester.RegisterRunner("karma", &karmaTestRunner{})
}
