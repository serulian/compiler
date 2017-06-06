// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Some code based on the Golang VCS portion of the cmd package:
// https://golang.org/src/cmd/go/vcs.go

// vcs package defines helpers functions and interfaces for working with Version Control Systems
// such as git, including discovery of VCS information based on the Golang VCS discovery protocol.
//
// The discovery protocol is documented loosely here: https://golang.org/cmd/go/#hdr-Remote_import_paths
package vcs

import (
	"fmt"
	"log"
	"os"
	"path"
)

// IsVCSRootDirectory returns true if the given local file system path is a VCS root directory.
// Note that this method will return false if the path does not exist locally.
func IsVCSRootDirectory(localPath string) bool {
	_, hasHandler := detectHandler(localPath)
	return hasHandler
}

// GetVCSCheckoutDirectory returns the path of the directory into which the given VCS path will checked out,
// if PerformVCSCheckout is called.
func GetVCSCheckoutDirectory(vcsPath string, pkgCacheRootPath string, vcsDevelopmentDirectories ...string) (string, error) {
	// Parse the VCS path.
	parsedPath, perr := ParseVCSPath(vcsPath)
	if perr != nil {
		return "", perr
	}

	// Check for a local development cache.
	if parsedPath.isRepoOnlyReference() {
		for _, directoryPath := range vcsDevelopmentDirectories {
			fullCheckDirectory := path.Join(directoryPath, parsedPath.url)
			if _, err := os.Stat(fullCheckDirectory); err == nil {
				return fullCheckDirectory, nil
			}
		}
	}

	fullCacheDirectory := path.Join(pkgCacheRootPath, parsedPath.cacheDirectory())
	return fullCacheDirectory, nil
}

// VCSCacheOption defines the caching options for VCS checkout.
type VCSCacheOption int

const (
	// VCSFollowNormalCacheRules indicates that VCS checkouts will be pulled from cache unless a HEAD
	// reference.
	VCSFollowNormalCacheRules VCSCacheOption = iota

	// VCSAlwaysUseCache indicates that VCS checkouts will always use the cache if available.
	VCSAlwaysUseCache
)

// PerformVCSCheckout performs the checkout and updating of the given VCS path and returns
// the local system directory at which the package was checked out.
//
// pkgCacheRootPath holds the path of the root directory that forms the package cache.
//
// vcsDevelopmentDirectories specifies optional directories to check for branchless and tagless copies
// of the repository first. If found, the copy will be used in lieu of a normal checkout.
func PerformVCSCheckout(vcsPath string, pkgCacheRootPath string, cacheOption VCSCacheOption, vcsDevelopmentDirectories ...string) (string, error, string) {
	// Parse the VCS path.
	parsedPath, perr := ParseVCSPath(vcsPath)
	if perr != nil {
		return "", perr, ""
	}

	var err error
	var warning string

	// Check for a local development cache.
	if parsedPath.isRepoOnlyReference() {
		for _, directoryPath := range vcsDevelopmentDirectories {
			fullCheckDirectory := path.Join(directoryPath, parsedPath.url)
			if _, err := os.Stat(fullCheckDirectory); err == nil {
				// Found a local copy. Returning it.
				log.Printf("Found a local HEAD copy of package %s under development directory %s", parsedPath.url, directoryPath)
				warning = fmt.Sprintf(
					`Package '%s' does not specify a tag, commit or branch and a local copy was found under development directory '%s'. VCS checkout will be skipped and the local copy used instead. To return to normal VCS behavior for this package, remove the --vcs-dev-dir flag or specify the package's tag, commit or branch.`, parsedPath.String(), directoryPath)
				return fullCheckDirectory, nil, warning
			}
		}
	}

	// Conduct the checkout or pull.
	fullCacheDirectory := path.Join(pkgCacheRootPath, parsedPath.cacheDirectory())
	err, warning = checkCacheAndPull(parsedPath, fullCacheDirectory, cacheOption)

	// Warn if the package is a HEAD checkout.
	if err == nil && warning == "" && parsedPath.isHEAD() {
		warning = fmt.Sprintf("Package '%s' points to HEAD of a branch or commit and will be updated on every build", parsedPath.String())
	}

	// If the parsed path is a subdirectory of the checkout, return a reference to it.
	if err == nil && parsedPath.subpackage != "" {
		subpackageCacheDirectory := path.Join(fullCacheDirectory, parsedPath.subpackage)
		if _, serr := os.Stat(subpackageCacheDirectory); os.IsNotExist(serr) {
			return "", fmt.Errorf("Subpackage '%s' does not exist under VCS package '%s'", parsedPath.subpackage, parsedPath.url), warning
		}
		return subpackageCacheDirectory, nil, warning
	}

	return fullCacheDirectory, err, warning
}

// InspectInfo holds all the data returned from a call to PerformVCSCheckoutAndInspect.
type InspectInfo struct {
	Engine    string
	CommitSHA string
	Tags      []string
}

// PerformVCSCheckoutAndInspect performs the checkout and updating of the given VCS path and returns
// the commit SHA of the package, as well as its tags.
//
// pkgCacheRootPath holds the path of the root directory that forms the package cache.
//
// vcsDevelopmentDirectories specifies optional directories to check for branchless and tagless copies
// of the repository first. If found, the copy will be used in lieu of a normal checkout.
func PerformVCSCheckoutAndInspect(vcsPath string, pkgCacheRootPath string, cacheOption VCSCacheOption, vcsDevelopmentDirectories ...string) (InspectInfo, error, string) {
	directory, err, warning := PerformVCSCheckout(vcsPath, pkgCacheRootPath, cacheOption, vcsDevelopmentDirectories...)
	if err != nil {
		return InspectInfo{}, err, warning
	}

	handler, _ := detectHandler(directory)
	sha, err := handler.inspect(directory)
	if err != nil {
		return InspectInfo{}, err, warning
	}

	tags, err := handler.listTags(directory)
	if err != nil {
		return InspectInfo{}, err, warning
	}

	return InspectInfo{string(handler.kind), sha, tags}, nil, warning
}

// checkCacheAndPull conducts the cache check and necessary pulls.
func checkCacheAndPull(parsedPath vcsPackagePath, fullCacheDirectory string, cacheOption VCSCacheOption) (error, string) {
	// TODO(jschorr): Should we delete the package cache directory here if there was an error?

	// Check the package cache for the path.
	log.Printf("Checking cache directory %s", fullCacheDirectory)
	if _, err := os.Stat(fullCacheDirectory); os.IsNotExist(err) {
		// Do a full checkout.
		return performFullCheckout(parsedPath, fullCacheDirectory)
	}

	// If the cache exists, we only perform an update to the VCS package if the package
	// is marked as pointing to HEAD of a branch or commit. Tagged VCS packages are always left alone.
	log.Printf("Cache directory %s exists", fullCacheDirectory)
	if parsedPath.isHEAD() && cacheOption == VCSFollowNormalCacheRules {
		return performUpdateCheckout(parsedPath, fullCacheDirectory)
	}

	log.Printf("Cache directory %s exists and points to tag %s; no update needed", fullCacheDirectory, parsedPath.tag)
	return nil, ""
}

// performFullCheckout performs a full VCS checkout of the given package path.
func performFullCheckout(path vcsPackagePath, fullCacheDirectory string) (error, string) {
	// Lookup the VCS discovery information.
	discovery, err := DiscoverVCSInformation(path.url)
	if err != nil {
		return err, ""
	}

	// Perform a full checkout.
	handler, ok := vcsByKind[discovery.Kind]
	if !ok {
		panic("Unknown VCS handler")
	}

	warning, err := handler.checkout(path, discovery, fullCacheDirectory)
	return err, warning
}

// performUpdateCheckout performs a VCS update of the given package path.
func performUpdateCheckout(path vcsPackagePath, fullCacheDirectory string) (error, string) {
	// Detect the kind of VCS based on the checkout.
	handler, ok := detectHandler(fullCacheDirectory)
	if !ok {
		return fmt.Errorf("Could not detect VCS for directory: %s", fullCacheDirectory), ""
	}

	// If the checkout has changes, warn, but nothing more to do.
	if handler.check(fullCacheDirectory) {
		warning := fmt.Sprintf("VCS Package '%s' has changes on the local file system and will therefore not be updated", path.String())
		return nil, warning
	}

	// Otherwise, perform a pull to update.
	err := handler.update(fullCacheDirectory)
	return err, ""
}
