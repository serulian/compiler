// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package packageloader defines functions and types for loading and parsing source from disk or VCS.
package packageloader

import (
	"fmt"
	"log"
	"path"
	"strings"
	"sync"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/vcs"

	cmap "github.com/streamrail/concurrent-map"
)

// SerulianPackageDirectory is the directory under the root directory holding cached packages.
const SerulianPackageDirectory = ".pkg"

// PackageInfo holds information about a loaded package.
type PackageInfo struct {
	kind        string                       // The source kind of the package.
	referenceId string                       // The unique ID for this package.
	modulePaths []compilercommon.InputSource // The module paths making up this package.
}

// PackageInfoForTesting returns a PackageInfo for testing.
func PackageInfoForTesting(kind string, modulePaths []compilercommon.InputSource) PackageInfo {
	return PackageInfo{kind, "testpackage", modulePaths}
}

// Kind returns the source kind of this package.
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

// pathInformation holds information about a path to load.
type pathInformation struct {
	// referenceId is the unique reference ID for the path. Typically the path
	// of the package or module itself, but can be anything, as long as it uniquely
	// identifies this path.
	referenceId string

	// kind is the kind of path (source file, local package or vcs package).
	kind pathKind

	// path is the path being represented.
	path string

	// sourceKind is the source kind (empty string for `.seru`, `webidl` for "webidl")
	sourceKind string

	// sal is the source and location that *referenced* this path.
	sal compilercommon.SourceAndLocation
}

// Returns the string representation of the given path.
func (p *pathInformation) String() string {
	return fmt.Sprintf("%v::%s::%s", int(p.kind), p.sourceKind, p.path)
}

// PackageLoader helps to fully and recursively load a Serulian package and its dependencies
// from a directory or set of directories.
type PackageLoader struct {
	rootSourceFile            string     // The root source file location.
	vcsDevelopmentDirectories []string   // Directories to check for VCS packages before VCS checkout.
	pathLoader                PathLoader // The path loaders to use.
	alwaysValidate            bool       // Whether to always run validation, regardless of errors. Useful to IDE tooling.
	skipVCSRefresh            bool       // Whether to skip VCS refresh if cache exists. Useful to IDE tooling.

	errors   chan compilercommon.SourceError   // Errors are reported on this channel
	warnings chan compilercommon.SourceWarning // Warnings are reported on this channel

	handlers map[string]SourceHandler // The handlers for each of the supported package kinds.

	pathsToLoad          chan pathInformation // The paths to load
	pathKindsEncountered cmap.ConcurrentMap   // The path+kinds processed by the loader goroutine
	vcsPathsLoaded       cmap.ConcurrentMap   // The VCS paths that have been loaded, mapping to their checkout dir
	vcsLockMap           lockMap              // LockMap for ensuring single loads of all VCS paths.
	packageMap           *mutablePackageMap   // The package map

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
	PackageMap LoadedPackageMap               // Map of packages loaded
}

// Config defines configuration for a PackageLoader.
type Config struct {
	// The path of the root source file from which loading will begin.
	RootSourceFilePath string

	// Paths of directories to check for local copies of remote packages
	// before performing VCS checkout.
	VCSDevelopmentDirectories []string

	// The source handlers to use to parse and import the loaded source files.
	SourceHandlers []SourceHandler

	// The path loader to use to load source files and directories.
	PathLoader PathLoader

	// If true, validation will always be called.
	AlwaysValidate bool

	// If true, VCS forced refresh (for branches and HEAD) will be skipped if cache
	// exists.
	SkipVCSRefresh bool
}

// NewBasicConfig returns PackageLoader Config for a root source file and source handlers.
func NewBasicConfig(rootSourceFilePath string, sourceHandlers ...SourceHandler) Config {
	return Config{
		RootSourceFilePath:        rootSourceFilePath,
		VCSDevelopmentDirectories: []string{},
		SourceHandlers:            sourceHandlers,
		PathLoader:                LocalFilePathLoader{},
		AlwaysValidate:            false,
		SkipVCSRefresh:            false,
	}
}

// NewPackageLoader creates and returns a new package loader for the given config.
func NewPackageLoader(config Config) *PackageLoader {
	handlersMap := map[string]SourceHandler{}
	for _, handler := range config.SourceHandlers {
		handlersMap[handler.Kind()] = handler
	}

	pathLoader := config.PathLoader
	if pathLoader == nil {
		pathLoader = LocalFilePathLoader{}
	}

	return &PackageLoader{
		rootSourceFile:            config.RootSourceFilePath,
		vcsDevelopmentDirectories: config.VCSDevelopmentDirectories,
		pathLoader:                pathLoader,
		alwaysValidate:            config.AlwaysValidate,
		skipVCSRefresh:            config.SkipVCSRefresh,

		errors:   make(chan compilercommon.SourceError),
		warnings: make(chan compilercommon.SourceWarning),

		handlers: handlersMap,

		pathKindsEncountered: cmap.New(),
		packageMap:           newMutablePackageMap(),
		pathsToLoad:          make(chan pathInformation),

		vcsPathsLoaded: cmap.New(),
		vcsLockMap:     createLockMap(),

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
	result.PackageMap = p.packageMap.Build()

	// Apply all handler changes.
	for _, handler := range p.handlers {
		handler.Apply(result.PackageMap)
	}

	// Perform verification in all handlers.
	if p.alwaysValidate || len(result.Errors) == 0 {
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

// PathLoader returns the path loader used by this package manager.
func (p *PackageLoader) PathLoader() PathLoader {
	return p.pathLoader
}

// ModuleOrPackage defines a reference to a module or package.
type ModuleOrPackage struct {
	// Name is the name of the module or package.
	Name string

	// Path is the on-disk path of the module or package.
	Path string

	// SourceKind is the kind source for the module or package. Packages will always be
	// empty.
	SourceKind string
}

// ListSubModulesAndPackages lists all modules or packages found *directly* under the given path.
func (p *PackageLoader) ListSubModulesAndPackages(packagePath string) ([]ModuleOrPackage, error) {
	directoryContents, err := p.pathLoader.LoadDirectory(packagePath)
	if err != nil {
		return []ModuleOrPackage{}, err
	}

	var modulesOrPackages = make([]ModuleOrPackage, 0, len(directoryContents))
	for _, entry := range directoryContents {
		if entry.IsDirectory {
			modulesOrPackages = append(modulesOrPackages, ModuleOrPackage{entry.Name, path.Join(packagePath, entry.Name), ""})
			continue
		}

		for _, handler := range p.handlers {
			if strings.HasSuffix(entry.Name, handler.PackageFileExtension()) {
				name := entry.Name[0 : len(entry.Name)-len(handler.PackageFileExtension())]
				modulesOrPackages = append(modulesOrPackages, ModuleOrPackage{name, path.Join(packagePath, entry.Name), handler.Kind()})
				break
			}
		}
	}

	return modulesOrPackages, nil
}

// LocalPackageInfoForPath returns the package information for the given path. Note that VCS paths will
// be converted into their local package equivalent. If the path refers to a source file instead of a
// directory, a package containing the single module will be returned.
func (p *PackageLoader) LocalPackageInfoForPath(path string, sourceKind string, isVCSPath bool) (PackageInfo, error) {
	if isVCSPath {
		localPath, err := p.getVCSDirectoryForPath(path)
		if err != nil {
			return PackageInfo{}, err
		}

		path = localPath
	}

	// Find the source handler matching the source kind.
	handler, ok := p.handlers[sourceKind]
	if !ok {
		return PackageInfo{}, fmt.Errorf("Unknown source kind %s", sourceKind)
	}

	// Check for a single module.
	filePath := path + handler.PackageFileExtension()
	if p.pathLoader.IsSourceFile(filePath) {
		return PackageInfo{
			kind:        sourceKind,
			referenceId: filePath,
			modulePaths: []compilercommon.InputSource{compilercommon.InputSource(filePath)},
		}, nil
	}

	// Otherwise, read the contents of the directory.
	return p.packageInfoForPackageDirectory(path, sourceKind)
}

// packageInfoForDirectory returns a PackageInfo for the package found at the given path.
func (p *PackageLoader) packageInfoForPackageDirectory(packagePath string, sourceKind string) (PackageInfo, error) {
	directoryContents, err := p.pathLoader.LoadDirectory(packagePath)
	if err != nil {
		return PackageInfo{}, err
	}

	handler, ok := p.handlers[sourceKind]
	if !ok {
		return PackageInfo{}, fmt.Errorf("Unknown source kind %s", sourceKind)
	}

	packageInfo := &PackageInfo{
		kind:        sourceKind,
		referenceId: packagePath,
		modulePaths: make([]compilercommon.InputSource, 0),
	}

	// Find all source files in the directory and add them to the paths list.
	for _, entry := range directoryContents {
		if !entry.IsDirectory && path.Ext(entry.Name) == handler.PackageFileExtension() {
			filePath := path.Join(packagePath, entry.Name)

			// Add the source file to the package information.
			packageInfo.modulePaths = append(packageInfo.modulePaths, compilercommon.InputSource(filePath))
		}
	}

	return *packageInfo, nil
}

// getVCSDirectoryForPath returns the directory on disk where the given VCS path will be placed, if any.
func (p *PackageLoader) getVCSDirectoryForPath(vcsPath string) (string, error) {
	rootDirectory := path.Dir(p.rootSourceFile)
	pkgDirectory := path.Join(rootDirectory, SerulianPackageDirectory)
	return vcs.GetVCSCheckoutDirectory(vcsPath, pkgDirectory, p.vcsDevelopmentDirectories...)
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

	// Ensure we have not already seen this path and kind.
	pathKey := currentPath.String()
	if !p.pathKindsEncountered.SetIfAbsent(pathKey, true) {
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
	// Lock on the package path to ensure no other checkouts occur for this path.
	pathLock := p.vcsLockMap.getLock(packagePath.path)
	pathLock.Lock()
	defer pathLock.Unlock()

	existingCheckoutDir, exists := p.vcsPathsLoaded.Get(packagePath.path)
	if exists {
		// Note: existingCheckoutDir will be empty if there was an error loading the VCS.
		if existingCheckoutDir != "" {
			// Push the now-local directory onto the package loading channel.
			p.pushPathWithId(packagePath.referenceId, packagePath.sourceKind, pathLocalPackage, existingCheckoutDir.(string), packagePath.sal)
			return
		}
	}

	rootDirectory := path.Dir(p.rootSourceFile)
	pkgDirectory := path.Join(rootDirectory, SerulianPackageDirectory)

	// Perform the checkout of the VCS package.
	var cacheOption = vcs.VCSFollowNormalCacheRules
	if p.skipVCSRefresh {
		cacheOption = vcs.VCSAlwaysUseCache
	}

	checkoutDirectory, err, warning := vcs.PerformVCSCheckout(packagePath.path, pkgDirectory, cacheOption, p.vcsDevelopmentDirectories...)
	if err != nil {
		p.vcsPathsLoaded.Set(packagePath.path, "")
		p.errors <- compilercommon.SourceErrorf(packagePath.sal, "Error loading VCS package '%s': %v", packagePath.path, err)
		return
	}

	p.vcsPathsLoaded.Set(packagePath.path, checkoutDirectory)
	if warning != "" {
		p.warnings <- compilercommon.NewSourceWarning(packagePath.sal, warning)
	}

	// Push the now-local directory onto the package loading channel.
	p.pushPathWithId(packagePath.referenceId, packagePath.sourceKind, pathLocalPackage, checkoutDirectory, packagePath.sal)
}

// loadLocalPackage loads the package found at the path relative to the package directory.
func (p *PackageLoader) loadLocalPackage(packagePath pathInformation) {
	packageInfo, err := p.packageInfoForPackageDirectory(packagePath.path, packagePath.sourceKind)
	if err != nil {
		p.errors <- compilercommon.SourceErrorf(packagePath.sal, "Could not load directory '%s'", packagePath.path)
		return
	}

	// Add the module paths to be parsed.
	var moduleFound = false
	for _, modulePath := range packageInfo.ModulePaths() {
		p.pushPath(pathSourceFile, packagePath.sourceKind, string(modulePath), packagePath.sal)
		moduleFound = true
	}

	// Add the package itself to the package map.
	p.packageMap.Add(packagePath.sourceKind, packagePath.referenceId, packageInfo)
	if !moduleFound {
		p.warnings <- compilercommon.SourceWarningf(packagePath.sal, "Package '%s' has no source files", packagePath.path)
		return
	}
}

// conductParsing performs parsing of a source file found at the given path.
func (p *PackageLoader) conductParsing(sourceFile pathInformation) {
	inputSource := compilercommon.InputSource(sourceFile.path)

	// Add the file to the package map as a package of one file.
	p.packageMap.Add(sourceFile.sourceKind, sourceFile.referenceId, PackageInfo{
		kind:        sourceFile.sourceKind,
		referenceId: sourceFile.referenceId,
		modulePaths: []compilercommon.InputSource{inputSource},
	})

	// Load the source file's contents.
	contents, err := p.pathLoader.LoadSourceFile(sourceFile.path)
	if err != nil {
		p.errors <- compilercommon.SourceErrorf(sourceFile.sal, "Could not load source file '%s': %v", sourceFile.path, err)
		return
	}

	// Parse the source file.
	handler, hasHandler := p.handlers[sourceFile.sourceKind]
	if !hasHandler {
		log.Fatalf("Missing handler for source file of kind: [%v]", sourceFile.sourceKind)
	}

	handler.Parse(inputSource, string(contents), p.handleImport)
}

// verifyNoVCSBoundaryCross does a check to ensure that walking from the given start path
// to the given end path does not cross a VCS boundary. If it does, an error is returned.
func (p *PackageLoader) verifyNoVCSBoundaryCross(startPath string, endPath string, title string, importInformation PackageImport) *compilercommon.SourceError {
	var checkPath = startPath
	for {
		if checkPath == endPath {
			return nil
		}

		if vcs.IsVCSRootDirectory(checkPath) {
			err := compilercommon.SourceErrorf(importInformation.SourceLocation,
				"Import of %s '%s' crosses VCS boundary at package '%s'", title,
				importInformation.Path, checkPath)
			return &err
		}

		nextPath := path.Dir(checkPath)
		if checkPath == nextPath {
			return nil
		}

		checkPath = nextPath
	}
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
	}

	// Check the path to see if it exists as a single source file. If so, we add it
	// as a source file instead of a local package.
	currentDirectory := path.Dir(sourcePath)
	dirPath := path.Join(currentDirectory, importInformation.Path)
	filePath := dirPath + handler.PackageFileExtension()

	var importedDirectoryPath = dirPath
	var title = "package"

	// Determine if path refers to a single source file. If so, it is imported rather than
	// the entire directory.
	isSourceFile := p.pathLoader.IsSourceFile(filePath)
	if isSourceFile {
		title = "module"
		importedDirectoryPath = path.Dir(filePath)
	}

	// Check to ensure we are not crossing a VCS boundary.
	if currentDirectory != importedDirectoryPath {
		// If the imported directory is underneath the current directory, we need to walk upward.
		if strings.HasPrefix(importedDirectoryPath, currentDirectory) {
			err := p.verifyNoVCSBoundaryCross(importedDirectoryPath, currentDirectory, title, importInformation)
			if err != nil {
				p.errors <- *err
				return ""
			}
		} else {
			// Otherwise, we walk upward from the current directory to the imported directory.
			err := p.verifyNoVCSBoundaryCross(currentDirectory, importedDirectoryPath, title, importInformation)
			if err != nil {
				p.errors <- *err
				return ""
			}
		}
	}

	// Push the imported path.
	if isSourceFile {
		return p.pushPath(pathSourceFile, handler.Kind(), filePath, importInformation.SourceLocation)
	}

	return p.pushPath(pathLocalPackage, handler.Kind(), dirPath, importInformation.SourceLocation)
}
