// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package tester implements support for testing Serulian code.
package tester

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/serulian/compiler/builder"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/generator/es5"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// runners defines the map of test runners by name.
var runners = map[string]TestRunner{}

// TestRunner defines an interface for the test runner.
type TestRunner interface {
	// Title is a human-readable title for the test runner.
	Title() string

	// DecorateCommand decorates the cobra command for the runner with the runner-specific
	// options.
	DecorateCommand(command *cobra.Command)

	// SetupIfNecessary is run before any test runs occur to run the setup process
	// for the runner (if necessary). This method should no-op if all necessary
	// dependencies are in place.
	SetupIfNecessary() error

	// Run runs the test runner over the generated ES path.
	Run(generatedFilePath string) (bool, error)
}

// runTestsViaRunner runs all the tests at the given source path via the runner.
func runTestsViaRunner(runner TestRunner, path string) bool {
	log.Printf("Starting test run of %s via %v runner", path, runner.Title())

	// Run setup for the runner.
	err := runner.SetupIfNecessary()
	if err != nil {
		errHighlight := color.New(color.FgRed, color.Bold)
		errHighlight.Print("ERROR: ")

		text := color.New(color.FgWhite)
		text.Printf("Could not setup %s runner: %v\n", runner.Title(), err)
		return false
	}

	// Iterate over each test file in the source path. For each, compile the code into
	// JS at a temporary location and then pass the temporary location to the test
	// runner.
	return compilerutil.WalkSourcePath(path, func(currentPath string, info os.FileInfo) (bool, error) {
		if !strings.HasSuffix(info.Name(), "_test"+parser.SERULIAN_FILE_EXTENSION) {
			return false, nil
		}

		success, err := buildAndRunTests(currentPath, runner)
		if err != nil {
			return true, err
		}

		if success {
			return true, nil
		} else {
			return true, fmt.Errorf("Failure in test of file %s", currentPath)
		}
	}, packageloader.SerulianPackageDirectory)
}

// buildAndRunTests builds the source found at the given path and then runs its tests via the runner.
func buildAndRunTests(filePath string, runner TestRunner) (bool, error) {
	log.Printf("Building %s...", filePath)

	filename := path.Base(filePath)

	scopeResult := scopegraph.ParseAndBuildScopeGraph(filePath,
		[]string{},
		builder.CORE_LIBRARY)

	if !scopeResult.Status {
		// TODO: better output
		return false, fmt.Errorf("Compilation errors for test %s: %v", filePath, scopeResult.Errors)
	}

	// Generate the source.
	generated, sourceMap, err := es5.GenerateES5(scopeResult.Graph, filename+".js", "")
	if err != nil {
		log.Fatal(err)
	}

	// Save the source (with an adjusted call), in a temporary directory.
	moduleName := filename[0 : len(filename)-len(parser.SERULIAN_FILE_EXTENSION)]
	adjusted := fmt.Sprintf(`
		%s

		window.Serulian.then(function(global) {
			global.%s.TEST().then(function(a) {
			}).catch(function(err) {
		    throw err;     
		  })
		})

		//# sourceMappingURL=/%s.js.map
	`, generated, moduleName, filename)

	dir, err := ioutil.TempDir("", "testing")
	if err != nil {
		log.Fatal(err)
	}

	// Clean up once complete.
	defer os.RemoveAll(dir)

	// Write the source and map into the directory.
	marshalled, err := sourceMap.Build().Marshal()
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(path.Join(dir, filename+".js"), []byte(adjusted), 0777)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(path.Join(dir, filename+".js.map"), marshalled, 0777)
	if err != nil {
		log.Fatal(err)
	}

	// Call the runner with the test file.
	return runner.Run(path.Join(dir, filename+".js"))
}

// DecorateRunners decorates the test command with a command for each runner.
func DecorateRunners(command *cobra.Command) {
	for name, runner := range runners {
		var runnerCmd = &cobra.Command{
			Use:   fmt.Sprintf("%s [source path]", name),
			Short: "Runs the tests defined at the given source path via " + runner.Title(),
			Long:  "Runs the tests found in any *_test.seru files at the given source path",
			Run: func(cmd *cobra.Command, args []string) {
				if len(args) != 1 {
					fmt.Println("Expected source path")
					os.Exit(-1)
				}

				if runTestsViaRunner(runner, args[0]) {
					os.Exit(1)
				} else {
					os.Exit(-1)
				}
			},
		}

		runner.DecorateCommand(runnerCmd)
		command.AddCommand(runnerCmd)
	}
}

// RegisterRunner registers a test runner with the specific name.
func RegisterRunner(name string, runner TestRunner) {
	if name == "" {
		panic("Test runner must have a name")
	}

	if runner == nil {
		panic("Cannot register nil runner")
	}

	if _, exists := runners[name]; exists {
		panic(fmt.Sprintf("Test runner with name %s already exists", name))
	}

	runners[name] = runner
}
