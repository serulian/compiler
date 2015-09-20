// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// packageloader package defines functions and types for loading Serulian source from disk or VCS.
package packageloader

import (
	"fmt"
	"io/ioutil"
	"path"
	"sync"

	"github.com/serulian/compiler/parser"
	"github.com/serulian/compiler/vcs"
)

// pathKind identifies the supported kind of paths
type pathKind int

const (
	pathSourceFile pathKind = iota
	pathLocalPackage
	pathVCSPackage
)

// serulianFileExtension is the extension for Serulian source files.
const serulianSourceFileExtension = ".seru"

// serulianPackageDirectory is the directory under the root directory holding cached packages.
const serulianPackageDirectory = ".pkg"

// pathInformation holds information about a path to load.
type pathInformation struct {
	kind   pathKind
	path   string
	source string
}

// Returns the string representation of the given path.
func (p *pathInformation) String() string {
	return fmt.Sprintf("%s://%s", p.kind, p.path)
}

// PackageLoader helps to fully and recursively load a Serulian package and its dependencies
// from a directory or set of directories.
type PackageLoader struct {
	rootSourceFile string // The root source file location.

	errors   chan error  // Errors are reported on this channel
	warnings chan string // Warnings are reported on this channel

	nodeBuilder parser.NodeBuilder // Builder to use for constructing the parse nodes

	pathsToLoad      chan pathInformation // The paths to load
	pathsEncountered map[string]bool      // The paths processed by the loader goroutine
	workTracker      sync.WaitGroup       // WaitGroup used to wait until all loading is complete
	finished         chan bool            // Channel used to tell background goroutines to quit
}

type LoadResult struct {
	Status   bool     // True on success, false otherwise
	Errors   []error  // The errors encountered, if any
	Warnings []string // The warnings encountered, if any
}

// NewPackageLoader creates and returns a new package loader for the given path.
func NewPackageLoader(rootSourceFile string, nodeBuilder parser.NodeBuilder) *PackageLoader {
	return &PackageLoader{
		rootSourceFile: rootSourceFile,

		errors:   make(chan error),
		warnings: make(chan string),

		nodeBuilder: nodeBuilder,

		pathsEncountered: map[string]bool{},
		pathsToLoad:      make(chan pathInformation),

		finished: make(chan bool, 2),
	}
}

// Load performs the loading of a Serulian package found at the directory path.
// Any dependencies will be loaded as well.
func (p *PackageLoader) Load() *LoadResult {
	// Start the loading goroutine.
	go p.loadAndParse()

	// Start the error/warning collection goroutine.
	result := &LoadResult{
		Status:   true,
		Errors:   make([]error, 0),
		Warnings: make([]string, 0),
	}

	go p.collectIssues(result)

	// Add the root source file as the first package to be parsed.
	p.pushPath(pathSourceFile, p.rootSourceFile, "(root)")

	// Wait for all packages and source files to be completed.
	p.workTracker.Wait()

	// Tell the goroutines to quit.
	p.finished <- true
	p.finished <- true
	return result
}

// pushPath adds a path to be processed by the package loader.
func (p *PackageLoader) pushPath(kind pathKind, path string, source string) {
	p.workTracker.Add(1)
	p.pathsToLoad <- pathInformation{kind, path, source}
}

// collectIssues watches the errors and warnings channels to collect those issues as they
// are added.
func (p *PackageLoader) collectIssues(result *LoadResult) {
	for {
		select {
		case newError := <-p.errors:
			result.Errors = append(result.Errors, newError)
			result.Status = false

		case newWarnings := <-p.warnings:
			result.Warnings = append(result.Warnings, newWarnings)

		case <-p.finished:
			return
		}
	}
}

// loadAndParse watches the pathsToLoad channel and parses load operations off to their
// own goroutines.
func (p *PackageLoader) loadAndParse() {
	for {
		select {
		case currentPath := <-p.pathsToLoad:
			go p.loadAndParsePath(currentPath)

		case <-p.finished:
			return
		}
	}
}

// loadAndParsePath parses or loads a specific path.
func (p *PackageLoader) loadAndParsePath(currentPath pathInformation) {
	defer p.workTracker.Done()

	// Ensure we have not already seen this path.
	pathKey := currentPath.String()
	if _, ok := p.pathsEncountered[pathKey]; ok {
		return
	}

	p.pathsEncountered[pathKey] = true

	// Perform parsing/loading.
	switch currentPath.kind {
	case pathSourceFile:
		p.conductParsing(currentPath)

	case pathLocalPackage:
		p.loadLocalPackage(currentPath)

	case pathVCSPackage:
		p.loadVCSPackage(currentPath)
	}
}

// loadVCSPackage loads the package found at the given VCS path.
func (p *PackageLoader) loadVCSPackage(packagePath pathInformation) {
	rootDirectory := path.Dir(p.rootSourceFile)
	pkgDirectory := path.Join(rootDirectory, serulianPackageDirectory)

	// Perform the checkout of the VCS package.
	checkoutDirectory, err, warning := vcs.PerformVCSCheckout(packagePath.path, pkgDirectory)
	if err != nil {
		p.errors <- fmt.Errorf("Error loading VCS package '%s' referenced by '%s': %v", packagePath.path, packagePath.source, err)
		return
	}

	if warning != "" {
		p.warnings <- warning
	}

	// Push the now-local directory onto the package loading channel.
	p.pushPath(pathLocalPackage, checkoutDirectory, packagePath.source)
}

// loadLocalPackage loads the package found at the path relative to the package directory.
func (p *PackageLoader) loadLocalPackage(packagePath pathInformation) {
	// Ensure the directory exists.
	if ok, _ := exists(packagePath.path); !ok {
		p.errors <- fmt.Errorf("Path '%s' referenced by '%s' could not be found", packagePath.path, packagePath.source)
		return
	}

	// Read the contents of the directory.
	directoryContents, err := ioutil.ReadDir(packagePath.path)
	if err != nil {
		p.errors <- err
		return
	}

	// Find all source files in the directory and add them to the paths list.
	var fileFound bool

	for _, fileInfo := range directoryContents {
		if path.Ext(fileInfo.Name()) == serulianSourceFileExtension {
			fileFound = true
			filePath := path.Join(packagePath.path, fileInfo.Name())
			p.pushPath(pathSourceFile, filePath, packagePath.path)
		}
	}

	if !fileFound {
		p.warnings <- fmt.Sprintf("Package '%s' has no source files", packagePath.path)
		return
	}
}

// conductParsing performs parsing of a source file found at the given path.
func (p *PackageLoader) conductParsing(sourceFile pathInformation) {
	// Ensure the file exists.
	if ok, _ := exists(sourceFile.path); !ok {
		p.errors <- fmt.Errorf("Path '%s' referenced by '%s' could not be found", sourceFile.path, sourceFile.source)
		return
	}

	// Load the source file's contents.
	contents, err := ioutil.ReadFile(sourceFile.path)
	if err != nil {
		p.errors <- fmt.Errorf("Path '%s' referenced by '%s' could not be loaded: %v", sourceFile.path, sourceFile.source, err)
		return
	}

	// Parse the source file.
	parser.Parse(p.nodeBuilder, p.handleImport, parser.InputSource(sourceFile.path), string(contents))
}

// handleImport queues an import found in a source file.
func (p *PackageLoader) handleImport(importInformation parser.PackageImport) {
	sourcePath := string(importInformation.SourceFile)

	if importInformation.ImportType == parser.ImportTypeVCS {
		// VCS paths get added directly.
		p.pushPath(pathVCSPackage, importInformation.Path, sourcePath)
	} else {
		// Otherwise, check the path to see if it exists as a single source file. If so, we add it
		// as a source file instead of a local package.
		dirPath := path.Join(path.Dir(sourcePath), importInformation.Path)
		filePath := dirPath + serulianSourceFileExtension
		if ok, _ := exists(filePath); ok {
			p.pushPath(pathSourceFile, filePath, sourcePath)
		} else {
			p.pushPath(pathLocalPackage, dirPath, sourcePath)
		}
	}
}
