// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// packageloader package defines functions and types for loading and parsing source from disk or VCS.
package packageloader

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"sync"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/vcs"

	"github.com/streamrail/concurrent-map"
)

// PackageInfo holds information about a loaded package.
type PackageInfo struct {
	kind        string                       // The kind of the package.
	referenceId string                       // The unique ID for this package.
	modulePaths []compilercommon.InputSource // The module paths making up this package.
}

// Kind returns the kind code this package.
func (pi PackageInfo) Kind() string {
	return pi.kind
}

// ReferenceId returns the unique reference ID for this package.
func (pi PackageInfo) ReferenceId() string {
	return pi.referenceId
}

// ModulePaths returns the list of full paths of the modules in this package.
func (pi PackageInfo) ModulePaths() []compilercommon.InputSource {
	return pi.modulePaths
}

// pathKind identifies the supported kind of paths
type pathKind int

const (
	pathSourceFile pathKind = iota
	pathLocalPackage
	pathVCSPackage
)

// serulianPackageDirectory is the directory under the root directory holding cached packages.
const serulianPackageDirectory = ".pkg"

// pathInformation holds information about a path to load.
type pathInformation struct {
	referenceId string
	kind        pathKind
	path        string
	sourceKind  string
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

	errors   chan compilercommon.SourceError   // Errors are reported on this channel
	warnings chan compilercommon.SourceWarning // Warnings are reported on this channel

	handlers map[string]SourceHandler // The handlers for each of the supported package kinds.

	pathsToLoad      chan pathInformation // The paths to load
	pathsEncountered cmap.ConcurrentMap   // The paths processed by the loader goroutine
	packageMap       cmap.ConcurrentMap   // The package map

	workTracker sync.WaitGroup // WaitGroup used to wait until all loading is complete
	finished    chan bool      // Channel used to tell background goroutines to quit
}

// Library contains a reference to an external library to load, in addition to those referenced
// by the root source file.
type Library struct {
	PathOrURL string // The file location or SCM URL of the library's package.
	IsSCM     bool   // If true, the PathOrURL is treated as a remote SCM package.
	Kind      string // The kind of the library. Leave empty for Serulian files.
}

// LoadResult contains the result of attempting to load all packages and source files for this
// project.
type LoadResult struct {
	Status     bool                           // True on success, false otherwise
	Errors     []compilercommon.SourceError   // The errors encountered, if any
	Warnings   []compilercommon.SourceWarning // The warnings encountered, if any
	PackageMap map[string]PackageInfo         // Map of packages loaded
}

// NewPackageLoader creates and returns a new package loader for the given path.
func NewPackageLoader(rootSourceFile string, handlers ...SourceHandler) *PackageLoader {
	handlersMap := map[string]SourceHandler{}
	for _, handler := range handlers {
		handlersMap[handler.Kind()] = handler
	}

	return &PackageLoader{
		rootSourceFile: rootSourceFile,

		errors:   make(chan compilercommon.SourceError),
		warnings: make(chan compilercommon.SourceWarning),

		handlers: handlersMap,

		pathsEncountered: cmap.New(),
		packageMap:       cmap.New(),
		pathsToLoad:      make(chan pathInformation),

		finished: make(chan bool, 2),
	}
}

// Load performs the loading of a Serulian package found at the directory path.
// Any libraries specified will be loaded as well.
func (p *PackageLoader) Load(libraries ...Library) LoadResult {
	// Start the loading goroutine.
	go p.loadAndParse()

	// Start the error/warning collection goroutine.
	result := &LoadResult{
		Status:   true,
		Errors:   make([]compilercommon.SourceError, 0),
		Warnings: make([]compilercommon.SourceWarning, 0),
	}

	go p.collectIssues(result)

	// Add the root source file as the first package to be parsed.
	sal := compilercommon.NewSourceAndLocation(compilercommon.InputSource(p.rootSourceFile), 0)

	var added = false
	for _, handler := range p.handlers {
		if strings.HasSuffix(p.rootSourceFile, handler.PackageFileExtension()) {
			p.pushPath(pathSourceFile, handler.Kind(), p.rootSourceFile, sal)
			added = true
			break
		}
	}

	if !added {
		log.Fatalf("Could not find handler for root source file: %v", p.rootSourceFile)
	}

	// Add the libraries to be parsed.
	for _, library := range libraries {
		if library.IsSCM {
			sal := compilercommon.NewSourceAndLocation(compilercommon.InputSource(library.PathOrURL), 0)
			p.pushPath(pathVCSPackage, library.Kind, library.PathOrURL, sal)
		} else {
			sal := compilercommon.NewSourceAndLocation(compilercommon.InputSource(library.PathOrURL), 0)
			p.pushPath(pathLocalPackage, library.Kind, library.PathOrURL, sal)
		}
	}

	// Wait for all packages and source files to be completed.
	p.workTracker.Wait()

	// Tell the goroutines to quit.
	p.finished <- true
	p.finished <- true

	// Save the package map.
	packageMap := map[string]PackageInfo{}
	for entry := range p.packageMap.Iter() {
		packageMap[entry.Key] = entry.Val.(PackageInfo)
	}

	result.PackageMap = packageMap

	// Apply all handler changes.
	for _, handler := range p.handlers {
		handler.Apply(packageMap)
	}

	// Perform verification in all handlers.
	if len(result.Errors) == 0 {
		errorReporter := func(err compilercommon.SourceError) {
			result.Errors = append(result.Errors, err)
			result.Status = false
		}

		warningReporter := func(warning compilercommon.SourceWarning) {
			result.Warnings = append(result.Warnings, warning)
		}

		for _, handler := range p.handlers {
			handler.Verify(errorReporter, warningReporter)
		}
	}

	return *result
}

// pushPath adds a path to be processed by the package loader.
func (p *PackageLoader) pushPath(kind pathKind, sourceKind string, path string, sal compilercommon.SourceAndLocation) string {
	return p.pushPathWithId(path, sourceKind, kind, path, sal)
}

// pushPathWithId adds a path to be processed by the package loader, with the specified ID.
func (p *PackageLoader) pushPathWithId(pathId string, sourceKind string, kind pathKind, path string, sal compilercommon.SourceAndLocation) string {
	p.workTracker.Add(1)
	p.pathsToLoad <- pathInformation{pathId, kind, path, sourceKind, sal}
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
	if !p.pathsEncountered.SetIfAbsent(pathKey, true) {
		return
	}

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
	p.pushPathWithId(packagePath.referenceId, packagePath.sourceKind, pathLocalPackage, checkoutDirectory, packagePath.sal)
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
		kind:        packagePath.sourceKind,
		referenceId: packagePath.referenceId,
		modulePaths: make([]compilercommon.InputSource, 0),
	}

	// Load the handler for the package.
	handler, hasHandler := p.handlers[packagePath.sourceKind]
	if !hasHandler {
		log.Fatalf("Missing handler for source file of kind: [%v]", packagePath.sourceKind)
	}

	// Find all source files in the directory and add them to the paths list.
	var fileFound bool
	for _, fileInfo := range directoryContents {
		if path.Ext(fileInfo.Name()) == handler.PackageFileExtension() {
			fileFound = true
			filePath := path.Join(packagePath.path, fileInfo.Name())
			p.pushPath(pathSourceFile, packagePath.sourceKind, filePath, packagePath.sal)

			// Add the source file to the package information.
			packageInfo.modulePaths = append(packageInfo.modulePaths, compilercommon.InputSource(filePath))
		}
	}

	p.packageMap.Set(packagePath.referenceId, *packageInfo)
	if !fileFound {
		p.warnings <- compilercommon.SourceWarningf(packagePath.sal, "Package '%s' has no source files", packagePath.path)
		return
	}
}

// conductParsing performs parsing of a source file found at the given path.
func (p *PackageLoader) conductParsing(sourceFile pathInformation) {
	inputSource := compilercommon.InputSource(sourceFile.path)

	// Add the file to the package map as a package of one file.
	p.packageMap.Set(sourceFile.referenceId, PackageInfo{
		kind:        sourceFile.sourceKind,
		referenceId: sourceFile.referenceId,
		modulePaths: []compilercommon.InputSource{inputSource},
	})

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
	handler, hasHandler := p.handlers[sourceFile.sourceKind]
	if !hasHandler {
		log.Fatalf("Missing handler for source file of kind: [%v]", sourceFile.sourceKind)
	}

	handler.Parse(inputSource, string(contents), p.handleImport)
}

// handleImport queues an import found in a source file.
func (p *PackageLoader) handleImport(importInformation PackageImport) string {
	handler, hasHandler := p.handlers[importInformation.Kind]
	if !hasHandler {
		p.errors <- compilercommon.SourceErrorf(importInformation.SourceLocation, "Unknown kind of import '%s'", importInformation.Kind)
		return ""
	}

	sourcePath := string(importInformation.SourceLocation.Source())

	if importInformation.ImportType == ImportTypeVCS {
		// VCS paths get added directly.
		return p.pushPath(pathVCSPackage, importInformation.Kind, importInformation.Path, importInformation.SourceLocation)
	} else {
		// Otherwise, check the path to see if it exists as a single source file. If so, we add it
		// as a source file instead of a local package.
		dirPath := path.Join(path.Dir(sourcePath), importInformation.Path)
		filePath := dirPath + handler.PackageFileExtension()

		if ok, _ := exists(filePath); ok {
			return p.pushPath(pathSourceFile, handler.Kind(), filePath, importInformation.SourceLocation)
		} else {
			return p.pushPath(pathLocalPackage, handler.Kind(), dirPath, importInformation.SourceLocation)
		}
	}
}
