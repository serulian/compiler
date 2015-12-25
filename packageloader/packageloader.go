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

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/parser"
	"github.com/serulian/compiler/vcs"
)

// PackageInfo holds information about a loaded package.
type PackageInfo struct {
	referenceId string                       // The unique ID for this package.
	modulePaths []compilercommon.InputSource // The module paths making up this package.
}

// ReferenceId returns the unique reference ID for this package.
func (pi *PackageInfo) ReferenceId() string {
	return pi.referenceId
}

// ModulePaths returns the list of full paths of the modules in this package.
func (pi *PackageInfo) ModulePaths() []compilercommon.InputSource {
	return pi.modulePaths
}

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
	referenceId string
	kind        pathKind
	path        string
	sal         compilercommon.SourceAndLocation
}

// Returns the string representation of the given path.
func (p *pathInformation) String() string {
	return fmt.Sprintf("%s://%s", p.kind, p.path)
}

// PackageLoader helps to fully and recursively load a Serulian package and its dependencies
// from a directory or set of directories.
type PackageLoader struct {
	rootSourceFile string // The root source file location.

	errors   chan *compilercommon.SourceError   // Errors are reported on this channel
	warnings chan *compilercommon.SourceWarning // Warnings are reported on this channel

	nodeBuilder parser.NodeBuilder // Builder to use for constructing the parse nodes

	pathsToLoad      chan pathInformation    // The paths to load
	pathsEncountered map[string]bool         // The paths processed by the loader goroutine
	packageMap       map[string]*PackageInfo // The package map

	workTracker sync.WaitGroup // WaitGroup used to wait until all loading is complete
	finished    chan bool      // Channel used to tell background goroutines to quit
}

// Library contains a reference to an external library to load, in addition to those referenced
// by the root source file.
type Library struct {
	PathOrURL string // The file location or SCM URL of the library's package.
	IsSCM     bool   // If true, the PathOrURL is treated as a remote SCM package.
}

// LoadResult contains the result of attempting to load all packages and source files for this
// project.
type LoadResult struct {
	Status     bool                            // True on success, false otherwise
	Errors     []*compilercommon.SourceError   // The errors encountered, if any
	Warnings   []*compilercommon.SourceWarning // The warnings encountered, if any
	PackageMap map[string]*PackageInfo         // Map of packages loaded
}

// NewPackageLoader creates and returns a new package loader for the given path.
func NewPackageLoader(rootSourceFile string, nodeBuilder parser.NodeBuilder) *PackageLoader {
	return &PackageLoader{
		rootSourceFile: rootSourceFile,

		errors:   make(chan *compilercommon.SourceError),
		warnings: make(chan *compilercommon.SourceWarning),

		nodeBuilder: nodeBuilder,

		pathsEncountered: map[string]bool{},
		packageMap:       map[string]*PackageInfo{},
		pathsToLoad:      make(chan pathInformation),

		finished: make(chan bool, 2),
	}
}

// Load performs the loading of a Serulian package found at the directory path.
// Any libraries specified will be loaded as well.
func (p *PackageLoader) Load(libraries ...Library) *LoadResult {
	// Start the loading goroutine.
	go p.loadAndParse()

	// Start the error/warning collection goroutine.
	result := &LoadResult{
		Status:   true,
		Errors:   make([]*compilercommon.SourceError, 0),
		Warnings: make([]*compilercommon.SourceWarning, 0),
	}

	go p.collectIssues(result)

	// Add the root source file as the first package to be parsed.
	sal := compilercommon.NewSourceAndLocation(compilercommon.InputSource(p.rootSourceFile), 0)
	p.pushPath(pathSourceFile, p.rootSourceFile, sal)

	// Add the libraries to be parsed.
	for _, library := range libraries {
		if library.IsSCM {
			sal := compilercommon.NewSourceAndLocation(compilercommon.InputSource(library.PathOrURL), 0)
			p.pushPath(pathVCSPackage, library.PathOrURL, sal)
		} else {
			sal := compilercommon.NewSourceAndLocation(compilercommon.InputSource(library.PathOrURL), 0)
			p.pushPath(pathLocalPackage, library.PathOrURL, sal)
		}
	}

	// Wait for all packages and source files to be completed.
	p.workTracker.Wait()

	// Tell the goroutines to quit.
	p.finished <- true
	p.finished <- true

	// Save the package map.
	result.PackageMap = p.packageMap
	return result
}

// pushPath adds a path to be processed by the package loader.
func (p *PackageLoader) pushPath(kind pathKind, path string, sal compilercommon.SourceAndLocation) string {
	return p.pushPathWithId(path, kind, path, sal)
}

// pushPathWithId adds a path to be processed by the package loader, with the specified ID.
func (p *PackageLoader) pushPathWithId(pathId string, kind pathKind, path string, sal compilercommon.SourceAndLocation) string {
	p.workTracker.Add(1)
	p.pathsToLoad <- pathInformation{pathId, kind, path, sal}
	return pathId
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
		p.errors <- compilercommon.SourceErrorf(packagePath.sal, "Error loading VCS package '%s'", packagePath.path)
		return
	}

	if warning != "" {
		p.warnings <- compilercommon.NewSourceWarning(packagePath.sal, warning)
	}

	// Push the now-local directory onto the package loading channel.
	p.pushPathWithId(packagePath.referenceId, pathLocalPackage, checkoutDirectory, packagePath.sal)
}

// loadLocalPackage loads the package found at the path relative to the package directory.
func (p *PackageLoader) loadLocalPackage(packagePath pathInformation) {
	// Ensure the directory exists.
	if ok, _ := exists(packagePath.path); !ok {
		p.errors <- compilercommon.SourceErrorf(packagePath.sal, "Could not find directory '%s'", packagePath.path)
		return
	}

	// Read the contents of the directory.
	directoryContents, err := ioutil.ReadDir(packagePath.path)
	if err != nil {
		p.errors <- compilercommon.SourceErrorf(packagePath.sal, "Could not load directory '%s'", packagePath.path)
		return
	}

	// Add the package information to the map.
	packageInfo := &PackageInfo{
		referenceId: packagePath.referenceId,
		modulePaths: make([]compilercommon.InputSource, 0),
	}

	// Find all source files in the directory and add them to the paths list.
	var fileFound bool

	for _, fileInfo := range directoryContents {
		if path.Ext(fileInfo.Name()) == serulianSourceFileExtension {
			fileFound = true
			filePath := path.Join(packagePath.path, fileInfo.Name())
			p.pushPath(pathSourceFile, filePath, packagePath.sal)

			// Add the source file to the package information.
			packageInfo.modulePaths = append(packageInfo.modulePaths, compilercommon.InputSource(filePath))
		}
	}

	p.packageMap[packagePath.referenceId] = packageInfo

	if !fileFound {
		p.warnings <- compilercommon.SourceWarningf(packagePath.sal, "Package '%s' has no source files", packagePath.path)
		return
	}
}

// conductParsing performs parsing of a source file found at the given path.
func (p *PackageLoader) conductParsing(sourceFile pathInformation) {
	inputSource := compilercommon.InputSource(sourceFile.path)

	// Add the file to the package map as a package of one file.
	p.packageMap[sourceFile.referenceId] = &PackageInfo{
		referenceId: sourceFile.referenceId,
		modulePaths: []compilercommon.InputSource{inputSource},
	}

	// Ensure the file exists.
	if ok, _ := exists(sourceFile.path); !ok {
		p.errors <- compilercommon.SourceErrorf(sourceFile.sal, "Could not find source file '%s'", sourceFile.path)
		return
	}

	// Load the source file's contents.
	contents, err := ioutil.ReadFile(sourceFile.path)
	if err != nil {
		p.errors <- compilercommon.SourceErrorf(sourceFile.sal, "Could not load source file '%s'", sourceFile.path)
		return
	}

	// Parse the source file.
	parser.Parse(p.nodeBuilder, p.handleImport, inputSource, string(contents))
}

// handleImport queues an import found in a source file.
func (p *PackageLoader) handleImport(importInformation parser.PackageImport) string {
	sourcePath := string(importInformation.SourceLocation.Source())

	if importInformation.ImportType == parser.ImportTypeVCS {
		// VCS paths get added directly.
		return p.pushPath(pathVCSPackage, importInformation.Path, importInformation.SourceLocation)
	} else {
		// Otherwise, check the path to see if it exists as a single source file. If so, we add it
		// as a source file instead of a local package.
		dirPath := path.Join(path.Dir(sourcePath), importInformation.Path)
		filePath := dirPath + serulianSourceFileExtension
		if ok, _ := exists(filePath); ok {
			return p.pushPath(pathSourceFile, filePath, importInformation.SourceLocation)
		} else {
			return p.pushPath(pathLocalPackage, dirPath, importInformation.SourceLocation)
		}
	}
}
