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
	"strings"

	"github.com/serulian/compiler/tester"

	"github.com/spf13/cobra"
)

var (
	browsers []string
)

// karmaPackages defines the NPM packages necessary to run tests via the Karma runner.
var karmaPackages = []string{"karma", "jasmine-core", "karma-jasmine", "karma-sourcemap-loader"}

// karmaRunnerPackages defines the packages necessary for each browser runner.
var karmaRunnerPackages = map[string]string{
	"Chrome": "karma-chrome-launcher",
}

// karmaConfigTemplate defines a template for emitting the config ES code.
const karmaConfigTemplate = `
module.exports = function(config) {
  config.set(%s);
};
`

// karmaConfigFilename is the filename at which the generated Karma configuration will be emitted.
const karmaConfigFilename = "karma.conf.js"

// karmaTestRunner defines the karma test runner.
type karmaTestRunner struct{}

func (ktr *karmaTestRunner) Title() string {
	return "Karma"
}

func (ktr *karmaTestRunner) DecorateCommand(command *cobra.Command) {
	command.PersistentFlags().StringSliceVar(&browsers, "browser", []string{"Chrome"},
		"If specified, browsers to use to run the tests")
}

func (ktr *karmaTestRunner) SetupIfNecessary(testingEnvDirectoryPath string) error {
	log.Printf("Verifying installation of Yarn")

	// Ensure Yarn is installed and available.
	_, err := runCommand(testingEnvDirectoryPath, "yarn", "--version")
	if err != nil {
		return err
	}

	// Check for a yarn.lock file in the testing environment directory. If none, generate an empty one.
	lockFilePath := path.Join(testingEnvDirectoryPath, "yarn.lock")
	_, err = ioutil.ReadFile(lockFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		err := ioutil.WriteFile(lockFilePath, []byte(""), 0777)
		if err != nil {
			return err
		}
	}

	// Collect the packages that need to be installed.
	var packages = make([]string, 0)
	for _, packageName := range karmaPackages {
		packages = append(packages, packageName)
	}

	for _, browser := range browsers {
		browserPackage, exists := karmaRunnerPackages[browser]
		if !exists {
			return fmt.Errorf("Unknown browser %s", browser)
		}

		packages = append(packages, browserPackage)
	}

	// Filter out any packages already installed.
	log.Printf("Verifying installation of Yarn packages")
	var requiredPackages = make([]string, 0, len(packages))

	args := append([]string{"list"}, packages...)
	output, err := runCommand(testingEnvDirectoryPath, "yarn", args...)
	if err != nil {
		return err
	}

	for _, packageName := range packages {
		if !strings.Contains(output, packageName+"@") {
			requiredPackages = append(requiredPackages, packageName)
		}
	}

	// Check if we need to install some packages.
	if len(requiredPackages) > 0 {
		fmt.Printf("The following Yarn packages are required to coninue: %v. Install now? [Y/n]\n", requiredPackages)
		var input string
		fmt.Scanln(&input)

		if strings.ToLower(input) != "y" && input != "" {
			return fmt.Errorf("Installation of required packages canceled")
		}

		// Install the necessary packages.
		log.Printf("Installing Yarn packages: %s\n", requiredPackages)
		args := append([]string{"add"}, requiredPackages...)
		err := runStreamingCommand(testingEnvDirectoryPath, "yarn", args...)
		if err != nil {
			return err
		}
	}

	return nil
}

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

func (ktr *karmaTestRunner) Run(testingEnvDirectoryPath string, generatedFilePath string) (bool, error) {
	generatedDirectory := path.Dir(generatedFilePath)
	config := karmaConfig{
		BasePath:   generatedDirectory,
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

	// Write karma configuration into the testing environment directory.
	karmaConfigPath := path.Join(testingEnvDirectoryPath, karmaConfigFilename)
	configStr := fmt.Sprintf(karmaConfigTemplate, string(marshalledConfig))
	err = ioutil.WriteFile(karmaConfigPath, []byte(configStr), 0777)
	if err != nil {
		log.Fatal(err)
	}

	// Execute karma over the config.
	log.Printf("Running tests via Karma...")
	err = runStreamingCommand(testingEnvDirectoryPath, "node_modules/karma/bin/karma", "start", karmaConfigFilename)
	if err != nil {
		// Handle exits that are expected.
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func runStreamingCommand(workDir string, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = workDir

	err := cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func runCommand(workDir string, command string, args ...string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Stdout = &out
	cmd.Dir = workDir

	err := cmd.Run()
	return out.String(), err
}

func init() {
	tester.RegisterRunner("karma", &karmaTestRunner{})
}
